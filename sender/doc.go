/*
 * skogul, sender boilerplate
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
Package sender is a set of types that implement skogul.Sender. A Sender in
skogul is a simple primitive that receives skogul metrics and "does
something with them".

The traditional and obvious sender accepts metrics and uses and external
service to persist them to disk. E.g.: the InfluxDB sender stores the
metrics to influxdb. The postgres sender accepts metrics and stores to
postgres, and so forth.

The other type of senders are "internal", typical for routing. The
classic examples are the "dupe" sender that accepts metrics and passes
them on to multiple other senders - e.g.: Store to both postgres and
influxdb. An other classic is the "fallback" sender: It has a list of
senders and tries each one in order until one succeeds, allowing you to
send to a primary influxdb normally - if influx fails, write to local disk,
if that fails, write a message to the log.

The only thing a sender "must" do is implement Send(c *skogul.Container),
and it is disallowed to modify the container in Send(), since multiple
senders might be working on it at the same time.

To make a sender configurable, simply ensure data types in the type
definition can be Unmarshalled from JSON. A small note on that is that
it is necessary to use "SenderRef" and "HandlerRef" objects instead of
Sender and Handler directly for now. This is to let the config engine
track references that haven't resolved yet.

It also means certain data types need to be avoided or worked around.
Currently, time.Duration is such an example, as it is missing a JSON
unmrashaller. For such data types, a simple wrapper will do the trick,
e.g. skogul.Duration wraps time.Duration.
*/
package sender
