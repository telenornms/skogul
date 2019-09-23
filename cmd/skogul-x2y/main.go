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

A more advanced example of how to use Skogul is provided in the package
documentation, where multiple receivers, multiple senders, failover and
more is covered.

What you are mainly missing with skogul-x2y is advanced error-handling,
load balancing, graceful failure, etc.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/config"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"github.com/KristianLyng/skogul/transformer"
)

var flisten = flag.String("receiver", "http://[::1]:8080", "Where to receive data from. See -help for details.")
var frecvhelp = flag.String("receiver-help", "", "Print extra options for receiver")
var fsendhelp = flag.String("sender-help", "", "Print extra options for sender")
var ftarget = flag.String("sender", "debug://", "Where to send data. See -help for details.")
var fhelp = flag.Bool("help", false, "Print extensive help/usage")
var fbatch = flag.Int("batch", 0, "Number of messages to batch up before passing them on as a single entity.")
var fcount = flag.String("count", "", "Print periodic stats using the count sender in addition to regular sender - same syntax as -sender (tip: -count debug://)")
var ferr = flag.String("errors", "null://", "Sender to divert errors to.")
var fretries = flag.Uint64("retries", 0, "Number of retries before giving up")
var fretrydelay = flag.Duration("retry-delay", time.Second, "Initial delay between retries. Retry time is exponential, so first retry iis retry-delay, second is 2xretry-delay, etc")

// Max width of help text before wrapping, should be some number lower than
// expected terminal size. 66 is nice for 80x25 terminals.
const helpWidth = 66

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
	for s, m := range sender.Auto {
		prettyPrint(s, m.Help)
	}
	fmt.Printf("\n\n")
	prettyHeader("receivers")
	for s, m := range receiver.Auto {
		prettyPrint(s, m.Help)
	}
}

func helpReceiver(s string) {
	if receiver.Auto[s] == nil {
		fmt.Printf("No such receiver %s\n", s)
		return
	}
	m := receiver.Auto[s]
	fmt.Printf("%s: %s\n\n", s, m.Help)
	fmt.Printf("Additional URL parameters for %s:\n\n", s)
	if m.Flags != nil {
		x := m.Flags()
		x.VisitAll(func(f *flag.Flag) {
			fmt.Printf("%s [%s]: %s\n", f.Name, f.DefValue, f.Usage)
		})
	}
}
func helpSender(s string) {
	if sender.Auto[s] == nil {
		fmt.Printf("No such sender %s\n", s)
		return
	}
	sh, _ := config.HelpSender(s)
	sh.Print()
}

func main() {
	flag.Parse()
	if *fhelp {
		help()
		os.Exit(0)
	}
	if *frecvhelp != "" {
		helpReceiver(*frecvhelp)
		os.Exit(0)
	}
	if *fsendhelp != "" {
		helpSender(*fsendhelp)
		os.Exit(0)
	}

	s, err := sender.New(*ftarget)
	if err != nil {
		log.Fatal(err)
	}

	if *fretries > 0 {
		ret := sender.Backoff{Next: s, Retries: *fretries, Base: *fretrydelay}
		s = &ret
	}
	errors := sender.ErrDiverter{}
	errors.Next = s
	errors.Err, err = sender.New(*ferr)
	if err != nil {
		log.Fatal("Failed to create sender where errors are diverted")
	}

	s = &errors

	if *fbatch > 0 {
		b := sender.Batch{Threshold: *fbatch}
		b.Next(s)
		s = &b
	}

	if *fcount != "" {
		cs, cerr := sender.New(*fcount)
		if cerr != nil {
			log.Fatal("Count sender failed to be created", cerr)
		}
		ch := skogul.Handler{Sender: cs}
		c := sender.Counter{Next: s, Stats: ch}
		s = &c
	}

	h := skogul.Handler{
		Parser:       parser.JSON{},
		Sender:       s,
		Transformers: []skogul.Transformer{transformer.Templater{}}}

	receiver, err := receiver.New(*flisten, h)
	if err != nil {
		log.Fatal(err)
	}

	receiver.Start()
}
