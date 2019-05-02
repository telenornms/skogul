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
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"
)

/*
InfluxDB posts data to the provided URL and measurement, using the InfluxDB
line format over HTTP.
*/
type InfluxDB struct {
	URL         string
	Measurement string
	client      *http.Client
	mux         sync.Mutex
}

func init() {
	addAutoSender("influx", NewInflux, "Send InfluxDB data to a HTTP endpoint, using the first element of the path as db and second as measurement, e.g: influx://host/db/measurement")
}

/*
NewInflux returns a new InfluxDB sender, parsing the URL more than the regular
InfluxDB.URL field. Since it is mapped to the "influx" schema, the url
provided is not used as-is, and instead follows the spec of
influx://host[:port]/database/measurement
*/
func NewInflux(ul url.URL) skogul.Sender {
	db, measurement := path.Split(ul.Path)
	t := fmt.Sprintf("http://%s/write?db=%s", ul.Host, db[1:len(db)-1])
	x := InfluxDB{URL: t, Measurement: measurement}
	return &x
}

// Send data to Influx, re-using idb.client.
func (idb *InfluxDB) Send(c *skogul.Container) error {
	var buffer bytes.Buffer
	for _, m := range c.Metrics {
		fmt.Fprintf(&buffer, "%s", idb.Measurement)
		for key, value := range m.Metadata {
			fmt.Fprintf(&buffer, ",%s=%#v", key, value)
		}
		fmt.Fprintf(&buffer, " ")
		comma := ""
		for key, value := range m.Data {
			fmt.Fprintf(&buffer, "%s%s=%#v", comma, key, value)
			comma = ","
		}
		fmt.Fprintf(&buffer, " %d\n", m.Time.UnixNano())
	}
	if idb.client == nil {
		idb.mux.Lock()
		// Recheck after acquiring lock
		if idb.client == nil {
			idb.client = &http.Client{Timeout: 5 * time.Second}
		}
		idb.mux.Unlock()
	}
	resp, err := idb.client.Post(idb.URL, "text/plain", &buffer)
	if err != nil {
		e := skogul.Error{Source: "influxdb sender", Reason: "unable to POST data", Next: err}
		log.Print(e)
		return e
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		e := skogul.Error{Source: "influxdb sender", Reason: fmt.Sprintf("bad response code from InfluxDB: %d", resp.StatusCode)}
		log.Print(e)
		return e
	}
	return nil
}
