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
of using Skogul as a simple gateway from X to Y - it only utilizes a
small subset of the possible senders and receivers provided by Skogul,
but should prove sufficient for many scenarios.

XXX: Future work on Senders and Receivers will include a convention/method
of automating this process. Quite possibly by providing a FromURL(url)
function which will (attempt to) generate the appropriate Sender
or Receiver based on URL.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/KristianLyng/skogul/pkg"
	"github.com/KristianLyng/skogul/pkg/parsers"
	"github.com/KristianLyng/skogul/pkg/receivers"
	"github.com/KristianLyng/skogul/pkg/senders"
	"github.com/KristianLyng/skogul/pkg/transformers"
	"log"
	"net/url"
)

var flisten = flag.String("listen", "http://[::1]:8080", "Where to listen. Supports schemas http:// mqtt:// and tcp://")
var ftarget = flag.String("target", "debug://", "Target address. Currently supported schemes: debug://, http://, influx://")

func main() {
	flag.Parse()
	var target skogul.Sender
	var receiver skogul.Receiver
	turl, err := url.Parse(*ftarget)
	if err != nil {
		log.Print("Failed to parse target url: %v", err)
		return
	}
	rurl, err := url.Parse(*flisten)
	if err != nil {
		log.Print("Failed to parse target url: %v", err)
		return
	}

	if turl.Scheme == "influx" {
		turl.Scheme = "http"
		target = &senders.InfluxDB{URL: turl.String()}
		log.Print("influx sender")
	} else if turl.Scheme == "http" {
		target = &senders.HTTP{URL: turl.String()}
		log.Print("HTTP sender")
	} else if turl.Scheme == "debug" {
		log.Print("debug sender")
		target = senders.Debug{}
	}

	h := skogul.Handler{
		Parser:       parsers.JSON{},
		Sender:       target,
		Transformers: []skogul.Transformer{transformers.Templater{}}}

	if rurl.Scheme == "http" {
		hl := &receivers.HTTP{Address: rurl.Host}
		hl.Handle(fmt.Sprintf("/%s", rurl.Path), &h)
		receiver = hl
	} else if rurl.Scheme == "mqtt" {
		receiver = &receivers.MQTT{Address: rurl.String(), Handler: &h}
	} else if rurl.Scheme == "tcp" {
		receiver = &receivers.TCPLine{Address: rurl.String(), Handler: &h}
	}

	receiver.Start()
}
