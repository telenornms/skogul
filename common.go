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
Package skogul is a framework for receiving, processing and forwarding
data, typically metric data or event-oriented data, at high throughput.

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
rules. Etc.

A general-purpose command is provided under cmd/skogul-x2y, which provides
a mechanism for starting a simplistic 1-to-1 receiver/sender chain, and
should satisfy many common use-cases. E.g.: Accepting data from MQTT and
passing it to influx. Or accepting from HTTP and printing to stdout for
debug purposes.

A more complete example is also provided further down.
*/
package skogul

import (
	"flag"
	"fmt"
	"net/url"
)

/*
Handler determines what a receiver will do with data received. It requires
a parser to interperet the raw data, 0 or more transformers to mutate
Containers, and a sender to call after data is parsed and mutated and ready
to be dealt with.
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
Sender accepts data through Send() - and "sends it off". The canonical
sender is one that implements a storage backend or outgoing API. E.g.:
accept data, send to influx.

Senders are not allowed to modify the Container - there could be multiple
goroutines running with same Container. If modification is required, the
Sender needs to take a copy.

A sender should assume that the container has been validated, and is
non-null. Slightly counter to common sense, it is NOT recommended to
verify the input data again, since multiple senders are likely chained
and will thus likely redo the same verifications.

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

Currently, the only transformer is the template transformer, that expands a
template in a Container so underlying senders need not worry about the
existence of a template or not.
*/
type Transformer interface {
	Transform(c *Container) error
}

/*
Receiver is how we get data. At the point of this writing, a HTTP
receiver, line-based TCP receiver, line-based FIFO receiver, MQTT receiver
and "tester"-receiver exists. A receiver uses a Handler to deal with
data, including how to parse it. Meaning: The HTTP receiver could support
both Skogul JSON-data and other formats, and the same goes for any other
receiver.
*/
type Receiver interface {
	Start() error
}

/*
Error is a typical skogul error. All Skogul functions should provide Source
and Reason.

If the Next field is provided, error messages will recurse to the bottom, thus
propagating errors from the bottom and up.
*/
type Error struct {
	Reason string
	Source string
	Next   error
}

// Error for use in regular error messages. Also outputs to log.Print().
// Will also include e.Next, if present.
func (e Error) Error() string {
	src := "<nil>"
	if e.Source != "" {
		src = e.Source
	}
	tail := ""
	if e.Next != nil {
		tail = fmt.Sprint(": ", e.Next.Error())
	}
	return fmt.Sprintf("%s: %s%s", src, e.Reason, tail)
}

// URLParse parses a url's "GET parameters" into the provided FlagSet.
func URLParse(u url.URL, fs *flag.FlagSet) error {
	vs := u.Query()
	for i, v := range vs {
		for _, e := range v {
			err := fs.Set(i, e)
			if err != nil {
				return Error{Source: "auto receiver", Reason: fmt.Sprintf("failed to parse argument %s value %s", i, e), Next: err}
			}
		}
	}
	return nil
}
