/*
 * skogul, help/man-file stuff
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.com>
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

package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/config"
	"github.com/telenornms/skogul/encoder"
	"github.com/telenornms/skogul/parser"
	"github.com/telenornms/skogul/receiver"
	"github.com/telenornms/skogul/sender"
	"github.com/telenornms/skogul/transformer"
)

// man generates an RST document suited for converting to a manual page
// using rst2man. The RST document itself is also valid, but some short
// cuts have been made, e.g., cutting long lines is not done, so the
// raw rst document might seem a bit rough, but translated to a manual page
// it looks fine.
//
// Also includes help for all senders and receivers, and uses flag to print
// the command line flag options as well.
func man() {
	fmt.Printf(`
======
skogul
======

------
Skogul
------

:Manual section: 1
:Date: %s
:Version: skogul %s
:Authors: Kristian Lyngstøl <kly@kly.no>

SYNOPSIS
========

::

	skogul -f config-file [-show]
	
	skogul [-help | -show | -make-man]

DESCRIPTION
===========

Skogul is a generic tool for moving metric data around. It can serve as a
collector of data, but is primarily designed to be a framework for building
bridges between data collectors and storage engines.

These bridges can be simple - accept data on HTTP, write to influxdb - or
complex: Accept data on unencrypted http, batch data together, forward it
to a remote skogul-instance over a password-protected, encrypted HTTPS
channel, if that fails, write to a local queue and retry periodically.

To facilitate this, Skogul has three core components:

1. Receivers acquire raw data
2. Handlers turns raw data into meaningful content, using a parser and 0 or
   more transformers.
3. Senders determine what happens to the data

A single instance of Skogul must have at least one receiver, but can have
multiple. It also, typically, must have at least one handler and sender.

Senders come in two distinct but interchangeable variants: Storage-oriented
senders are used to send the data to some external resource, e.g., a time
series database like InfluxDB. Utility-oriented senders are used to add
logic, such as error handling or duplicating data to multiple storage
systems.

There are more examples in the the "examples/" directory.

OPTIONS
=======

`, time.Now().Format("2006-01-02"), versionNo)

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
the -f option. You need to specify at least one receiver and handler to
make something sensible, you probably also want a sender.

The base configuration set is::

  {
    "parsers": {
      "xxx": {
	"type": "type-of-parser",
	type-specific-options
      },
      "other-parser...": ...
    },
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
        "transformers": [...],
        "sender": "reference-to-sender"
      }
    },
    "transformers": {
      "rrr": {
        "type": "type-of-transformer",
        type-specific-options
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

In the above pseudo-config, "xxx", "yyy", "zzz", "rrr", and "qqq" are
arbitrary names you chose that are how to reference that specific item
within the same configuration. The "type" field references what
implementation to use - each implementation have different configuration
options. You can specify as many senders, receivers and handlers as you
want, and they can cross-reference each other.

Upon start-up, all receivers are started.

It is valid to have multiple receivers use the same handler. It is also
valid for multiple senders to reference the same sender. It is up to the
operator to avoid setting up feedback loops. In other words: If you tell
a dupe-sender to send data back to itself, it will happily do so until
it runs out of memory and dies.

Numerous parsers, transformers, senders and receivers exist and they are
each documented below. Some of these can be referenced by implementation
name without defining them in the configuration - this will create a
"blank" variant of the module with default options. This is noted for each
module, and is provided to avoid bloating your configuration with modules
that don't have any required options (e.g.: the debug sender or the
templater transformer, or most parsers). Any explicitly defined module
will always take precedence over the implicitly created ones.

The documentation for each module also lists all options. In general, you
do not need to specify all options.

CONFIGURATION DATA TYPES
========================

Data types are parsed into Go types. In most cases, the the type is self
explanatory (e.g: a "string" is a typical text string, "int" is an integer,
and so on).

However, here are some examples that might not be obvious.

HandlerRef
	This is a text string referencing a named handler, specified in
	"handlers".

SenderRef
	A text string referencing a named sender, specified in "senders".

ParserRef
	A text string referencing a named parser, specified in "parsers".

[]string
	An array of text strings. E.g. ["foo","bar"].

[]*skogul.SenderRef
	An array of SenderRef, so similar to the above ["foo", "bar"], where "foo"
	and "bar" are senders named in the "senders" section of the configuration.

map[string]*skogul.HandlerRef
	This is a map of strings to handler references. For example, { "/some/path": "aHandler",
	"/other/path": "bHandler"}.

interface{}
	This is a generic "anything"-structure that can hold any arbitrary
	value. Can be any structure or variable, including nested
	variables. Used in the data/metadata transformers, among others.

SENDERS
=======

The following senders exist.

`)
	helpModules(sender.Auto)
	fmt.Print(`
RECEIVERS
=========

The following receivers exist.

`)
	helpModules(receiver.Auto)
	fmt.Print(`
ENCODERS
========

Encoders are used by a few senders to standardize how data is encoded, they
are the logical counter-point to parsers. They are mostly optional and most senders
don't use them, and those who do generally have a sensible default.

`)
	helpModules(encoder.Auto)
	fmt.Print(`
TRANSFORMERS
============

Transformers are the only tools that can actively modify a metric. See the
"HANDLERS" section for more discussion. The available transformers are:

`)
	helpModules(transformer.Auto)
	fmt.Print(`

PARSERS
=======

Parsers are responsible for parsing the data after it has been received by
a receiver. The reference-parser for Skogul is the JSON parser, but skogul
supports a number of other parsers as well.

`)
	helpModules(parser.Auto)
	fmt.Print(`
HANDLERS
========

There is only one type of handler. It accepts four arguments: A parser to
parse data, a list of optional transformers, an option to ignore partial
metric validation failures, and the first sender that will receive the
parsed container(s).

E.g.:::

  "handlers": {
	  "foo": {
		  "parser": "json",
		  "transformers": [],
		  "IgnorePartialFailures": true,
		  "sender": "blatt"
	  }
  }

IgnorePartialFailures
---------------------

If IgnorePartialFailures is set to true, invalid metrics will be removed
from a container before a sender is called, and as long as at least 1 valid
metric remains, the sender will be called. Under normal log levels, this
filtering is suppressed as long as at least 1 metric is valid, but debug
logging will reveal it.

JSON parsing
------------

If the "json" parser is used , data sent to Skogul will be parsed to fit
the internal data model of Skogul. The JSON representation is roughly
thus::

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

The "template" is optional, see the "Templater"-documentation above for an
in-depth description.

If the format of the incoming data does not conform to the Skogul JSON
structure it is possible to use a custom JSON parser ("custom-json")
which puts all the incoming data into the data-field in the Container.
After this it is possible to apply transformers to process the data further.

The primary difference between metadata and data is searchability,
and it will depend on storage engines. Typically this means the name
of a server is metadata, but the load average is data. Skogul itself
does not much care.

Juniper Telemetry Parsing
-------------------------

If the "protobuf" parser is used, the Juniper Telemetry-specific protobuf
parser is used to decode streaming telemetry from Juniper devices. Details
on how to configure your Juniper device for streaming telemetry is outside
the scope of this document.

Since streaming telemetry is sent on UDP, you need to also use the UDP
receiver. An example configuration::

	{
	  "receivers": {
	      "udp": {
		"type": "udp",
		"address": ":5015",
		"handler": "protobuf",
		"packetsize": 9000
	      }
	  },
	  "handlers": {
	    "protobuf": {
	      "parser": "protobuf",
	      "transformers": [],
	      "sender": "print"
	    }
	  },
	  "senders": {
	    "print": {
	      "type": "debug"
	    }
	  }
	}

Since the protobuf data is typically nested, you may need to use one or
more transformer before passing it on. However, senders such as the
debug-sender, HTTP-sender and SQL-sender can be used.

An example that writes to postgres::

	{
	  "receivers": {
	      "udp": {
		"type": "udp",
		"address": ":5015",
		"handler": "protobuf",
		"packetsize": 9000
	      }
	  },
	  "handlers": {
	    "protobuf": {
	      "parser": "protobuf",
	      "transformers": [],
	      "sender": "batch"
	    }
	  },
	  "senders": {
	    "batch": {
	      "type": "batch",
	      "interval": "2s",
	      "threshold": 1000,
	      "next": "psql"
	    },
	    "psql": {
	      "type": "sql",
	      "driver": "postgres",
	      "connstr": "user=skogul password=hunter2 database=telemetry sslmode=disable",
	      "query": "INSERT INTO telemetry VALUES(${timestamp}, ${json.metadata}, ${json.data})"
	    }
	  }
	}

Minimalistic schema::

			       Table "public.telemetry"
	  Column  |           Type           | Collation | Nullable | Default
	----------+--------------------------+-----------+----------+---------
	 ts       | timestamp with time zone |           |          |
	 metadata | jsonb                    |           |          |
	 data     | jsonb                    |           |          |

Custom JSON
-----------

If the "custom-json" parser is used, all data is parsed as if it was JSON,
but put into a single metric with the current time as timestamp. This
allows sending generic JSON to Skogul and using transformers to extract
relevant data (or just store everything as data).

InfluxDB parser
---------------

A simple influx line protocol parser is provided and can be used by setting
the parser to "influx".

Templating
----------

The templating-transformer is useful for adding identical fields to all
metrics in a collection. If a template is provided, and the
templater-transformer is applied, all metrics are initialized with whatever
value the template came with.

This is inteded for when you are sending multiple metrics that share
certain attributes, e.g, they are all from the same machine and all
collected at the same time. Or they are all from the same data center
or region.

Templates are shallow. If your metric has nested fields, they will not
be merged with what the template provides. For example::

   {
     "template": {
       "timestamp": "2019-09-27T15:42:00Z",
       "metadata": {
         "site": "naboo",
         "machine": {
           "os": "Debian"
         }
       }
     },
     "metrics": [
       {
         "metadata": {
           "machine": {
             "hostname": "r2d2"
           }
         },
         "data": {
           "something": "blah"
         }
       },
       {
         "metadata": {
           "machine": {
             "hostname": "c3po"
           }
         },
         "data": {
           "something": "duck"
         }
       }
     ]
   }

Here, the template provides three items: a timestamp, the "site" field and
the "machine" field of metadata. Once transformed, the result will be::

   {
     "metrics": [
       {
         "timestamp": "2019-09-27T15:42:00Z",
         "metadata": {
           "site": "naboo",
           "machine": {
             "hostname": "r2d2"
           }
         },
         "data": {
           "something": "blah"
         }
       },
       {
         "timestamp": "2019-09-27T15:42:00Z",
         "metadata": {
           "site": "naboo",
           "machine": {
             "hostname": "c3po"
           }
         },
         "data": {
           "something": "duck"
         }
       }
     ]
   }

Since each metric also provided a "machine"-field, it overwrote the value
from the template, even if there were no overlapping fields.


EXAMPLES
========

A minimalistic example that accepts data on HTTP and prints it to standard
output::

  { 
    "receivers": { 
      "api": { 
        "type": "http", 
        "address": ":8080", 
        "handlers": { "/": "myhandler" }
      }
    },
    "handlers": {
      "myhandler": {
        "parser": "json", 
        "transformers": ["templater"], 
        "sender": "debug"
      }
    }
  }

Note that parser, transformers and senders used here can be used without
being explicitly configured. This is only true for very simple modules and is
noted in each module's documentation section.

The same example can be written more verbosly::

  { 
    "receivers": { 
      "api": { 
        "type": "http", 
        "address": ":8080", 
        "handlers": { "/": "myhandler" }
      }
    },
    "handlers": {
      "myhandler": {
        "parser": "json", 
        "transformers": ["templater"], 
        "sender": "debug"
      }
    },
    "parsers": {
      "json": {
        "type": "json"
      }
    },
    "transformers": {
      "templater": {
	"type": "templater"
      }
    },
    "senders": {
      "debug": {
	"type": "debug",
	"prefix": "foobar"
      }
    }
  }

This is the exact same as the first example, with one minor exception: We
give the "debug" sender an optional prefix.

You can have multiple variants of the same modules with different options,
which obviously requires you to configure each with unique names::

  { 
    "receivers": { 
      "api": { 
        "type": "http", 
        "address": ":8080", 
        "handlers": { 
		"/": "silent",
		"/noise": "noisy" 
	}
      }
    },
    "handlers": {
      "silent": {
        "parser": "json", 
        "transformers": ["templater"], 
        "sender": "debug-clean"
      },
      "noisy": {
        "parser": "json", 
        "transformers": ["templater"], 
        "sender": "debug-noisy"
      }
    },
    "senders": {
      "debug-clean": {
	"type": "debug",
      },
      "debug-noisy": {
	"type": "debug",
	"prefix": "noise"
      }
    }
  }

The above sets up a HTTP receiver on port 8080 that will use the "silent"
handler for data sent to  /, which in turn will use the "debug-clean"
sender without a prefix, while data sent to /noise will use the "noisy"
handler and "debug-noisy" sender with a prefixed message of "noise".

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

To add a metadata field to signal where data came from before passing it on
to a central instance::

  {
    "receivers": {
      "local": {
        "type": "http",
        "address": "[::1]:8080",
        "handlers": {
          "/": "jsontemplating"
        }
      }
    },
    "transformers": {
      "origin": {
        "type": "metadata",
        "set": {
          "dc": "bergen1",
          "collector": "serverX"
        }
      }
    },
    "handlers": {
      "jsontemplating": {
        "parser": "json",
        "transformers": [ "templater","metadata" ],
        "sender": "batch"
      }
    },
    "senders": {
      "batch": {
        "type": "batch",
        "interval": "5s",
        "threshold": 1000,
        "next": "central"
      },
      "central": {
        "type": "http",
        "url": "https://bergen1X:hunter2@central-skogul.example.com/",
        "Timeout": "10s"
      }
    }
  }

More examples are provided in the examples/ directory of the Skogul source
package.

SEE ALSO
========

https://github.com/telenornms/skogul

BUGS
====

While Skogul tries to provide sensible feedback on cofiguration errors, it
CAN be complicated.

Workaround: Use the "-show" option to display the parsed configuration.

COPYRIGHT
=========

This document is licensed under the same license as Skogul itself, which
happens to be LGPLv2 (or later). See LICENSE for details.

* Copyright (c) 2019-2020 - Telenor Norge AS

`)

}

// helpModules iterates over a ModuleMap, printing rst-formatted help for
// each module.
func helpModules(mmap skogul.ModuleMap) {
	mods := []string{}
	for idx := range mmap {
		if mmap[idx].Name != idx {
			continue // alias
		}
		mods = append(mods, idx)
	}
	sort.Strings(mods)
	for _, mod := range mods {
		mh, _ := config.HelpModule(mmap, mod)
		thingMan(mh)
	}
}

// fieldDoc iterates of FieldDoc to pretty-print it for rst.
func fieldDoc(inFields map[string]config.FieldDoc) {
	fields := []string{}
	doit := false
	for n := range inFields {
		fields = append(fields, n)
		doit = true

	}
	if doit {
		fmt.Printf("Settings:\n\n")
	}
	sort.Strings(fields)
	for _, n := range fields {
		f := inFields[n]
		fmt.Printf("``%s - %s``\n\t", strings.ToLower(n), f.Type)
		fmt.Printf("%s\n\n", strings.Replace(f.Doc, "\n", "\n\t", -1))
		if f.Example != "" {
			fmt.Printf("\tExample(s): %s\n\n", f.Example)
		}
	}
}

// thingMan is thus named because of reasons. It prints RST-formatted
// documentation for a sender or receiver, whatever config.Help has.
func thingMan(thing config.Help) {
	fmt.Printf("%s\n", thing.Name)
	for l := len(thing.Name); l > 0; l-- {
		fmt.Print("-")
	}
	fmt.Printf("\n\n")
	fmt.Printf("%s\n\n", thing.Doc)
	if thing.Aliases != "" {
		fmt.Printf("Aliases: %s\n\n", thing.Aliases)
	}
	if thing.AutoMake {
		fmt.Printf("This module does not have to be explicitly defined in your configuration file. If you refer to it by one of it's implementation names, a default-variant of it will be used.\n\n")
	}
	fieldDoc(thing.Fields)
	if len(thing.CustomTypes) > 0 {
		for ctype, n := range thing.CustomTypes {
			fmt.Printf("Custom type ``%s``\n\n", ctype)
			fieldDoc(n)
		}
	}
}
