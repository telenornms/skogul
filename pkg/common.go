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
	"log"
)

/*
 * The idea is that setting up new binaries should be trivial. It may be
 * necessary to write some Go-code, but given that, the fundamental
 * distinction is that receiving data should be disconnected from storing
 * it. E.g.: We might be collecting interface statistics with a home-brew
 * SNMP collector that can post to us, and a different collector based on
 * streaming telemetry instead of SNMP, and a third source might be a
 * different security domain where they might even use the TICK-stack and
 * export to use. So we have a generic receiver-stack (so far, only HTTP),
 * a reasonably generic writer/sender stack, but the trick is to be able to
 * put this all together.
 *
 * In all likelihood, these interfaces will change over time as we add more
 * actual use cases, but the idea is the above.
 *
 * Things that WILL come in to play: What parts are concurrent and what are
 * not? Queues? What are idempotent, what are not?
 */

/*
 * The Handler is inteded to be the "you've got the data... now what?"
 * part. It most certainly will not look like this at the end of the day.
 */
type Handler struct {
	Transformers []Transformer
	Senders      []Sender
}

/*
 * A Sender is _seemingly_ the final step in the chain. It accepts data,
 * and that's the end of it. Anything that writes to storage will be a
 * sender, but it can also be things such as queues (e.g.: Send() to a
 * queue which only ensures a quick accept, then sends it on to a slower
 * storage up to a threshold, after which it starts blocking or dropping,
 * depending on the queue-setup). This will allow us to do things like have
 * a queue for streaming telemetry that is lossy, while the queue for
 * SNMP-data is not.
 */
type Sender interface {
	Send(c *Container) error
}

/*
 * A transformer is a fast(!) way to modify a collection of metrics. Not
 * sure exactly how useful it will be, but it could convecivably be things
 * such as ASN1-lookup, lower-casing(?) of names, or the whole
 * template-logic.
 *
 * The key difference between a transformer and a sender is that a
 * transformer just modifies an existing collection, while a sender
 * consumes it and (optionally) produces others if further processing is
 * needed.
 */
type Transformer interface {
	Transform(c *Container) error
}

/*
 * Receiver is how we get data. The only current implementation is a HTTP
 * interface, but we should also expect UDP-receivers, line-based
 * TCP-receivers and even things such as influxdb-format receivers.
 *
 * The exact details of the interface will most likely change once we see
 * more how it will be used in real deployments.
 */
type Receiver interface {
	Start() error
}

/*
 * Not sure we really need these, but here you are...
 */
type Gerror struct {
	Reason string
}

func (e Gerror) Error() string {
	log.Printf("Error: %v", e.Reason)
	return e.Reason
}
