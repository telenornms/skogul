/*
 * skogul, influxdb writer
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.com>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
 * 02110-1301  USA
 */

package sender

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/telenornms/skogul"
)

var influxLog = skogul.Logger("sender", "influxdb")

/*
InfluxDB posts data to the provided URL and measurement, using the InfluxDB
line format over HTTP.
*/
type InfluxDB struct {
	URL                     string          `doc:"URL to InfluxDB API. Must include write end-point and database to write to." example:"http://[::1]:8086/write?db=foo"`
	Measurement             string          `doc:"Measurement name to write to."`
	MeasurementFromMetadata string          `doc:"Metadata key to read the measurement from. Either this or 'measurement' must be set. If both are present, 'measurement' will be used if the named metadatakey is not found."`
	Timeout                 skogul.Duration `doc:"HTTP timeout"`
	Insecure                bool            `doc:"Disable TLS certificate validation."`
	ConnsPerHost            int             `doc:"Max concurrent connections per host. Should reflect ulimit -n. Defaults to unlimited."`
	IdleConnsPerHost        int             `doc:"Max idle connections retained per host. Should reflect expected concurrency. Defaults to 2 + runtime.NumCPU."`
	RootCA                  string          `doc:"Path to an alternate root CA used to verify server certificates. Leave blank to use system defaults."`
	ConvertIntToFloat       bool            `doc:"Convert all integers to floats. Don't do this unless you really know why you're doing this."`
	Token                   skogul.Secret   `doc:"Authorization token used in InfluxDB 2.0"`
	client                  *http.Client
	replacer                *strings.Replacer
	once                    sync.Once
	stats                   influxStats
}

// influxStats holds runtime counters for the InfluxDB sender.
type influxStats struct {
	Received      atomic.Uint64 // Containers received via Send.
	Sent          atomic.Uint64 // Containers successfully POSTed to InfluxDB.
	Written       atomic.Uint64 // Individual metrics included in the line protocol body.
	Skipped       atomic.Uint64 // Metrics skipped (no data, missing measurement, bad tag/field).
	Errors        atomic.Uint64 // Generic errors (request build etc.).
	RequestErrors atomic.Uint64 // HTTP transport-level errors.
	HttpErrors    atomic.Uint64 // Non-2XX responses from InfluxDB.
}

// checkVariable verifies that the relevant variable is of a type we can
// handle.
func checkVariable(category string, field string, idx string, value interface{}) error {
	t := reflect.TypeOf(value)
	if t == nil {
		influxLog.Warnf("invalid %s: %s: %s: %v", category, field, idx, value)
		return fmt.Errorf("invalid %s: %s: %s: %v", category, field, idx, value)
	}

	k := t.Kind()

	switch k {
	case reflect.Bool:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Uintptr:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.String:
	default:
		influxLog.WithFields(logrus.Fields{
			"category": category,
			"field":    field,
			"index":    idx,
			"kind":     k,
		}).Info("Invalid tag/field data type. Flatten/convert data first.")
		return fmt.Errorf("bad tag/field")
	}
	return nil
}

