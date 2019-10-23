/*
 * skogul, detach sender
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
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/telenornms/skogul"
)

/*
These senders are somewhat controversial, as they affect the fanout/fan-in
behavior of Skogul. If you use detacher on a receiver that already does multiple
go routines, you end up with fan-in. You could couple it with the Fanout sender
to go back to multiple go routines, but yeah.... this is not without its drawbacks.

However, they ARE useful, just make sure you think through side effects of using
these senders.
*/

/*
Detacher accepts a message, sends it to a channel, then picks it up on the
other end in a separate go routine. This, unfortunately, leads to fan-in:
if used in conjunction with HTTP receiver, for example, you end up going from
multiple independent go routines to a single one, which is probably not what
you want.

The purpose is to smooth out reading.
*/
type Detacher struct {
	Next  skogul.SenderRef `doc:"Sender that receives the metrics."`
	Depth int              `doc:"How many containers can be pending delivery before we start blocking. Defaults to 1000."`
	ch    chan *skogul.Container
	once  sync.Once
}

// consume is the detached go routine that picks up containers and passes
// them on.
func (de *Detacher) consume() {
	for c := range de.ch {
		de.Next.S.Send(c)
	}
}

func (de *Detacher) doInit() {
	if de.Depth == 0 {
		de.Depth = 1000
		log.WithField("depth", de.Depth).Debug("No detach depth/queue depth set. Using default value.")
	}
	de.ch = make(chan *skogul.Container, de.Depth)
	go de.consume()
}

// Send ensures a consumer exists, then transmits the container on a
// channel and returns immediately.
func (de *Detacher) Send(c *skogul.Container) error {
	de.once.Do(func() {
		de.doInit()
	})
	de.ch <- c
	return nil
}

/*
Fanout sender implements a worker pool for passing data on. This SHOULD be
unnecessary, as the receiver should ideally do this for us (e.g.: the
HTTP receiver does this natively). However, there might be times
where it makes sense, specially since this can be used in reverse too:
you can use the Fanout sender to limit the degree of concurrency that
downstream is exposed to.

Again, this should really not be needed. If you use the fanout sender, be
sure you understand why.

There only settings provided is "Next" to provide the next sender, and
"Workers", that defines the size of the worker pool.

*/
type Fanout struct {
	Next    skogul.SenderRef `doc:"Sender receiving the metrics"`
	Workers int              `doc:"Number of worker threads in use. To _fan_in_ you can set this to 1."`
	once    sync.Once
	workers chan chan *skogul.Container
}

func (fo *Fanout) doInit() {
	if fo.Workers == 0 {
		n := runtime.NumCPU()
		fo.Workers = n
		log.WithField("fanout", fo.Workers).Debug("No fanout size set. Using default value")
	}
	fo.workers = make(chan chan *skogul.Container)
	for i := 0; i < fo.Workers; i++ {
		go fo.worker()
	}
}

// Send ensures the workers are booted, then picks up a channel from
// available workers and sends the container to that container.
func (fo *Fanout) Send(c *skogul.Container) error {
	fo.once.Do(func() {
		fo.doInit()
	})
	x := <-fo.workers
	x <- c
	return nil
}

// worker makes a channel for work, makes that channel available on the
// shared fo.workers channel, then reads from it.
func (fo *Fanout) worker() {
	c := make(chan *skogul.Container)
	for {
		fo.workers <- c
		con := <-c
		fo.Next.S.Send(con)
	}
}
