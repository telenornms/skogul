/*
 * skogul, batch sender
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

// BUG(kly): The interval for the batch sender is really a timeout, not an
// interval at the moment.
import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/KristianLyng/skogul"
)

/*
Batch sender collects metrics into a single container then passes them on
after Threshold number of metrics are collected. In case Threshold is
"never" reached, it will periodically flush metrics if no message has been
received in Interval time.
*/
type Batch struct {
	Next      skogul.Sender
	Interval  time.Duration
	Threshold int
	allocSize int
	ch        chan *skogul.Container
	once      sync.Once
	metrics   int
	cont      *skogul.Container
}

func (bat *Batch) setup() {
	if bat.Threshold == 0 {
		bat.Threshold = 10
	}
	bat.ch = make(chan *skogul.Container, 10)
	if bat.Interval == 0 {
		bat.Interval = time.Duration(1 * time.Second)
	}
	slack := int(bat.Threshold / 5)
	bat.allocSize = bat.Threshold + slack
	go bat.run()
}

// add is a custom-written append(), so to speak, since we know that the
// maximum size of the relevant slice is Threshold.
func (bat *Batch) add(c *skogul.Container) {
	if bat.cont == nil {
		bat.cont = &skogul.Container{}
	}

	cl := len(bat.cont.Metrics)
	cc := cap(bat.cont.Metrics)
	nl := len(c.Metrics) + cl

	if nl > cc {
		newlen := bat.allocSize
		// It's allowed to exceed Threshold - but only once.
		if newlen < nl {
			newlen = nl
			log.Print((skogul.Error{Source: "batch sender", Reason: fmt.Sprintf("Warning: slice too small for 20%% slack - need to resize/copy. Performance hit :D(Default alloc size is %d, need %d)", bat.allocSize, newlen)}))
		}
		x := make([]*skogul.Metric, newlen)
		copy(x, bat.cont.Metrics)
		bat.cont.Metrics = x
	}

	bat.cont.Metrics = bat.cont.Metrics[0:nl]
	copy(bat.cont.Metrics[cl:nl], c.Metrics)
}

func (bat *Batch) flush() {
	err := bat.Next.Send(bat.cont)
	if err != nil {
		log.Print(skogul.Error{Source: "batch sender", Reason: "down stream error", Next: err})
	}
	bat.cont = nil
}

func (bat *Batch) run() {
	for {
		select {
		case c := <-bat.ch:
			bat.add(c)
			if len(bat.cont.Metrics) >= bat.Threshold {
				bat.flush()
			}
		case <-time.After(bat.Interval):
			bat.flush()
		}
	}
}

// Send batches up multiple metrics and passes them on after an interval or
// a set size is reached. It never returns error, since there is no way to
// know.
func (bat *Batch) Send(c *skogul.Container) error {
	bat.once.Do(func() {
		bat.setup()
	})

	bat.ch <- c
	return nil
}
