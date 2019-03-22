/*
 * skogul, debug sender
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
	"encoding/json"
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"math/rand"
	"time"
)

/*
Debug sender simply prints the metrics in json-marshalled format to
stdout.
*/
type Debug struct {
}

func (db Debug) Send(c *skogul.Container) error {
	b, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		log.Panic("Unable to marshal json for debug output: %s", err)
		return err
	}
	log.Printf("Debug: \n%s", b)
	return nil
}

/*
The Sleeper-sender injects a random delay between 0 and MaxDelay before
passing execution over to the Next sender.

The purpose is testing.
*/
type Sleeper struct {
	Next     skogul.Sender
	MaxDelay time.Duration
	Verbose  bool
}

func (sl *Sleeper) Send(c *skogul.Container) error {
	d := rand.Float64() * float64(sl.MaxDelay)
	if sl.Verbose {
		log.Printf("Sleeping for %v", time.Duration(d))
	}
	time.Sleep(time.Duration(d))
	return sl.Next.Send(c)
}

/*
The Counter sender emits, periodically, the flow-rate of metrics through
it. The stats are sent on to the Stats-sender every Period.
*/
type Counter struct {
	Next      skogul.Sender
	Stats     skogul.Sender
	Period    time.Duration
	container *skogul.Container
	metric    *skogul.Metric
	last      time.Time
	current   count
	total     count
}

type count struct {
	containers int64
	values     int64
	metrics    int64
}

// Just set up the basics....
func (co *Counter) makeContainer() {
	co.container = &skogul.Container{}
	metric := skogul.Metric{}
	metric.Metadata = make(map[string]interface{})
	metric.Data = make(map[string]interface{})
	metric.Metadata["skogul"] = "counter"
	co.container.Metrics = append(co.container.Metrics, metric)
}

func (co *Counter) Send(c *skogul.Container) error {
	co.current.containers++
	for _, m := range c.Metrics {
		co.current.metrics++
		for range m.Data {
			co.current.values++
		}
	}
	now := time.Now()
	if now.Sub(co.last) > co.Period {
		if !co.last.IsZero() {
			co.total.containers += co.current.containers
			co.total.metrics += co.current.metrics
			co.total.values += co.current.values
			rate := count{
				containers: co.current.containers * int64(time.Second) / int64(now.Sub(co.last)),
				metrics:    co.current.metrics * int64(time.Second) / int64(now.Sub(co.last)),
				values:     co.current.values * int64(time.Second) / int64(now.Sub(co.last))}
			if co.container == nil {
				co.makeContainer()
			}
			co.container.Metrics[0].Time = &now
			co.container.Metrics[0].Data["total_containers"] = co.total.containers
			co.container.Metrics[0].Data["total_metrics"] = co.total.metrics
			co.container.Metrics[0].Data["total_values"] = co.total.values
			co.container.Metrics[0].Data["rate_containers"] = rate.containers
			co.container.Metrics[0].Data["rate_metrics"] = rate.metrics
			co.container.Metrics[0].Data["rate_values"] = rate.values
			co.Stats.Send(co.container)
		} else if co.Period == 0 {
			// Stupid way to set a default.... Must be a better
			// way?
			co.Period = 5 * time.Second
		}
		co.current = count{0, 0, 0}
		co.last = now
	}
	return co.Next.Send(c)
}
