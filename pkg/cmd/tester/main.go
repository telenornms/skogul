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
	"flag"
	"fmt"
	gollector "github.com/KristianLyng/gollector/pkg"
	senders "github.com/KristianLyng/gollector/pkg/senders"
	"log"
	"math/rand"
	"time"
)

var metrics = flag.Int64("metrics", 1000, "Number of metrics per HTTP post")
var values = flag.Int64("values", 5, "Number of values per metric")

func generate(t time.Time) (gollector.Container, int64) {
	c := gollector.Container{}
	c.Template.Time = &t
	c.Metrics = make([]gollector.Metric, *metrics)
	for i := int64(0); i < *metrics; i++ {
		m := gollector.Metric{}
		m.Metadata = map[string]interface{}{}
		m.Metadata["key1"] = i
		m.Data = map[string]interface{}{}
		for key := int64(0); key < *values; key++ {
			m.Data[fmt.Sprintf("metric%d", key)] = rand.Int63()
		}
		c.Metrics[i] = m
	}
	return c, *metrics * *values
}

func main() {
	flag.Parse()
	t := time.Now().Add(time.Second * -3600)
	sender := senders.HTTP{"http://[::1]:8080"}
	start := time.Now()
	end := time.Now()
	var mets int64
	mets = 0
	for runs := 0; runs < 3600; runs++ {
		c, nm := generate(t.Add(time.Second * time.Duration(runs)))
		mets += nm
		err := sender.Send(&c)
		if err != nil {
			log.Print(err)
		}
		if (runs % 10) == 0 {
			end = time.Now()
			log.Printf("Run %d, %d metrics/s", runs, mets*int64(time.Second)/int64(end.Sub(start)))
			mets = 0
			start = time.Now()

		}
	}
}
