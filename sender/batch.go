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

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/KristianLyng/skogul"
)

/*
Batch sender collects metrics into a single container then passes them on
after Threshold number of metrics are collected. In case Threshold is
"never" reached, it will periodically flush metrics if no message has been
received in Interval time.

Internally, the Batch sender consists of three parts. The first part is
the Send() part, which just pushes the received container onto a channel.

The second part, which is a single, dedicated go routine, picks up said
container and adds it to a batch-container. When the batch container is
"full" (e.g.: exceeds Threshold) - or a timeout is reached - the batch
container is pushed onto a second channel and a new, empty, batch container
is created.

The third part picks up the ready-to-send containers and issues next.Send()
on them. This is a separate go routine, one per NumCPU.

This means that:

1. Batch sender will both do a "fan-in" from potentially multiple Send()
   calls.
2. ... and do a fan-out afterwards.
3. Send() will only block if two channels are full.
*/
type Batch struct {
	Next      skogul.SenderRef `doc:"Sender that will receive batched metrics"`
	Interval  skogul.Duration  `doc:"Flush the bucket after this duration regardless of how full it is"`
	Threshold int              `doc:"Flush the bucket after reaching this amount of metrics"`
	allocSize int
	ch        chan *skogul.Container
	once      sync.Once
	metrics   int
	timer     *time.Timer
	cont      *skogul.Container
	out       chan *skogul.Container
}

func (bat *Batch) setup() {
	if bat.Threshold == 0 {
		bat.Threshold = 10
	}
	bat.ch = make(chan *skogul.Container, 10)
	if bat.Interval.Duration == 0 {
		bat.Interval.Duration = time.Duration(1 * time.Second)
	}
	// The slack is just an array of pointers, so it's more important
	// to avoid unnecessary allocation than save a few bytes. With a
	// 1000-size threshold, the allocation would be for 1500 - or
	// around 6kB.
	slack := int(bat.Threshold / 2)
	bat.allocSize = bat.Threshold + slack
	if bat.allocSize < 100 {
		bat.allocSize = 100
	}
	bat.out = make(chan *skogul.Container, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go bat.flusher()
	}
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
			log.Print((skogul.Error{Source: "batch sender", Reason: fmt.Sprintf("Warning: slice too small for 50%% slack - need to resize/copy. Performance hit :D(Default alloc size is %d, need %d)", bat.allocSize, newlen)}))
		}
		x := make([]*skogul.Metric, newlen)
		copy(x, bat.cont.Metrics)
		bat.cont.Metrics = x
	}

	bat.cont.Metrics = bat.cont.Metrics[0:nl]
	copy(bat.cont.Metrics[cl:nl], c.Metrics)
}

// flusher fetches a ready-to-ship container and issues send(). One flusher
// is run per NumCPU
func (bat *Batch) flusher() {
	for {
		c := <-bat.out
		err := bat.Next.S.Send(c)
		if err != nil {
			log.Print(skogul.Error{Source: "batch sender", Reason: "down stream error", Next: err})
		}
	}
}

func (bat *Batch) flush() {
	bat.out <- bat.cont
	bat.cont = nil
}

func (bat *Batch) timerReschedule() {
	if !bat.timer.Stop() {
		<-bat.timer.C
	}
	bat.timer = time.NewTimer(bat.Interval.Duration)
}

// run runs forever, listening for containers, doing the actual batching,
// and triggering a flush if either the timer or threshold is reached. The
// actual sending is done in the flusher() go routines.
func (bat *Batch) run() {
	bat.timer = time.NewTimer(bat.Interval.Duration)
	for {
		select {
		case c := <-bat.ch:
			bat.add(c)
			if len(bat.cont.Metrics) >= bat.Threshold {
				bat.flush()
				bat.timerReschedule()
			}
		case <-bat.timer.C:
			bat.timer = time.NewTimer(bat.Interval.Duration)
			if bat.cont != nil {
				bat.flush()
			}
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
