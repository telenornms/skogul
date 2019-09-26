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
	"sort"
	"strings"

	"github.com/KristianLyng/skogul/config"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
)

var ffile = flag.String("file", "~/.config/skogul.json", "Path to skogul config to read.")
var fhelp = flag.Bool("help", false, "Print more help")
var frecvhelp = flag.String("receiver-help", "", "Print extra options for receiver")
var fsendhelp = flag.String("sender-help", "", "Print extra options for sender")
var fconf = flag.Bool("show", false, "Print the parsed JSON config instead of starting")
var fman = flag.Bool("make-man", false, "Output RST documentation suited for rst2man")

func man() {
	fmt.Print(`
======
skogul
======

------
Skogul
------

:Manual section: 1
:Authors: Kristian Lyngstøl <kly@kly.no>

SYNOPSIS
========

::

	skogul -file config-file [-help] [-sender-help sender] 
	       [-receiver-help receiver] [-show] [-make-man]

DESCRIPTION
===========

Skogul is a framework for normalising and transporting event-oriented data
between different systems. It operates using a concept of a receiver, handler
and sender. A receiver will accept (or in special cases: generate) data, it
will then use a handler to prosess that data, e.g. parse JSON data. The handler
will forward the processed data to a sender.

Senders come in two distinct but interchangable variants: Storage-oriented
senders are used to send the data to some external resource, e.g., a time series
database like InfluxDB. Utility-oriented senders are used to do things like
route data two multiple senders (e.g.: store it in InfluxDB locally, but also
forward it to a remote database), handle failures by adding retry-mechanics or
fallback logic, and much more.

The simplest possible example is to receive data over HTTP, parse it as JSON and
store it to InfluxDB. A more valuable example is to receive data over insecure
HTTP locally, batch together a collection of data, then forward it over a
TLS-encrypted and authenticated HTTPS channel to an other skogul instance in
a different security domain, which can then store it to disk.

There are more examples in the the "examples/" directory.

Configuration is written as JSON as well, and dynamically parsed.

OPTIONS
=======

`)

	f := flag.CommandLine
	f.VisitAll(func(fl *flag.Flag) {
		s := fmt.Sprintf("``-%s``", fl.Name)
		name, usage := flag.UnquoteUsage(fl)
		if len(name) > 0 {
			s += " " + name
		}
		s += "\n\t"
		s += strings.ReplaceAll(usage, "\n", "\n\t")
		if fl.DefValue != "" {
			s += fmt.Sprintf(" (default %v)", fl.DefValue)
		}
		fmt.Print(s, "\n\n")
	})

	fmt.Print(`
CONFIGURATION
=============

Configuration of skogul is done with a json config file, referenced with
the -file option. You need to specify at least one receiver, handler and
sender to make something sensible.

The base configuration set is::

  {
    "receivers": {
      "xxx": {
        "type": "type-of-receiver",
        type-specific-options
      },
      "other-receiver...": ...
    },
    "handlers": {
      "yyy": {
        "parser": "json", // other options might come
        "transformers": [...], // only valid option today is [] or ["templater"]
        "sender": "reference-to-sender"
      }
    },
    "senders": {
      "zzz": {
        "type": "type-of-sender",
        type-specific-options
      },
      "qqq": {
        "type": "type-of-sender",
        type-specific-options
      },
      ...
    }
  }

In the above pseudo-config, "xxx", "yyy", "zzz" and "qqq" are arbitrary
names you chose that are how to reference that specific item within the same
configuration. The "type" field referneces what implementation to use - the list
of different sender-types and receiver-types are below. Each type of sender
and receiver require different options (e.g.: MySQL sender will require a
connection string, while InfluxDB sender will require a URL).

At present time, there is only a single parser and a single transformer, so
handlers mainly serve to name the next/initial sender for a receiver.

The documentation for each sender and receiver also lists all options. In
general, you do not need to specify all options. For formatting, the settings
use whatever JSON unmarshalling logic that Go provides, but it should be self
explanatory or explained in the documentation for the relevant option.

SENDERS
=======

Senders move parsed data around according to internal logic. A sender will typically
either ensure data is stored, or forward it according to some internal logic to an
other sender with that goal.

The following senders exist. A list can also be retrieved by using the "-help"
option, and -sender-help or -receiver-help.

`)
	senders := []string{}
	for idx := range sender.Auto {
		if sender.Auto[idx].Name != idx {
			continue // alias
		}
		senders = append(senders, idx)
	}
	sort.Strings(senders)
	for _, s := range senders {
		sh, _ := config.HelpSender(s)
		thingMan(sh)
	}
	fmt.Print(`
RECEIVERS
=========

Receivers accept data from the outside world - or in special cases,
generate the data themself. Receivers do not typically deal with how
individual collections of data is handled, but leaves that specific task
to a handler.

The following receivers exist.

`)
	receivers := []string{}
	for idx := range receiver.Auto {
		if receiver.Auto[idx].Name != idx {
			continue // alias
		}
		receivers = append(receivers, idx)
	}
	sort.Strings(receivers)
	for _, r := range receivers {
		sh, _ := config.HelpReceiver(r)
		thingMan(sh)
	}
	fmt.Print(`
HANDLERS
========

There is only one type of handler. It accepts three arguments: A parser to
parse data, a list of optional transformers, and the first sender that will
receive the parsed container(s).

Currently the only valid parser is "json" and the only valid transformer is
"templating".

FIXME: Templating

JSON FORMAT
===========

Data sent to Skogul will be parsed to fit the internal data model of Skogul. The
JSON representation is roughly thus::

  {
    "template": { 
      "timestamp": "iso8601-time",
      "metadata": { 
        "key": value, 
        ...
      },
      "data": {
        "key": value,
        ...
      }
    },
    "metrics": [
      {
        "timestamp": "iso8601-time",
        "metadata": { 
          "key": value, 
          ...
        },
        "data": {
          "key": value,
          ...
        }
      },
      { ...}
    ]
  }

The entire "template" is optional. If the "templater" transformer is
applied, all metrics will start with whatever value is present in the
template, and then overwrite with "local" variables. E.g.: If all your
metrics share timestamp in a collection, you can specify that in the
template. Or if they share some metadata.

The primary difference between metadata and data is searchability,
and it will depend on storage engines. Typically this means the name
of a server is metadata, but the load average is data. Skogul itself
does not much care.

EXAMPLES
========

The following specifies an insecure HTTP-based receiver that will wait up
to 5 seconds or 1000 metrics before writing data to InfluxDB::

  {
    "receivers": {
      "api": {
        "type": "http",
        "address": "[::1]:8080",
        "handlers": {
          "/": "jsontemplating"
        }
      }
    },
    "handlers": {
      "jsontemplating": {
        "parser": "json",
        "transformers": [ "templater" ],
        "sender": "batch"
      }
    },
    "senders": {
      "batch": {
        "type": "batch",
        "interval": "5s",
        "threshold": 1000,
        "next": "influx"
      },
      "influx": {
        "type": "influx",
        "URL": "http://[::1]:8086/write?db=testdb",
        "measurement": "demo",
        "Timeout": "10s"
      }
    }
  }

More examples are provided in the examples/ directory of the Skogul source
package.

SEE ALSO
========

https://github.com/KristianLyng/skogul

BUGS
====

The biggest known issue right now is that the configuration engine is a bit
horrible at giving constructive error message, and will silently ignore
unkown (or misspelled) variable names. Work in progress.

A tip for working around this is to compare your configuration with what
skogul outputs when you run it with -show, as that is a representation of
the parsed configuration.

COPYRIGHT
=========

This document is licensed under the same license as Skogul itself, which
happens to be GPLv2 (or later). See LICENSE for details.

* Copyright (c) 2019 - Telenor Norge AS


`)

}

