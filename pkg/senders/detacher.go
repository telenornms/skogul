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

package senders

import (
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"sync"
)

/*
XXX: The senders in detacher.go should probably be ignored for now. There are some
conceptual issues here that need to be resolved. E.g.: Should "fanout" even
exist? Should "detacher" just use "go next.Send(...)" - if so, how to avoid
exploding? etc.
*/

/*
Deatcher accepts a message, sends it to a channel, then picks it up on the
other end in a separate go routine. This, unfortunately, leades to fan-in:
if used in conjunction with HTTP receiver, for example, you end up going from
multiple independent go routines to a single one, which is probably not what
you want.

The purpose is to smooth out reading.
*/
type Detacher struct {
	Next  skogul.Sender
	Depth int
	init  bool
	ch    chan *skogul.Container
	mux   sync.Mutex
}

func (de *Detacher) Consume() {
	for c := range de.ch {
		de.Next.Send(c)
	}
}

func (de *Detacher) Init() {
	if de.init {
		return
	}
	de.mux.Lock()
	defer de.mux.Unlock()
	if de.init {
		return
	}
	if de.Depth == 0 {
		log.Print("No detach depth/queue depth set. Using default of 1000.")
		de.Depth = 1000
	}
	de.ch = make(chan *skogul.Container, de.Depth)
	go de.Consume()
	de.init = true
}

func (de *Detacher) Send(c *skogul.Container) error {
	de.Init()
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
	Next    skogul.Sender
	Workers int
	init    bool
	mux     sync.Mutex
	workers chan chan *skogul.Container
}

func (fo *Fanout) Init() {
	if fo.init {
		return
	}
	fo.mux.Lock()
	defer fo.mux.Unlock()
	if fo.init {
		return
	}
	if fo.Workers == 0 {
		log.Print("No fanout size set. Using default of 10.")
		fo.Workers = 10
	}
	fo.workers = make(chan chan *skogul.Container)
	for i := 0; i < fo.Workers; i++ {
		go fo.worker()
	}
	fo.init = true
}

func (fo *Fanout) Send(c *skogul.Container) error {
	fo.Init()
	x := <-fo.workers
	x <- c
	return nil
}

func (fo *Fanout) worker() {
	c := make(chan *skogul.Container)
	for {
		fo.workers <- c
		con := <-c
		fo.Next.Send(con)
	}
}
