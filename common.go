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

package skogul

import (
	"fmt"
	"time"
)

/*
Handler determines what a receiver will do with data received. It requires
a parser to interperet the raw data, 0 or more transformers to mutate
Containers, and a sender to call after data is parsed and mutated and ready
to be dealt with.

Whenever a new Container is created, it should pass that to a Handler, not
directly to a Sender. This goes for artificially created data too, e.g. if
a sender wants to emit statistics. This ensures that transformers can be
used in the future.

To make it configurable, a HandlerRef should be used.
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

Senders that pass data off to other senders should use a SenderRef instead,
to facilitate configuration.
*/
type Sender interface {
	Send(c *Container) error
}

/*
SenderVerification is an *optional* interface for senders. If a sender
implements it, the configuration engine will issue Verify() after all
configuration is parsed. The sender should never modify state upon
Verify(), but should simply check that internal state is usable.
*/
type SenderVerification interface {
	Verify() error
}

/*
Transformer mutates a collection before it is passed to a sender. Transformers
should be very fast, but are the only means to modifying the data.
*/
type Transformer interface {
	Transform(c *Container) error
}

/*
Receiver is how we get data. Receivers are responsible for getting raw data and the
outer boundaries of a Container, but should explicitly avoid parsing raw data.
This ensures that how data is transported is not bound by how it is parsed.
*/
type Receiver interface {
	Start() error
}

/*
Error is a typical skogul error. All Skogul functions should provide Source
and Reason. I'm not entirely sure why, except that it allows chaining errors?

If the Next field is provided, error messages will recurse to the bottom, thus
propagating errors from the bottom and up.
*/
type Error struct {
	Reason string
	Source string
	Next   error
}

/*
SenderRef is a reference to a named sender. This is required to allow
references to be resolved after all senders are loaded. Wherever a
Sender is loaded from configuration, a SenderRef should be used in its
place. The maintenance of the sender is handled in the configuration
system.
*/
type SenderRef struct {
	S    Sender
	Name string
}

// HandlerRef references a named handler. Used whenever a handler is
// defined by configuration.
type HandlerRef struct {
	H    *Handler
	Name string
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

// Container returns a skogul container representing the error
func (e Error) Container() Container {
	c := Container{}
	now := time.Now()
	c.Metrics = make([]*Metric, 1)
	m := Metric{}
	m.Metadata = make(map[string]interface{})
	m.Data = make(map[string]interface{})
	m.Time = &now
	m.Metadata["source"] = e.Source
	m.Data["reason"] = e.Reason
	m.Data["description"] = e.Error()
	c.Metrics[0] = &m
	return c
}

// Verify the basic integrity of a handler. Quite shallow.
func (h Handler) Verify() error {
	if h.Parser == nil {
		return Error{Reason: "Missing parser for Handler"}
	}
	if h.Transformers == nil {
		return Error{Reason: "Missing parser for Handler"}
	}
	for _, t := range h.Transformers {
		if t == nil {
			return Error{Reason: "nil-transformer for Handler"}
		}
	}
	if h.Sender == nil {
		return Error{Reason: "Missing parser for Handler"}
	}
	return nil
}
