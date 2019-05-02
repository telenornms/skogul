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
skogul-x2y exposes a subset of Skogul's metric receivers and senders
in a single binary, supporting the most trivial use case of just moving
data from A to B. This will allow you to move data from Skogul's
HTTP API to influx, or read from file-mapped fifo and pass the data to
M&R's "port collector", and so on.

skogul-x2y only scratches the surface of what Skogul can do, but will probably
be sufficient for 95% of all deployments. Any sender or receiver in Skogul that
participate in the "auto"-scheme of configuration is supported implicitly. See
-help for an actual list.

A more advanced example of how to use Skogul is provided in cmd/skogul-demo,
where multiple receivers, multiple senders, failover and more is covered.

What you are mainly missing with skogul-x2y is advanced error-handling,
load balancing, graceful failure, etc.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/KristianLyng/skogul/pkg"
	"github.com/KristianLyng/skogul/pkg/parser"
	"github.com/KristianLyng/skogul/pkg/receiver"
	"github.com/KristianLyng/skogul/pkg/sender"
	"github.com/KristianLyng/skogul/pkg/transformer"
	"log"
	"net/url"
	"os"
	"strings"
)

var flisten = flag.String("receiver", "http://[::1]:8080", "Where to receive data from. See -help for details.")
var ftarget = flag.String("sender", "debug://", "Where to send data. See -help for details.")
var fhelp = flag.Bool("help", false, "Print extensive help/usage")

// Max width of help text before wrapping, should be some number lower than
// expected terminal size.
const helpWidth = 70

/*
Print a table of scheme | desc, wrapping the description at helpWidth.

E.g. assuming small helpWidth value:

Without prettyPrint:

foo:// | A very long line will be wrapped

With:

foo:// | A very long
       | line will
       | be wrapped

We wrap at word boundaries to avoid splitting words.
*/
func prettyPrint(scheme string, desc string) {
	fmt.Printf("%8s:// |", scheme)
	fields := strings.Fields(desc)
	l := 0
	for _, w := range fields {
		if (l + len(w)) > helpWidth {
			l = 0
			fmt.Printf("\n%11s |", "")
		}
		fmt.Printf(" %s", w)
		l += len(w) + 1
	}
	fmt.Printf("\n")
}

// Convenience function to avoid copy/paste
func prettyHeader(title string) {
	fmt.Printf("Available %s:\n", title)
	fmt.Printf("%8s:// | %s\n", "scheme", "Description")
	fmt.Printf("%8s----+------------\n", "--------")
}

func help() {
	flag.Usage()
	fmt.Printf("\n")
	fmt.Print("skogul-x2y sets up a skogul receiver, accepts data from it and passes it to the sender.")
	fmt.Printf("\n\n")
	prettyHeader("senders")
	for _, m := range sender.Auto {
		prettyPrint(m.Scheme, m.Help)
	}
	fmt.Printf("\n\n")
	prettyHeader("receivers")
	for _, m := range receiver.Auto {
		prettyPrint(m.Scheme, m.Help)
	}
}

func getUrls() (turl *url.URL, rurl *url.URL) {
	var err error
	turl, err = url.Parse(*ftarget)
	if err != nil {
		log.Printf("Failed to parse target url: %v", err)
		os.Exit(1)
	}
	rurl, err = url.Parse(*flisten)
	if err != nil {
		log.Printf("Failed to parse receiver url: %v", err)
		os.Exit(1)
	}
	return
}

func main() {
	flag.Parse()
	if *fhelp {
		help()
		os.Exit(0)
	}

	turl, rurl := getUrls()

	if sender.Auto[turl.Scheme] == nil {
		log.Fatalf("Unknown target scheme: %s", turl.Scheme)
	}

	target := sender.Auto[turl.Scheme].Init(*turl)

	h := skogul.Handler{
		Parser:       parser.JSON{},
		Sender:       target,
		Transformers: []skogul.Transformer{transformer.Templater{}}}

	if receiver.Auto[rurl.Scheme] == nil {
		log.Fatalf("Unknown receiver scheme: %s", rurl.Scheme)
	}

	receiver := receiver.Auto[rurl.Scheme].Init(*rurl, h)

	receiver.Start()
}
