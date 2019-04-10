/*
 * skogul, http receiver, influxdb writer, example
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
Skogul is primarily a framework, where you use it to build your OWN
binaries. This package is provided to satisfy the common deployment
of using Skogul to push to MnR.
*/
package main

import (
	"flag"
	"github.com/KristianLyng/skogul/pkg"
	"github.com/KristianLyng/skogul/pkg/parsers"
	"github.com/KristianLyng/skogul/pkg/receivers"
	"github.com/KristianLyng/skogul/pkg/senders"
	"github.com/KristianLyng/skogul/pkg/transformers"
)

var flisten = flag.String("listen", ":8080", "Address Skogul will listen for HTTP on")
var fmnr = flag.String("mnr", "127.0.0.1:1234", "Adress of M&R collector to send to")
var fdefaultGroup = flag.String("defaultgroup", "group", "Default storage group of MnR")

func main() {
	flag.Parse()
	mnr := &senders.MnR{Address: *fmnr, DefaultGroup: *fdefaultGroup}

	h := skogul.Handler{
		Parser:       parsers.JSON{},
		Sender:       mnr,
		Transformers: []skogul.Transformer{transformers.Templater{}}}

	receiver := receivers.HTTP{Address: *flisten}
	receiver.Handle("/", &h)
	receiver.Start()
}
