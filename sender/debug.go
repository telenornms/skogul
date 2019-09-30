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

package sender

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/KristianLyng/skogul"
)

/*
Debug sender simply prints the metrics in json-marshalled format to
stdout.
*/
type Debug struct {
	Prefix string `doc:"Prefix to print before any metric"`
}

// Send prints the JSON-formatted container to stdout
func (db *Debug) Send(c *skogul.Container) error {
	b, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		log.Panicf("Unable to marshal json for debug output: %s", err)
		return err
	}
	fmt.Printf("%s%s\n", db.Prefix, b)
	return nil
}

/*
The Sleeper sender injects a random delay between Base and Base+MaxDelay
before passing execution over to the Next sender.

The purpose is testing.
*/
type Sleeper struct {
	Next     skogul.SenderRef `doc:"Sender that will receive delayed metrics"`
	MaxDelay skogul.Duration  `doc:"The maximum delay we will suffer"`
	Base     skogul.Duration  `doc:"The baseline - or minimum - delay"`
	Verbose  bool             `doc:"If set to true, will log delay durations"`
}

// Send sleeps a random duration according to Sleeper spec, then passes the
// data to the next sender.
func (sl *Sleeper) Send(c *skogul.Container) error {
	d := sl.Base.Duration + time.Duration(rand.Float64()*float64(sl.MaxDelay.Duration))
	if sl.Verbose {
		log.Printf("Sleeping for %v", d)
	}
	time.Sleep(d)
	return sl.Next.S.Send(c)
}

/*
ForwardAndFail sender will pass the container to the Next sender, but
always returns an error. The use-case for this is to allow the fallback
Sender or similar to eventually send data to a sender that ALWAYS works,
e.g. the Debug-sender og just printing a message in the log, but we still
want to propagate the error upwards in the stack so clients can take
appropriate action.

Example use:

faf := sender.ForwardAndFail{Next: skogul.Debug{}}
fb := sender.Fallback{Next: []skogul.Sender{influx, faf}}

*/
type ForwardAndFail struct {
	Next skogul.SenderRef `doc:"Sender receiving the metrics"`
}

// Send forwards the data to the next sender and always returns an error.
func (faf *ForwardAndFail) Send(c *skogul.Container) error {
	err := faf.Next.S.Send(c)
	if err == nil {
		return skogul.Error{Reason: "Forced failure"}
	}
	return err
}

// ErrDiverter calls the Next sender, but if it fails, it will convert the
// error to a Container and send that to Err.
type ErrDiverter struct {
	Next   skogul.SenderRef `doc:"Send normal metrics here"`
	Err    skogul.SenderRef `doc:"If the sender under Next fails, convert the error to a metric and send it here"`
	RetErr bool             `doc:"If true, the original error from Next will be returned, if false, both Next AND Err has to fail for Send to return an error."`
}

// Send data to the next sender. If it fails, use the Err sender.
func (ed *ErrDiverter) Send(c *skogul.Container) error {
	err := ed.Next.S.Send(c)
	if err == nil {
		return nil
	}
	cerr, ok := err.(skogul.Error)
	if !ok {
		cerr = skogul.Error{Source: "errdiverter sender", Reason: "downstream error", Next: err}
	}
	container := cerr.Container()
	newerr := ed.Err.S.Send(&container)
	if newerr != nil {
		return newerr
	}
	if ed.RetErr {
		return err
	}
	return nil
}

// Null sender does nothing and returns nil - mainly for test-purposes
type Null struct{}

// Send just returns nil
func (n *Null) Send(c *skogul.Container) error {
	return nil
}

// Test sender is used to facilitate tests, and discards any metrics, but
// increments the Received counter.
type Test struct {
	received uint64
}

// Send discards data and increments the Received counter
func (rcv *Test) Send(c *skogul.Container) error {
	atomic.AddUint64(&rcv.received, 1)
	return nil
}

type failer interface {
	Errorf(string, ...interface{})
	Helper()
}

// TestTime sends the container on the specified sender, waits delay period
// of time, then verifies that rcv has received the expected number of
// containers.
func (rcv *Test) TestTime(t failer, s skogul.Sender, c *skogul.Container, received uint64, delay time.Duration) {
	t.Helper()
	rcv.Set(0)
	err := s.Send(c)
	if err != nil {
		t.Errorf("sending on %v failed: %v", s, err)
	}
	time.Sleep(delay)

	r := rcv.Received()
	if r != received {
		t.Errorf("sending on %v: wanted %d received, got %d", s, received, r)
	}
}

// Set atomicly sets the received counter to v
func (rcv *Test) Set(v uint64) {
	atomic.StoreUint64(&rcv.received, v)
}

// Received returns the amount of containers received
func (rcv *Test) Received() uint64 {
	return atomic.LoadUint64(&rcv.received)
}

// TestNegative sends data on s and expects to fail.
func (rcv *Test) TestNegative(t failer, s skogul.Sender, c *skogul.Container) {
	t.Helper()
	rcv.Set(0)
	err := s.Send(c)
	if err == nil {
		t.Errorf("sending on %v expected to fail, but didn't.", s)
	}
}

// TestQuick sends data on the sender and waits 5 milliseconds before
// checking that the data was received on the other end.
func (rcv *Test) TestQuick(t failer, sender skogul.Sender, c *skogul.Container, received uint64) {
	t.Helper()
	rcv.TestTime(t, sender, c, received, time.Duration(5*time.Millisecond))
}
