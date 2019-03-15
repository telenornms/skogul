/*
 * gollector, extremely dumb benchmarker
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

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/KristianLyng/gollector/pkg/common"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var metrics = flag.Int64("metrics", 1000, "Number of metrics per HTTP post")
var values = flag.Int64("values", 5, "Number of values per metric")

func meh(t time.Time) int64 {
	c := GollectorContainer{}
	c.Template.Time = &t
	c.Metrics = make([]GollectorMetric, *metrics)
	for i := int64(0); i < *metrics; i++ {
		m := GollectorMetric{}
		m.Metadata = map[string]interface{}{}
		m.Metadata["key1"] = i
		m.Data = map[string]interface{}{}
		for key := int64(0); key < *values; key++ {
			m.Data[fmt.Sprintf("metric%d", key)] = rand.Int63()
		}
		c.Metrics[i] = m
	}
	b, err := json.Marshal(c)
	var buffer bytes.Buffer
	buffer.Write(b)
	req, err := http.NewRequest("POST", "http://[::1]:8080", &buffer)
	req.Header.Set("Content-Type", "application/json")
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	} else {
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Print(resp)
		}
	}
	return *metrics * *values
}

func main() {
	flag.Parse()
	t := time.Now().Add(time.Second * -3600)
	start := time.Now()
	end := time.Now()
	var mets int64
	mets = 0
	for runs := 0; runs < 3600; runs++ {
		mets += meh(t.Add(time.Second * time.Duration(runs)))
		if (runs % 10) == 0 {
			end = time.Now()
			log.Printf("Run %d, %d metrics/s", runs, mets*int64(time.Second)/int64(end.Sub(start)))
			mets = 0
			start = time.Now()

		}
	}
}
