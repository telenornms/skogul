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
	"runtime"
	"sync"
	"time"

	"github.com/telenornms/skogul"
)

var batchLog = skogul.Logger("sender", "batch")

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
	Threads   int              `doc:"Number of threads for batch sender. Defaults to number of CPU cores."`
	Burner	  skogul.SenderRef `doc:"If the next sender is too slow and containers pile up beyond the backlog, instead of blocking waiting for the regular sender, redirect the overflown data to this burner. If left blank, the batcher will block, waiting for the regular sender. Note that there is no secondary overflow protection, so if the burner-sender is slow as well, the batcher will still block. To just discard the data, just use the null-sender. To measure how frequently this happens, use the count-sender in combination with the null-sender."`
	allocSize int // Precomputed of container to allocate initially
	ch        chan *skogul.Container // Initial channel used from Send()
	once      sync.Once
	timer     *time.Timer
	cont      *skogul.Container // Current container - used single threaded
	out       chan *skogul.Container // When Thershold/Timer is triggered, dump the container here
	burner    *chan *skogul.Container // Or burn it. Points to "out" if no burner is configured.
}

func (bat *Batch) setup() {
	if bat.Threads == 0 {
		bat.Threads = runtime.NumCPU()
	}
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
	bat.out = make(chan *skogul.Container, bat.Threads)
	for i := 0; i < bat.Threads; i++ {
		go bat.flusher(bat.out, bat.Next.S)
	}
	if bat.Burner.Name != "" {
		burner := make(chan *skogul.Container, bat.Threads)
		bat.burner = &burner
		go bat.flusher(burner, bat.Burner.S)
	} else {
		bat.burner = &bat.out
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
			batchLog.Warning((skogul.Error{Source: "batch sender", Reason: fmt.Sprintf("Warning: slice too small for 50%% slack - need to resize/copy. Performance hit :D(Default alloc size is %d, need %d)", bat.allocSize, newlen)}))
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
func (bat *Batch) flusher(ch chan *skogul.Container, sender skogul.Sender) {
	for {
		c := <-ch
		err := sender.Send(c)
		if err != nil {
			batchLog.WithError(err).Error(skogul.Error{Source: "batch sender", Reason: "down stream error", Next: err})
		}
	}
}

// flush is a "non-blocking" flush from the single-threaded part of the
// batcher. It just dumps the container on to a channel, if the channel is
// blocked, it will use an alternate channel. bat.burner will just point
// back to bat.out if no burner is present, thus block.
func (bat *Batch) flush() {
	select {
		case bat.out <- bat.cont:
		default:
			*bat.burner <- bat.cont
	}
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