// Send data to Influx, re-using idb.client.
func (idb *InfluxDB) Send(c *skogul.Container) error {
	var buffer bytes.Buffer
	idb.once.Do(func() {
		if idb.ConvertIntToFloat {
			influxLog.Warn("Influx sender is configured with 'ConvertIntToFloat'. This will convert *all* integers to floats.")
		}
		idb.replacer = strings.NewReplacer("\\", "\\\\", " ", "\\ ", ",", "\\,", "=", "\\=")
		if idb.Timeout.Duration == 0 {
			idb.Timeout.Duration = 20 * time.Second
		}
		cp, err := getCertPool(idb.RootCA)
		if err != nil {
			influxLog.Errorf("Failed to initialize root CA pool")
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: idb.Insecure,
			RootCAs:            cp,
		}
		iconsph := idb.IdleConnsPerHost
		if iconsph == 0 {
			iconsph = 2 + runtime.NumCPU()
		}
		tran := http.Transport{
			TLSClientConfig:     tlsConfig,
			MaxConnsPerHost:     idb.ConnsPerHost,
			MaxIdleConnsPerHost: iconsph,
		}

		idb.client = &http.Client{Transport: &tran, Timeout: idb.Timeout.Duration}
	})
	idb.stats.Received.Add(1)
	added := 0
	nmdata := 0
	ndata := 0
	for _, m := range c.Metrics {
		measurement := idb.Measurement
		if len(m.Data) == 0 {
			// must have SOME data
			// XXX: Should report.
			idb.stats.Skipped.Add(1)
			continue
		}
		if idb.MeasurementFromMetadata != "" {
			measure, ok := m.Metadata[idb.MeasurementFromMetadata].(string)
			if ok {
				measurement = measure
			}
			// The reason this isn't an else-if is because now
			// it also catches the scenario where the type cast
			// is successful, but the key is empty.
			if measurement == "" {
				// XXX:
				// How do we report issues of single
				// metrics failing, but not the container
				// in general? Failing the entire container
				// for just one failed metric is not really
				// acceptable...
				idb.stats.Skipped.Add(1)
				continue
			}
		}
		failed := 0
		for key, value := range m.Metadata {
			e1 := checkVariable("metadata", "key", "0", key)
			e2 := checkVariable("metadata", "value", key, value)
			if e1 != nil || e2 != nil {
				failed++
			}
		}
		for key, value := range m.Data {
			e1 := checkVariable("data", "key", "0", key)
			e2 := checkVariable("data", "value", key, value)
			if e1 != nil || e2 != nil {
				failed++
			}
		}
		if failed > 0 {
			idb.stats.Skipped.Add(1)
			continue
		}
		fmt.Fprintf(&buffer, "%s", measurement)
		for key, value := range m.Metadata {
			// Tag values and field values are handled differently;
			// A tag value is always a string, but if you wrap it in
			// quotes the quotes will be part of the tag value.
			// Therefore you need to escape any invalid character instead.
			// Run the replacer for tags (keys and values), and field keys,
			// but not for field values.
			var tagValue interface{}
			v, ok := value.(string)

			if ok {
				tagValue = idb.replacer.Replace(v)
				// Skip empty tag values, they are invalid
				// for Influx
				if tagValue == "" {
					continue
				}
			} else {
				tagValue = value
			}
			fmt.Fprintf(&buffer, ",%s=%v", idb.replacer.Replace(key), tagValue)
			nmdata++
		}
		fmt.Fprintf(&buffer, " ")
		comma := ""
		for key, value := range m.Data {

			fmt.Fprintf(&buffer, "%s%s=%s", comma, idb.replacer.Replace(key), idb.toInfluxValue(value))
			comma = ","
			ndata++
		}
		fmt.Fprintf(&buffer, " %d\n", m.Time.UnixNano())
		added++
	}
	idb.stats.Written.Add(uint64(added))
	if added == 0 {
		influxLog.Trace("Tried to send 0 metrics to influx. Probably no viable metrics after filtering out invalid tags and such. You may have to transform your data.")
		return nil
	}

	req, err := http.NewRequest("POST", idb.URL, &buffer)
	if err != nil {
		idb.stats.Errors.Add(1)
		return fmt.Errorf("unable to create request: %w", err)
	}
	if len(idb.Token) > 0 {
		req.Header.Add("authorization", fmt.Sprintf("Token %s", idb.Token.Expose()))
	}

	resp, err := idb.client.Do(req)
	if err != nil {
		idb.stats.RequestErrors.Add(1)
		return fmt.Errorf("unable to POST data: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		idb.stats.HttpErrors.Add(1)
		body, rerr := io.ReadAll(resp.Body)
		if rerr != nil {
			body = []byte("unable to read body")
		}
		if len(body) == 0 {
			body = fmt.Appendf(nil, "No reply body. Request: %s", buffer.Bytes())
		}

		return fmt.Errorf("influx sender(%s) failed to send container (%s). Bad response from InfluxDB: %s - %s", skogul.Identity[idb], c.Describe(), resp.Status, string(body))
	}
	idb.stats.Sent.Add(1)
	return nil
}

// GetStats prepares a skogul metric with stats for the InfluxDB sender.
func (idb *InfluxDB) GetStats() *skogul.Metric {
	now := skogul.Now()
	metric := skogul.Metric{
		Time:     &now,
		Metadata: make(map[string]interface{}),
		Data:     make(map[string]interface{}),
	}
	metric.Metadata["component"] = "sender"
	metric.Metadata["type"] = "influxdb"
	metric.Metadata["identity"] = skogul.Identity[idb]
	metric.Data["received"] = idb.stats.Received.Load()
	metric.Data["sent"] = idb.stats.Sent.Load()
	metric.Data["written"] = idb.stats.Written.Load()
	metric.Data["skipped"] = idb.stats.Skipped.Load()
	metric.Data["errors"] = idb.stats.Errors.Load()
	metric.Data["request_errors"] = idb.stats.RequestErrors.Load()
	metric.Data["http_errors"] = idb.stats.HttpErrors.Load()
	return &metric
}

// toInfluxValue handles converting values to values known by InfluxDB.
// E.g. an integer should end with the char 'i', so if the value is an int,
// we need to add that 'i'.
func (idb *InfluxDB) toInfluxValue(value interface{}) string {
	if !idb.ConvertIntToFloat {
		i, ok := value.(int64)
		if ok {
			return fmt.Sprintf("%di", i)
		}
		i2, ok := value.(int32)
		if ok {
			return fmt.Sprintf("%di", i2)
		}
		u, ok := value.(uint64)
		if ok {
			return fmt.Sprintf("%du", u)
		}
		u2, ok := value.(uint32)
		if ok {
			return fmt.Sprintf("%du", u2)
		}
		u3, ok := value.(uint)
		if ok {
			return fmt.Sprintf("%du", u3)
		}
	}
	return fmt.Sprintf("%#v", value)
}

// Verify does a shallow verification of settings
func (idb *InfluxDB) Verify() error {
	if idb.URL == "" {
		return skogul.MissingArgument("URL")
	}
	_, err := getCertPool(idb.RootCA)
	if err != nil {
		return fmt.Errorf("failed to read custom root CA (RootCA: %s): %w", idb.RootCA, err)
	}
	if idb.Measurement == "" && idb.MeasurementFromMetadata == "" {
		return skogul.MissingArgument("Measurement or MeasurementFromMetadata")
	}
	return nil
}
