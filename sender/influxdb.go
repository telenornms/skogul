/*
 * skogul, influxdb writer
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
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
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

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
	client                  *http.Client
	replacer                *strings.Replacer
	once                    sync.Once
}

// checkVariable verifies that the relevant variable is of a type we can
// handle.
func checkVariable(category string, field string, idx string, value interface{}) error {
	t := reflect.TypeOf(value)
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
		e := skogul.Error{Source: "influxdb sender", Reason: fmt.Sprintf("bad tag/field")}
		return e
	}
	return nil
}

// Send data to Influx, re-using idb.client.
func (idb *InfluxDB) Send(c *skogul.Container) error {
	var buffer bytes.Buffer
	idb.once.Do(func() {
		idb.replacer = strings.NewReplacer("\\", "\\\\", " ", "\\ ", ",", "\\,", "=", "\\=")
		if idb.Timeout.Duration == 0 {
			idb.Timeout.Duration = 20 * time.Second
		}
		idb.client = &http.Client{Timeout: idb.Timeout.Duration}
	})
	added := 0
	nmdata := 0
	ndata := 0
	for _, m := range c.Metrics {
		measurement := idb.Measurement
		if len(m.Data) == 0 {
			// must have SOME data
			// XXX: Should rneport.
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
			continue
		}
		fmt.Fprintf(&buffer, "%s", measurement)
		for key, value := range m.Metadata {
			var field interface{}
			v, ok := value.(string)
			if ok {
				field = idb.replacer.Replace(v)
			} else {
				field = value
			}
			fmt.Fprintf(&buffer, ",%s=%v", idb.replacer.Replace(key), field)
			nmdata++
		}
		fmt.Fprintf(&buffer, " ")
		comma := ""
		for key, value := range m.Data {

			fmt.Fprintf(&buffer, "%s%s=%#v", comma, idb.replacer.Replace(key), value)
			comma = ","
			ndata++
		}
		fmt.Fprintf(&buffer, " %d\n", m.Time.UnixNano())
		added++
	}
	if added == 0 {
		influxLog.Trace("Tried to send 0 metrics to influx. Probably no viable metrics after filtering out invalid tags and such. You may have to transform your data.")
		return nil
	}

	resp, err := idb.client.Post(idb.URL, "text/plain", &buffer)
	if err != nil {
		e := skogul.Error{Source: "influxdb sender", Reason: "unable to POST data", Next: err}
		influxLog.Trace(e)
		return e
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var body []byte
		if resp.ContentLength > 0 {
			body = make([]byte, resp.ContentLength)

			if _, err := io.ReadFull(resp.Body, body); err != nil {
				body = []byte(`unable to ready body`)
			}
		} else {
			body = []byte(fmt.Sprintf("No reply body. Request: %s", buffer.Bytes()))
		}

		influxLog.WithFields(logrus.Fields{
			"metricsSent": added,
			"body":        string(body),
			"status":      resp.Status,
			"tags":        nmdata,
			"fields":      ndata,
		}).Error("Failed to send data to influx. Bad response from backend.")
		e := skogul.Error{Source: "influxdb sender", Reason: fmt.Sprintf("bad response from InfluxDB: %s - %s", resp.Status, string(body))}
		return e
	}
	return nil
}

// Verify does a shallow verification of settings
func (idb *InfluxDB) Verify() error {
	if idb.URL == "" {
		return skogul.Error{Source: "influxdb sender", Reason: "no URL set"}
	}
	if idb.Measurement == "" && idb.MeasurementFromMetadata == "" {
		return skogul.Error{Source: "influxdb sender", Reason: "no Measurement set or MeasurementFromMetadata"}
	}
	return nil
}