func thingMan(thing config.Help) {
	fmt.Printf("%s\n", thing.Name)
	for l := len(thing.Name); l > 0; l-- {
		fmt.Print("-")
	}
	fmt.Printf("\n\n")
	if thing.Aliases != "" {
		fmt.Printf("Aliases: %s\n\n", thing.Aliases)
	}
	fmt.Printf("%s\n\nSettings:\n\n", thing.Doc)
	fields := []string{}
	for n := range thing.Fields {
		fields = append(fields, n)
	}
	sort.Strings(fields)
	for _, n := range fields {
		f := thing.Fields[n]
		fmt.Printf("``%s [%s]``\n\t", strings.ToLower(n), f.Type)
		fmt.Printf("%s %s\n\n", f.Doc, f.Example)
	}
}

func help() {
	flag.Usage()
	fmt.Println("\nSenders:")
	for idx, sen := range sender.Auto {
		config.PrettyPrint(idx, sen.Help)
	}
	fmt.Println("\nReceivers:")
	for idx, rcv := range receiver.Auto {
		config.PrettyPrint(idx, rcv.Help)
	}
	fmt.Println("\nYou can also see the skogul manual page. It can be generated with `./skogul -make-man > foo; rst2man < foo > skogul.1; man ./skogul.1'.")
}

func main() {
	flag.Parse()
	if *fhelp {
		help()
		os.Exit(0)
	}
	if *fman {
		man()
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
