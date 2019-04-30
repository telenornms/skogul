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
Package senders is a set of types that implement skogul.Sender. A Sender in
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

*/
package senders

import (
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"net/url"
)

// AutoSender is used to provide generic constructors by URL/Scheme.
// See Auto or cmd/skogul-x2y for more.
type AutoSender struct {
	Scheme string
	Init   func(url url.URL) skogul.Sender
	Help   string
}

// Auto maps schemas to senders and help text to make appropriate senders.
var Auto map[string]*AutoSender

func addAutoSender(scheme string, init func(url url.URL) skogul.Sender, help string) {
	if Auto == nil {
		Auto = make(map[string]*AutoSender)
	}
	if Auto[scheme] != nil {
		log.Fatalf("BUG: Attempting to overwrite existing auto-add sender %v", scheme)
	}
	Auto[scheme] = &AutoSender{scheme, init, help}
}
