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
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/KristianLyng/skogul"
)

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
	for _, m := range c.Metrics {
		measurement := idb.Measurement
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
		}
		fmt.Fprintf(&buffer, " ")
		comma := ""
		for key, value := range m.Data {
			fmt.Fprintf(&buffer, "%s%s=%#v", comma, idb.replacer.Replace(key), value)
			comma = ","
		}
		fmt.Fprintf(&buffer, " %d\n", m.Time.UnixNano())
	}
	resp, err := idb.client.Post(idb.URL, "text/plain", &buffer)
	if err != nil {
		e := skogul.Error{Source: "influxdb sender", Reason: "unable to POST data", Next: err}
		return e
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		e := skogul.Error{Source: "influxdb sender", Reason: fmt.Sprintf("bad response from InfluxDB: %s", resp.Status)}
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
