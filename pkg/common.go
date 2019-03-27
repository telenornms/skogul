/*
 * skogul, common trivialities
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

/*
Skogul is a framework for receiving, processing and forwarding data,
typically metric data or event-oriented data, at high throughput.

It is designed to be as agnostic as possible with regards to how it
transmits data and how it receives it, and the processors in between
need not worry about how the data got there or how it will be treated in
the next chain.

This means you can use Skogul to receive data on a influxdb-like
line-based TCP interface and send it on to postgres - or influxdb -
without having to write explicit support, just set up the chain.

The guiding principles of Skogul is:

- Make as few assumptions as possible about how data is received
- Be stupid fast

In the most simple setup, you can use Skogul simply to receive data from
a random shell script and send it to influxdb. In a more complex setup,
you can have multiple Skogul servers, each in different security zones
receiving subsets of a total data set, write it to a local queue, then
transmit - through strong authentication - to two central Skogul servers
that store the data to multiple influxdb instances based on sharding
rules.

A full blown example of Skogul is provided in cmd/skogul/main.go, which
is - as the file itself points out - meant as an EXAMPLE, not a final
solution.
*/
package skogul

import (
	"log"
)

/*
Handler determines what a receiver will do with data received. It requires a parser
to interperet the raw data, 0 or more transformers to mutate Containers, and a sender
to call after data is parsed and mutated and ready to be dealt with.
*/
type Handler struct {
	Parser       Parser
	Transformers []Transformer
	Sender       Sender
}

// Parser is the interface for parsing arbitrary data into a Container
type Parser interface {
	Parse(data []byte) (Container, error)
}

/*
Sender accepts data through Send() - and "sends it off". The canoncial
sender is one that implements a storage backend or outgoing API. E.g.:
accept data, send to influx.

Senders are not allowed to modify the Container - there could be multiple
goroutines running with same Container. If modification is required, the
Sender needs to take a copy.

While a single sender writing to a database is useful, the true power of
the Sender-pattern is chaining multiple, tiny, senders together to build
completely custom "sender chains" and thus provide site-specific handling
of data.

E.g.: The fallback sender is set up using a list of "down stream"
senders and will accept data and try the first sender on the list,
if that fail, try the next, and so on.
*/
type Sender interface {
	Send(c *Container) error
}

/*
Transformer mutates a collection before it is passed to a sender. Transformers
should be very fast, but are the only means to modifying the data.

Currently, the only transformer is the template transfromer, that expands a
template in a Container so underlying senders need not worry about the existence
of a template or not.
*/
type Transformer interface {
	Transform(c *Container) error
}

/*
Receiver is how we get data. The only current implementation is a HTTP
interface, but we should also expect UDP-receivers, line-based
TCP-receivers and even things such as influxdb-format receivers.

Currently, the included test-tool of skogul does not use a receiver, but
in the future there will probably be a "synthesizer"-receiver both to
demonstrate the concept and because it simplifies the tester.
*/
type Receiver interface {
	Start() error
}

/*
Not sure we really need these, but here you are...
*/
type Error struct {
	Reason string
}

func (e Error) Error() string {
	log.Printf("Error: %v", e.Reason)
	return e.Reason
}
