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

It also doesn't necessarily offer the full capabilities of the relevant
senders and receivers, but authors of senders and receivers are encouraged
to make it possible to expose as many features as possible in this fashion,
through the senders.Auto and receivers.Auto mechanisms.

What you are mainly missing with this package is advanced error-handling,
load balancing, graceful failure, etc.
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
	"os"
)

var flisten = flag.String("listen", "http://[::1]:8080", "Where to listen. See -help for details.")
var ftarget = flag.String("target", "debug://", "Target address. See -help for details.")
var fhelp = flag.Bool("help", false, "Print help")

func help() {
	fmt.Printf("Available senders:\n")
	fmt.Printf("%9s:// | %s\n", "scheme", "Description")
	fmt.Printf("%9s----+------------\n", "---------")
	for _, m := range senders.Auto {
		fmt.Printf("%9s:// | %s\n", m.Scheme, m.Help)
	}
	fmt.Printf("\n\n")
	fmt.Printf("Available receivers:\n")
	fmt.Printf("%9s:// | %s\n", "scheme", "Description")
	fmt.Printf("%9s----+------------\n", "---------")
	for _, m := range receivers.Auto {
		fmt.Printf("%9s:// | %s\n", m.Scheme, m.Help)
	}
}

func main() {
	flag.Parse()
	if *fhelp {
		help()
		os.Exit(0)
	}
	turl, err := url.Parse(*ftarget)
	if err != nil {
		log.Print("Failed to parse target url: %v", err)
		return
	}
	rurl, err := url.Parse(*flisten)
	if err != nil {
		log.Print("Failed to parse receiver url: %v", err)
		return
	}

	if senders.Auto[turl.Scheme] == nil {
		log.Fatalf("Unknown target scheme: %s", turl.Scheme)
	}

	target := senders.Auto[turl.Scheme].Init(*turl)

	h := skogul.Handler{
		Parser:       parsers.JSON{},
		Sender:       target,
		Transformers: []skogul.Transformer{transformers.Templater{}}}

	if receivers.Auto[rurl.Scheme] == nil {
		log.Fatalf("Unknown receiver scheme: %s", rurl.Scheme)
	}

	receiver := receivers.Auto[rurl.Scheme].Init(*rurl, h)

	receiver.Start()
}
