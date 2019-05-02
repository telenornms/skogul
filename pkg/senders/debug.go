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

package senders

import (
	"encoding/json"
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"math/rand"
	"net/url"
	"time"
)

/*
Debug sender simply prints the metrics in json-marshalled format to
stdout.
*/
type Debug struct {
}

func init() {
	addAutoSender("debug", NewDebug, "Debug sender prints received metrics to stdout")
}

/*
NewDebug creates a new Debug sender, ignoring the URL.
*/
func NewDebug(url url.URL) skogul.Sender {
	x := Debug{}
	return x
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
	Verbose  bool
}

// Send sleeps a random duration according to Sleeper spec, then passes the
// data to the next sender.
func (sl *Sleeper) Send(c *skogul.Container) error {
	d := rand.Float64() * float64(sl.MaxDelay)
	if sl.Verbose {
		log.Printf("Sleeping for %v", time.Duration(d))
	}
	time.Sleep(time.Duration(d))
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

faf := senders.ForwardAndFail{Next: skogul.Debug{}}
fb := senders.Fallback{Next: []skogul.Sender{influx, faf}}

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
