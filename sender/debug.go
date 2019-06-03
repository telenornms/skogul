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
	"log"
	"math/rand"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/KristianLyng/skogul"
)

/*
Debug sender simply prints the metrics in json-marshalled format to
stdout.
*/
type Debug struct {
}

func init() {
	addAutoSender("debug", newDebug, "Debug sender prints received metrics to stdout")
	addAutoSender("null", newNull, "Null discards the data")
}

/*
newDebug creates a new Debug sender, ignoring the URL.
*/
func newDebug(url url.URL) skogul.Sender {
	x := Debug{}
	return x
}

// newNull creates a Null sender that discards data
func newNull(url url.URL) skogul.Sender {
	x := Null{}
	return &x
}

// Send prints the JSON-formatted container to stdout
func (db Debug) Send(c *skogul.Container) error {
	b, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		log.Panicf("Unable to marshal json for debug output: %s", err)
		return err
	}
	log.Printf("Debug: \n%s", b)
	return nil
}

/*
The Sleeper sender injects a random delay between 0 and MaxDelay before
passing execution over to the Next sender.

The purpose is testing.
*/
type Sleeper struct {
	Next     skogul.Sender
	MaxDelay time.Duration
	Base     time.Duration
	Verbose  bool
}

// Send sleeps a random duration according to Sleeper spec, then passes the
// data to the next sender.
func (sl *Sleeper) Send(c *skogul.Container) error {
	d := sl.Base + time.Duration(rand.Float64()*float64(sl.MaxDelay))
	if sl.Verbose {
		log.Printf("Sleeping for %v", d)
	}
	time.Sleep(d)
	return sl.Next.Send(c)
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
	Next skogul.Sender
}

// Send forwards the data to the next sender and always returns an error.
func (faf *ForwardAndFail) Send(c *skogul.Container) error {
	err := faf.Next.Send(c)
	if err == nil {
		return skogul.Error{Reason: "Forced failure"}
	}
	return err
}

// ErrDiverter calls the Next sender, but if it fails, it will convert the
// error to a Container and send that to Err.
type ErrDiverter struct {
	// The ordinary sender to use
	Next skogul.Sender
	// If Next.Send() fails, create a container from the error and send
	// it to Err
	Err skogul.Sender
	// If true, the original error from Next will be returned, if false
	// both Next AND Err has to fail for Send to return an error.
	RetErr bool
}

func (ed *ErrDiverter) Send(c *skogul.Container) error {
	err := ed.Next.Send(c)
	if err == nil {
		return nil
	}
	cerr, ok := err.(skogul.Error)
	if !ok {
		cerr = skogul.Error{Source: "errdiverter sender", Reason: "downstream error", Next: err}
	}
	container := cerr.Container()
	newerr := ed.Err.Send(&container)
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
