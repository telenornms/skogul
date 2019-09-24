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
	URL         string `doc:"URL to InfluxDB API"`
	Measurement string `doc:"Measurement name to write to"`
	Timeout     skogul.Duration `doc:"HTTP timeout"`
	client      *http.Client
	replacer    *strings.Replacer
	once        sync.Once
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
		fmt.Fprintf(&buffer, "%s", idb.Measurement)
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
