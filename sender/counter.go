/*
 * skogul, counter sender
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
	"sync"
	"time"

	"github.com/telenornms/skogul"
)

var countLog = skogul.Logger("sender", "count")

/*
Counter sender emits, periodically, the flow-rate of metrics through
it. The stats are sent on to the Stats-sender every Period.

To avoid locks and support multiple go routines using the same counter,
stats are sent over a channel to a separate goroutine that does the
actual aggregation and calculation.
*/
type Counter struct {
	Next   skogul.SenderRef  `doc:"Reference to the next sender in the chain"`
	Stats  skogul.HandlerRef `doc:"Handler that will receive the stats periodically"`
	Period skogul.Duration   `doc:"How often to emit stats" example:"5s"`
	ch     chan count
	once   sync.Once
	up     bool
}

// Internal count,
type count struct {
	ts         *time.Time
	containers int64
	values     int64
	metrics    int64
}

// Set up a Counter for the first time.
func (co *Counter) init() {
	co.ch = make(chan count, 100)
	co.up = true
	if co.Period.Duration == 0 {
		countLog.Debug("No Period set for Counter-sender. Using 1 second intervals.")
		co.Period.Duration = 1 * time.Second
	}
	go co.getIt()
}

// Send counts metrics, sends the count on a channel, then executes
// the next sender in the chain.
func (co *Counter) Send(c *skogul.Container) error {
	co.once.Do(func() {
		co.init()
	})
	var tmpc count
	tmpc.containers = 1
	for _, m := range c.Metrics {
		tmpc.metrics++
		for range m.Data {
			tmpc.values++
		}
	}
	x := time.Now()
	tmpc.ts = &x
	co.ch <- tmpc
	return co.Next.S.Send(c)
}

// Eat count-objects, once co.Period has passed, send them on.
//
// XXX: The sending should probably be a separate go routine for optimal
// performance and reliable stats, since the current implementation will
// stop sending updates if no stats are received. The _math_ will still
// be correct (Afaik), but right now this can't be used to verify that
// skogul is actually running, since there's no way to tell if Skogul is
// simply not receiving data or if it died.
func (co *Counter) getIt() {
	var current count
	var total count
	last := time.Now()

	for {
		m := <-co.ch
		current.containers += m.containers
		current.metrics += m.metrics
		current.values += m.values
		if m.ts.Sub(last) > co.Period.Duration {
			container := skogul.Container{}
			metric := skogul.Metric{}
			metric.Metadata = make(map[string]interface{})
			metric.Data = make(map[string]interface{})
			metric.Metadata["skogul"] = "counter"
			container.Metrics = []*skogul.Metric{&metric}
			total.containers += current.containers
			total.metrics += current.metrics
			total.values += current.values
			rate := count{
				containers: current.containers * int64(time.Second) / int64(m.ts.Sub(last)),
				metrics:    current.metrics * int64(time.Second) / int64(m.ts.Sub(last)),
				values:     current.values * int64(time.Second) / int64(m.ts.Sub(last))}
			container.Metrics[0].Time = m.ts
			container.Metrics[0].Data["total_containers"] = total.containers
			container.Metrics[0].Data["total_metrics"] = total.metrics
			container.Metrics[0].Data["total_values"] = total.values
			container.Metrics[0].Data["rate_containers"] = rate.containers
			container.Metrics[0].Data["rate_metrics"] = rate.metrics
			container.Metrics[0].Data["rate_values"] = rate.values
			if err := co.Stats.H.TransformAndSend(&container); err != nil {
				countLog.WithError(err).Error("Unable to transform and send counter stats")
			}
			current = count{nil, 0, 0, 0}
			last = *m.ts
		}
	}
}
