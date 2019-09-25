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
cmd/skogul parses a json-based config file and starts skogul.
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
Skogul provides a set of receivers and senders of time-based data, and is
intended to bridge different systems together. In its simplest form, it will
read from, e.g. HTTP, and write to a database like InfluxDB. But it can also
be configured for considerably more complex chains, including failover, data
duplication, retries, batching of multiple data sets, and more.

It is designed to be very fast, and flexible.

To start skogul, you need to specify a chain. A chain starts with a
receiver - this is where data originates as far as Skogul is concerned.
Each receiver references one or more handler. A handler defines how
the raw data is parsed (currently only a JSON parser is supported),
any transformations of data (e.g.: filtering? Today, only templating
is supported) and the sender that will receive the data.

A sender accepts containers of metrics and Does Something With Them.
The most traditional senders will just write the data to a database.
Several databases are supported, including InfluxDB, MySQL, M&R and
more.

Other senders are meant to provide configuration options. For example,
you might want to write data to both MySQL and InfluxDB - this can be
achieved by sending to a "dupe" sender, which in turn will pass data
to both an InfluxDB-sender and a MySQL-sender. Other senders can be
used to batch data together before it is forwarded to storage, or
provide alternate paths if the preferred sender fails, etc. The list is
long.

A small example of a simple chain is provided here. See the individual
senders and receivers for more help, e.g. using -sender-help and
-receiver-help.

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
}`)
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
