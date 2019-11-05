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
	"math/rand"
	"runtime"
	"time"

	"github.com/telenornms/skogul"
)

var testerLog = skogul.Logger("receiver", "tester")

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
		testerLog.WithField("threads", tst.Threads).Debug("No threads set, defaulting to runtime.NumCPU()")
	}
	if tst.Metrics < 1 {
		tst.Metrics = 10
		testerLog.WithField("metrics", tst.Metrics).Debug("No Metrics specified for testing, defaulting to default value")
	}
	if tst.Values < 1 {
		tst.Values = 50
		testerLog.WithField("values", tst.Values).Debug("No Values specified for testing, defaulting to default value")
	}
	for i := 1; i < tst.Threads; i++ {
		go tst.run()
	}
	tst.run()
	return nil
}

// run() is a signle thread of the tester that runs for ever.
func (tst *Tester) run() {
	for {
		c := tst.generate(time.Now())
		if err := tst.Handler.H.TransformAndSend(&c); err != nil {
			testerLog.WithError(err).Errorf("Failed to transform and send metrics: %v", err)
		}
		time.Sleep(tst.Delay.Duration)
	}
}
