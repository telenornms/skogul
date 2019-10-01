/*
 * skogul, test receiver
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

package receiver

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/KristianLyng/skogul"
)

// Tester synthesise dummy data.
type Tester struct {
	Metrics int64             `doc:"Number of metrics in each container"`
	Values  int64             `doc:"Number of unique values for each metric"`
	Threads int               `doc:"Threads to spawn"`
	Delay   skogul.Duration   `doc:"Sleep time between each metric is generated, if any."`
	Handler skogul.HandlerRef `doc:"Reference to a handler where the data is sent"`
}

func (tst *Tester) generate(t time.Time) skogul.Container {
	c := skogul.Container{}
	c.Template = &skogul.Metric{}
	c.Template.Time = &t
	c.Metrics = make([]*skogul.Metric, tst.Metrics)
	for i := int64(0); i < tst.Metrics; i++ {
		m := skogul.Metric{}
		m.Metadata = map[string]interface{}{}
		m.Metadata["key1"] = i
		m.Data = map[string]interface{}{}
		for key := int64(0); key < tst.Values; key++ {
			m.Data[fmt.Sprintf("metric%d", key)] = rand.Int63()
		}
		c.Metrics[i] = &m
	}
	return c
}

// Start never returns.
func (tst *Tester) Start() error {
	if tst.Threads == 0 {
		tst.Threads = runtime.NumCPU()
		log.Printf("No threads set, defaulting to runtime.NumCPU() (%d)", tst.Threads)
	}
	if tst.Metrics < 1 {
		tst.Metrics = 10
		log.Printf("No Metrics specified for testing, defaulting to %d", tst.Metrics)
	}
	if tst.Values < 1 {
		tst.Values = 50
		log.Printf("No Values specified for testing, defaulting to %d", tst.Values)
	}
	for i := 1; i < tst.Threads; i++ {
		go tst.run()
	}
	tst.run()
	return nil
}

// innerRun is a single iteration of generating metrics and sending them.
// It is separated from run() to simplify defer/delay. Without it, it's
// easy to introduce bugs where the delay isn't executed on transform
// failure.
func (tst *Tester) innerRun() {
	if tst.Delay.Duration != 0 {
		defer func(d time.Duration) {
			time.Sleep(d)
		}(tst.Delay.Duration)
	}
	c := tst.generate(time.Now())
	for _, t := range tst.Handler.H.Transformers {
		if e := t.Transform(&c); e != nil {
			log.Printf("Failed to transform test-metric: %v", e)
			return
		}
	}
	err := tst.Handler.H.Sender.Send(&c)
	if err != nil {
		log.Print(err)
	}
}

// run() is a signle thread of the tester that runs for ever.
func (tst *Tester) run() {
	for {
		tst.innerRun()
	}
}
