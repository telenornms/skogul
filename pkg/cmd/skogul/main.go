/*
 * skogul, main method/init
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

package main

import (
	"github.com/KristianLyng/skogul/pkg"
	"github.com/KristianLyng/skogul/pkg/receivers"
	"github.com/KristianLyng/skogul/pkg/senders"
	"github.com/KristianLyng/skogul/pkg/transformers"
	"time"
)

/*
Please be aware: This is NOT inteded to be "The final perfect version"
of how to run Skogul.

Skogul is primarily a framework, where you use it to build your OWN
binaries.

This is meant more to show-case existing features. Right now it's merely
a "This is where I test everything" thing.

In the future, a few more examples will be provided in separate files.

It kinda reads a bit up-side down, starting with the final resting place
for data so to speak, then ending with starting a web server to receive
data.
*/
func main() {

	// Let's start by setting up two "final" storage senders
	influx := &senders.InfluxDB{URL: "http://127.0.0.1:8086/write?db=test", Measurement: "test"}
	postgres := &senders.Postgres{ConnStr: "user=postgres dbname=test host=localhost port=5432 sslmode=disable"}
	// Init is optional, see the skogul.senders.Postgres documentation
	postgres.Init()

	// Set up a duplicator and hook influx and postgres up to it -
	// Everything going to the duplicator will go to both influx and
	// postgres.
	dupe2 := senders.Dupe{Next: []skogul.Sender{influx, postgres}}

	// The counter generates statistics for us every Period time
	// (assuming data) and sends it to the Stats-sender (here: influx)
	counter := &senders.Counter{Next: dupe2, Stats: influx, Period: 1 * time.Second}

	// Let's also inject a random delay for testing
	delay := &senders.Sleeper{counter, 1 * time.Millisecond, false}

	// An other duplicator. This one just prints "The following failed"
	// and then uses the Debug-sender to print the metrics.
	dupe := senders.Dupe{Next: []skogul.Sender{senders.Log{"The following failed"}, senders.Debug{}}}

	// the Fallback sender tries to write to the delay-sender
	// (delay->counter->dupe2->{postgres,influx}), but if this
	// fails, it will write to the dupe-sender (print "the following
	// failed" and the request).
	fb := senders.Fallback{Next: []skogul.Sender{delay, dupe}}

	// That takes care of the sender-chains. Let's set up three
	// receiver handlers.

	// This is the "normal" one - send to the fallback sender and
	// that's it. It also has a single transformer that - prior to
	// sending the data on - expands any template provided.
	h := skogul.Handler{
		Senders:      []skogul.Sender{fb},
		Transformers: []skogul.Transformer{transformers.Templater{}}}

	// This is the same - but just print the request.
	debugtemplate := skogul.Handler{
		Senders:      []skogul.Sender{senders.Debug{}},
		Transformers: []skogul.Transformer{transformers.Templater{}}}

	// Print the request, but do NOT expand the template. Demonstrates
	// what a template does and what the template transformer does.
	debugnotemplate := skogul.Handler{
		Senders:      []skogul.Sender{senders.Debug{}},
		Transformers: []skogul.Transformer{}}

	// Set up a HTTP receiver
	receiver := receivers.HTTP{Address: "localhost:8080"}

	// Add the various handlers to relevant paths.
	receiver.Handle("/", &h)
	receiver.Handle("/debug", &debugtemplate)
	receiver.Handle("/debug/notemplate", &debugnotemplate)

	// Start it
	receiver.Start()
}
