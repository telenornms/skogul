/*
 * skogul, extremely dumb benchmarker
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
	"github.com/KristianLyng/skogul/pkg"
	"github.com/KristianLyng/skogul/pkg/senders"
	"log"
	"math/rand"
	"time"
)

var metrics = flag.Int64("metrics", 1000, "Number of metrics per HTTP post")
var values = flag.Int64("values", 5, "Number of values per metric")

func generate(t time.Time) skogul.Container {
	c := skogul.Container{}
	c.Template.Time = &t
	c.Metrics = make([]skogul.Metric, *metrics)
	for i := int64(0); i < *metrics; i++ {
		m := skogul.Metric{}
		m.Metadata = map[string]interface{}{}
		m.Metadata["key1"] = i
		m.Data = map[string]interface{}{}
		for key := int64(0); key < *values; key++ {
			m.Data[fmt.Sprintf("metric%d", key)] = rand.Int63()
		}
		c.Metrics[i] = m
	}
	return c
}

func main() {
	flag.Parse()
	t := time.Now().Add(time.Second * -3600)
	skoup := senders.HTTP{"http://[::1]:8080"}
	sender := &senders.Counter{Next: skoup, Stats: senders.Debug{}}
	for runs := 0; runs < 3600; runs++ {
		c := generate(t.Add(time.Second * time.Duration(runs)))
		err := sender.Send(&c)
		if err != nil {
			log.Print(err)
		}
	}
}
