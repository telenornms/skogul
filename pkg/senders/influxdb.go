/*
 * gollector, influxdb writer
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

package senders

import (
	"bytes"
	"fmt"
	. "github.com/KristianLyng/gollector/pkg/common"
	"log"
	"net/http"
	"time"
)

type InfluxDB struct {
	Url         string
	Measurement string
}

func (idb InfluxDB) Send(c *GollectorContainer) error {
	var buffer bytes.Buffer
	for _, m := range c.Metrics {
		fmt.Fprintf(&buffer, "%s", idb.Measurement)
		for key, value := range c.Template.Metadata {
			fmt.Fprintf(&buffer, ",%s=%#v", key, value)
		}
		for key, value := range m.Metadata {
			fmt.Fprintf(&buffer, ",%s=%#v", key, value)
		}
		fmt.Fprintf(&buffer, " ")
		comma := ""
		for key, value := range m.Data {
			fmt.Fprintf(&buffer, "%s%s=%#v", comma, key, value)
			comma = ","
		}
		lt := c.Template.Time
		if m.Time != nil {
			lt = m.Time
		}
		fmt.Fprintf(&buffer, " %d\n", lt.UnixNano())
	}
	req, err := http.NewRequest("POST", idb.Url, &buffer)
	req.Header.Set("Content-Type", "text/plain")
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Print(resp)
	}
	return nil
}

//req, err := http.NewRequest("POST", "http://127.0.0.1:8086/write?db=test", &buffer)
