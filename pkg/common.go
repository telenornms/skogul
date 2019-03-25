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
*/
package skogul

import (
	"log"
)

/*
The Handler is inteded to be the "you've got the data... now what?"
part. It most certainly will not look like this at the end of the day.
*/
type Handler struct {
	Transformers []Transformer
	Sender       Sender
}

/*
A Sender accepts data through Send() - and "sends it off". The canoncial
sender is one that implements a storage backend or outgoing API. E.g.:
accept data, send to influx.

But the real power of a sender is chaining them together using tiny,
single-purpose senders to build complicated logic.

E.g.: The fallback sender is set up using a list of "down stream"
senders and will accept data and try the first sender on the list,
if that fail, try the next, and so on.

A Sender should not modify the data it accepts. If it needs to do
that, it has to make a copy, as multiple senders may be accessing
the same data.
*/
type Sender interface {
	Send(c *Container) error
}

/*
A transformer is a fast(!) way to modify a collection of metrics. It
is the safe way to modify data before it is passed to the first Sender.
*/
type Transformer interface {
	Transform(c *Container) error
}

/*
Receiver is how we get data. The only current implementation is a HTTP
interface, but we should also expect UDP-receivers, line-based
TCP-receivers and even things such as influxdb-format receivers.

The exact details of the interface will most likely change once we see
more how it will be used in real deployments.
*/
type Receiver interface {
	Start() error
}

/*
Not sure we really need these, but here you are...
*/
type Gerror struct {
	Reason string
}

func (e Gerror) Error() string {
	log.Printf("Error: %v", e.Reason)
	return e.Reason
}
