/*
 * skogul, examples
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

package skogul_test

import (
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"github.com/KristianLyng/skogul/transformer"
	"time"
)

// Example of the simplest Skogul chain possible
func Example() {
	// Create a debug-sender. A Debug-sender just prints the metric to
	// stdout.
	s := sender.Debug{}

	// A handler is used to inform a receiver how to treat incoming
	// data. This one will parse it using the JSON parser, then send it
	// on to the above sender.
	h := skogul.Handler{Parser: parser.JSON{}, Sender: s}

	// Create a receiver. The receiver.New() will parse a URL to find
	// an underlying receiver that implements the schema. In this case,
	// it will use the HTTP receiver.
	r, err := receiver.New("http://localhost:1234", h)
	if err != nil {
		panic(err)
	}

	// Finally, start the receiver.
	r.Start()
}

/*
A complex skogul example, using multiple storage services, internal
accounting, fallback methods and more. This still just demonstrate
one way to implement a Skogul chain. Most trivial use-cases are covered
by skogul-x2y, see http://godoc.org/github.com/KristianLyng/skogul/cmd/skogul-x2y

While it is usually possible to write it "the right way up" - start with
where we receive data and add senders - it's easier to write it "up-side-down"
and start with the final resting place for data. So you might consider reading
this bottom-up if it looks weird.

We are going to set up a chain of senders, starting with HTTP receiver which has
three distinct paths and handlers/chains: / sends to the primary handler,
/debug sends to a debug handler that just echos the parsed JSON to log.Print, and
/debug/notemplate does the same, but does not expand any provided template.

The primary chain looks sort of like this:

	http -> fallback -> detach -> fanout -> delay -> counter ->  dupe2 -> postgres
	            \                                        \             `-> influx
	             \			                      `------------/
		      \
		       `- dupe --> Log(print "the following failed")
		               `-> debug(print json to stdout)
*/
func Example_complex() {
	// Let's start by setting up two "final" storage senders
	influx := &sender.InfluxDB{URL: "http://127.0.0.1:8086/write?db=test", Measurement: "test"}
	postgres := &sender.Postgres{ConnStr: "user=postgres dbname=test host=localhost port=5432 sslmode=disable"}
	// Init is optional, see the skogul.senders.Postgres documentation
	postgres.Init()

	// Set up a duplicator and hook influx and postgres up to it -
	// Everything going to the duplicator will go to both influx and
	// postgres.
	dupe2 := sender.Dupe{Next: []skogul.Sender{influx, postgres}}

	// Set up a handler for where to send statistics. In this case, we
	// just send it to influx.
	countHandler := skogul.Handler{
		Sender:       influx,
		Transformers: []skogul.Transformer{}}

	// The counter generates statistics for us every Period time
	// (assuming data) and sends it to the Stats-handler (here:
	// influx). While it might seem strange to have a handler instead
	// of just a Sender at first, this allows us to provide arbitrary
	// transformers to the stats, e.g.: add metadata.
	counter := &sender.Counter{Next: dupe2, Stats: countHandler, Period: 1 * time.Second}

	// Let's also inject a random delay for testing!
	delay := sender.Sleeper{Next: counter, MaxDelay: 5000 * time.Millisecond, Verbose: false}

	fanout := sender.Fanout{Next: &delay}

	// Let's detach
	detach := sender.Detacher{Next: &fanout}

	// An other duplicator. This one just prints "The following failed"
	// and then uses the Debug-sender to print the metrics.
	dupe := sender.Dupe{Next: []skogul.Sender{sender.Log{Message: "The following failed"}, sender.Debug{}}}

	// the Fallback sender tries to write to the delay-sender
	// (delay->counter->dupe2->{postgres,influx}), but if this
	// fails, it will write to the dupe-sender (print "the following
	// failed" and the request).
	fb := sender.Fallback{}
	fb.Add(&detach)
	fb.Add(&dupe)

	// That takes care of the sender-chains. Let's set up three
	// receiver handlers.

	// This is the "normal" one - send to the fallback sender and
	// that's it. It also has a single transformer that - prior to
	// sending the data on - expands any template provided.
	h := skogul.Handler{
		Parser:       parser.JSON{},
		Sender:       &fb,
		Transformers: []skogul.Transformer{transformer.Templater{}}}

	// This is the same - but just print the request.
	debugtemplate := skogul.Handler{
		Parser:       parser.JSON{},
		Sender:       sender.Debug{},
		Transformers: []skogul.Transformer{transformer.Templater{}}}

	// Print the request, but do NOT expand the template. Demonstrates
	// what a template does and what the template transformer does.
	debugnotemplate := skogul.Handler{
		Parser:       parser.JSON{},
		Sender:       sender.Debug{},
		Transformers: []skogul.Transformer{}}

	// Set up a HTTP receiver
	rcvr := receiver.HTTP{Address: "[::1]:8080"}

	// Add the various handlers to relevant paths.
	rcvr.Handle("/", &h)
	rcvr.Handle("/debug", &debugtemplate)
	rcvr.Handle("/debug/notemplate", &debugnotemplate)

	// Start it
	rcvr.Start()
}
