/*
 * skogul, documentation
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

End users should only need to worry about the cmd/skogul tool, which comes
fully equipped with self-contained documentation.

Adding new logic to Skogul should also be fairly easy. New developers should
focus on understanding two things:

1. The skogul.Container data structure - which is the heart of Skogul.
2. The relationship from receiver to handler to sender.

The Container is documented in this very package.

Receivers are where data originates within Skogul. The typical Receiver
will receive data from the outside world, e.g. by other tools posting
data to a HTTP endpoint. Receivers can also be used to "create" data,
either test data or, for example, log data. When skogul starts, it
will start all receivers that are configured.

Handlers determine what is done with the data once received. They are
responsible for parsing raw data and optionally transform it. This is the
only place where it is allowed to _modify_ data. Today, the only transformer
is the "templater", which allows a collection of metrics which share certain
attributes (e.g.: all collected at the same time and from the same machine)
to provide these shared attributes in a template which the "templater"
transformer then applies to all metrics.

Other examples of transformations that make sense are:

- Adding a metadata field
- Removing a metadata field
- Removing all but a specific set of fields
- Converting nested metrics to multiple metrics or flatten them

Once a handler has done its deed, it sends the Container to the sender,
and this is where "the fun begins" so to speak.

Senders consist of just a data structure that implements the Send()
interface. They are not allowed to change the container, but besides that,
they can do "whatever". The most obvious example is to send the container
to a suitable storage system - e.g., a time series database.

So if you want to add support for a new time series database in Skogul, you
will write a sender.

In addition to that, many senders serve only to add internal logic and pass
data on to other senders. Each sender should only do one specific thing.
For example, if you want to write data both to InfluxDB and MySQL, you need
three senders: The "MySQL" and "InfluxDB" senders, and the "dupe" sender,
which just takes a list of other senders and sends whatever it receives
on to all of them.

Today, Senders and Receivers both have an identical "Auto"-system, found in
auto.go of the relevant directories. This is how the individual
implementations are made discoverable to the configuration system, and how
documentation is provided. Documentation for the settings of a sender/receiver
is handled as struct tags.

Once more parsers/transformers are added, they will likely also use a similar
system.
*/
package skogul
