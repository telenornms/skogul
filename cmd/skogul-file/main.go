/*
 * skogul, using config file
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
skogul-file parses a json-based config file and starts skogul.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/KristianLyng/skogul/config"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
)

var ffile = flag.String("file", "~/.config/skogul.json", "Path to skogul config to read.")
var fhelp = flag.Bool("help", false, "Print more help")
var frecvhelp = flag.String("receiver-help", "", "Print extra options for receiver")
var fsendhelp = flag.String("sender-help", "", "Print extra options for sender")
var fconf = flag.Bool("show", false, "Print the parsed JSON config instead of starting")

func help() {
	flag.Usage()
	fmt.Println(`
The config file needs to specify three items as a minimum: senders,
receivers and handlers, each with at least one item. Data is received
(or generated) in a receiver, which then uses a handler to first parse
the data and optionally transform it. At present, only JSON-parsing
is implemented, and only templating for transformation. This will change
on demand.

Each receiver/sender needs to specify a type, and to find documentation for
that type, use -receiver-help <type> or -sender-help <type>.

Semi-complete example:

{
  "senders": {
    "print": {
      "type": "debug"
    }
  },
  "receivers": {
    "api": {
      "type": "http",
      "handlers": {
	"/": "api_handler"
      }
    }
  },
  "handlers": {
    "api_handler": {
      "parser": "json",
      "transformers": ["templater"],
      "sender": "print"
    }
  }
}
`)
	fmt.Println("\nSenders:")
	for idx, sen := range sender.Auto {
		config.PrettyPrint(idx, sen.Help)
	}
	fmt.Println("\nReceivers:")
	for idx, rcv := range receiver.Auto {
		config.PrettyPrint(idx, rcv.Help)
	}
}

func main() {
	flag.Parse()
	if *fhelp {
		help()
		os.Exit(0)
	}
	if *fsendhelp != "" {
		sh, _ := config.HelpSender(*fsendhelp)
		sh.Print()
		os.Exit(0)
	}
	if *frecvhelp != "" {
		sh, _ := config.HelpReceiver(*frecvhelp)
		sh.Print()
		os.Exit(0)
	}

	c, err := config.File(*ffile)
	if err != nil {
		log.Fatal(err)
	}

	if *fconf {
		out, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
		os.Exit(0)
	}

	for _, r := range c.Receivers {
		r.Receiver.Start()
	}
}
