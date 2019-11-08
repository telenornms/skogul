
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
2. Handlers turns raw data into meaningful content
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

``-f`` string
	Path to skogul config to read. (default ~/.config/skogul.json)

``-help``
	Print more help (default false)

``-loglevel`` string
	Minimum loglevel to display ([e]rror, [w]arn, [i]nfo, [d]ebug, [t]race/[v]erbose) (default warn)

``-make-man``
	Output RST documentation suited for rst2man (default false)

``-show``
	Print the parsed JSON config instead of starting (default false)

``-timestamp``
	Include timestamp in log entries (default true)


CONFIGURATION
=============

Configuration of skogul is done with a json config file, referenced with
the -f option. You need to specify at least one receiver, handler and
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
operator to avoid setting up feedback loops.

Two parsers exist: the JSON parser and a Juniper Telemetry protobuf parser.
Only three transformers exists, and to simplify configuration, the
"templater" transformer does not have to be explicitly defined to be
referenced.

The documentation for each sender and receiver also lists all options. In
general, you do not need to specify all options.

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

[]string
	An array of text strings. E.g. ["foo","bar"].

[]*skogul.HandlerRef
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

backoff
-------

Forwards data to the next sender, retrying after a delay upon failure. For each retry, the delay is doubled. Gives up after the set number of retries.

Aliases: retry 

Settings:

``base - Duration``
	Initial delay after a failure. Will double for each retry

``next - SenderRef``
	The sender to try

``retries - uint64``
	Number of retries before giving up

batch
-----

Accepts metrics and puts them in a shared container. When the container either has a set number of metrics (Threshold), or a timeout occurs, the entire container is forwarded. This allows down-stream senders to work with larger batches of metrics at a time, which is frequently more efficient. A side effect of this is that down-stream errors are not propogated upstream. That means any errors need to be dealt with down stream, or they will be ignored.

Aliases: batcher 

Settings:

``interval - Duration``
	Flush the bucket after this duration regardless of how full it is

``next - SenderRef``
	Sender that will receive batched metrics

``threshold - int``
	Flush the bucket after reaching this amount of metrics

counter
-------

Accepts metrics, counts them and passes them on. Then emits statistics to the Stats-handler on an interval.

Aliases: count 

Settings:

``next - SenderRef``
	Reference to the next sender in the chain

``period - Duration``
	How often to emit stats

	Example(s): 5s

``stats - HandlerRef``
	Handler that will receive the stats periodically

debug
-----

Prints received metrics to stdout.

Settings:

``prefix - string``
	Prefix to print before any metric

detacher
--------

Returns OK without waiting for the next sender to finish.

Aliases: detach 

Settings:

``depth - int``
	How many containers can be pending delivery before we start blocking. Defaults to 1000.

``next - SenderRef``
	Sender that receives the metrics.

dupe
----

Sends the same metrics to all senders listed in Next.

Aliases: duplicate dup 

Settings:

``next - []skogul.SenderRef``
	List of senders that will receive metrics, in order.

errdiverter
-----------

Forwards data to next sender. If an error is returned, the error is converted into a Skogul container and sent to the err-handler. This provides the means of logging errors through regular skogul-chains.

Aliases: errordivert errdivert errordiverter 

Settings:

``err - HandlerRef``
	If the sender under Next fails, convert the error to a metric and send it here.

``next - SenderRef``
	Send normal metrics here.

``reterr - bool``
	If true, the original error from Next will be returned, if false, both Next AND Err has to fail for Send to return an error.

fallback
--------

Tries the senders provided in Next, in order. E.g.: if the first responds OK, the second will never get data. Useful for diverting traffic to alternate paths upon failure.

Settings:

``next - []skogul.SenderRef``
	Ordered list of senders that will potentially receive metrics.

fanout
------

Fanout to a fixed number of threads before passing data on. This is rarely needed, as receivers should do this.

Settings:

``next - SenderRef``
	Sender receiving the metrics

``workers - int``
	Number of worker threads in use. To _fan_in_ you can set this to 1.

forwardfail
-----------

Forwards metrics, but always returns failure. Useful in complex failure handling involving e.g. fallback sender, where it might be used to write log or stats on failure while still propogating a failure upward.

Settings:

``next - SenderRef``
	Sender receiving the metrics

http
----

Sends Skogul-formatted JSON-data to a HTTP endpoint (e.g.: an other Skogul instance?). Highly useful in scenarios with multiple data collection methods spread over several servers.

Aliases: https 

Settings:

``connsperhost - int``
	Max concurrent connections per host. Should reflect ulimit -n. Defaults to unlimited.

``idleconnsperhost - int``
	Max idle connections retained per host. Should reflect expected concurrency. Defaults to 2 + runtime.NumCPU.

``insecure - bool``
	Disable TLS certificate validation.

``rootca - string``
	Path to an alternate root CA used to verify server certificates. Leave blank to use system defaults.

``timeout - Duration``
	HTTP timeout.

``url - string``
	Fully qualified URL to send data to.

	Example(s): http://localhost:6081/ https://user:password@[::1]:6082/

influx
------

Send to a InfluxDB HTTP endpoint. The sender can either always send the data to a single measurement, send it to a measurement extracted from the metadata of a metric, or a combination where the "measurement" serves as a default measurement to use if the metric doesn't have the key presented in "measurementfrommetadata".

Aliases: influxdb 

Settings:

``measurement - string``
	Measurement name to write to.

``measurementfrommetadata - string``
	Metadata key to read the measurement from. Either this or 'measurement' must be set. If both are present, 'measurement' will be used if the named metadatakey is not found.

``timeout - Duration``
	HTTP timeout

``url - string``
	URL to InfluxDB API. Must include write end-point and database to write to.

	Example(s): http://[::1]:8086/write?db=foo

log
---

Logs a message, mainly useful for enriching debug information in conjunction with, for example, dupe and debug.

Settings:

``message - string``
	Message to print.

mnr
---

Sends M&R line format to a TCP endpoint.

Aliases: m&r 

Settings:

``address - string``
	Address to send data to

	Example(s): 192.168.1.99:1234

``defaultgroup - string``
	Default group to use if the metadatafield group is missing.

mqtt
----

Publishes received metrics to an MQTT broker/topic.

Settings:

``address - string``
	URL-encoded address.

	Example(s): mqtt://user:password@server/topic

net
---

Sends json data to a network endpoint.

Settings:

``address - string``
	Address to send data to

	Example(s): 192.168.1.99:1234

``network - string``
	Network, according to net.Dial. Typically udp or tcp.

null
----

Discards all data. Mainly useful for testing.

sleep
-----

Injects a random delay before passing data on. Mainly for testing.

Settings:

``base - Duration``
	The baseline - or minimum - delay

``maxdelay - Duration``
	The maximum delay we will suffer

``next - SenderRef``
	Sender that will receive delayed metrics

``verbose - bool``
	If set to true, will log delay durations

sql
---

Execute a SQL query for each received metric, using a template. Any query can be run, and if multiple metrics are present in the same container, they are all executed in a single transaction, which means the batch-sender will greatly increase performance. Supported engines are MySQL/MariaDB and Postgres.

Settings:

``connstr - string``
	Connection string to use for database. Slight variations between database engines. For MySQL typically user:password@host/database.

	Example(s): mysql: 'root:lol@/mydb' postgres: 'user=pqgotest dbname=pqgotest sslmode=verify-full'

``driver - string``
	Database driver/system. Currently suported: mysql and postgres.

``query - string``
	Query run for each metric. The following expansions are made:
	
	${timestamp} is expanded to the actual metric timestamp.
	
	${metadata.KEY} will be expanded to the metadata with key name "KEY".
	
	${data.KEY} will be expanded to data[foo].
	
	${json.metadata} will be expanded to a json representation of all metadata.
	
	${json.data} will be expanded to a json representation of all data.
	
	Finally, ${KEY} is a shorthand for ${data.KEY}. Both methods are provided, to allow referencing data fields named "metadata.". E.g.: ${data.metadata.x} will match data["metadata.x"], while ${metadata.x} will match metadata["x"].

	Example(s): INSERT INTO test VALUES(${timestamp},${hei},${metadata.key1})

test
----

Used for internal testing. Basically just discards data but provides an internal counter of received data


RECEIVERS
=========

The following receivers exist.

fifo
----

Reads continuously from a file. Can technically read from any file, but since it will re-open and re-read the file upon EOF, it is best suited for reading a fifo. Assumes one collection per line.

Settings:

``delay - Duration``
	Delay before re-opening the file, if any.

``file - string``
	Path to the fifo or file from which to read from repeatedly.

``handler - HandlerRef``
	Handler used to parse and transform and send data.

file
----

Reads from a file, then stops. Assumes one collection per line.

Settings:

``file - string``
	Path to the file to read from once.

``handler - HandlerRef``
	Handler used to parse, transform and send data.

http
----

Listen for metrics on HTTP or HTTPS. Optionally requiring authentication. Each request received is passed to the handler.

Aliases: https 

Settings:

``address - string``
	Address to listen to.

	Example(s): [::1]:80 [2001:db8::1]:443

``certfile - string``
	Path to certificate file for TLS. If left blank, un-encrypted HTTP is used.

``handlers - map[string]*skogul.HandlerRef``
	Paths to handlers. Need at least one.

	Example(s): {"/": "someHandler" }

``keyfile - string``
	Path to key file for TLS.

``password - string``
	Password for basic authentication.

``username - string``
	Username for basic authentication. No authentication is required if left blank.

log
---

Log attaches to the internal logging of Skogul and diverts log messages.

Settings:

``echo - bool``
	Logs are also echoed to stdout.

``handler - HandlerRef``
	Reference to a handler where the data is sent. Parser will be overwritten.

mqtt
----

Listen for Skogul-formatted JSON on a MQTT endpoint

Settings:

``address - string``
	Address to connect to.

``handler - *skogul.HandlerRef``
	Handler used to parse, transform and send data.

``password - string``
	Username for authenticating to the broker.

``username - string``
	Password for authenticating.

stdin
-----

Reads from standard input, one collection per line, allowing you to pipe collections to Skogul on a command line or similar.

Settings:

``handler - HandlerRef``
	Handler used to parse, transform and send data.

tcp
---

Listen for Skogul-formatted JSON on a tcp socket, reading one collection per line.

Settings:

``address - string``
	Address and port to listen to.

	Example(s): [::1]:3306

``handler - HandlerRef``
	Handler used to parse, transform and send data.

test
----

Generate dummy-data. Useful for testing, including in combination with the http sender to send dummy-data to an other skogul instance.

Settings:

``delay - Duration``
	Sleep time between each metric is generated, if any.

``handler - HandlerRef``
	Reference to a handler where the data is sent

``metrics - int64``
	Number of metrics in each container

``threads - int``
	Threads to spawn

``values - int64``
	Number of unique values for each metric

udp
---

Accept UDP messages, parsed by specified handler. E.g.: Protobuf.

Settings:

``address - string``
	Address and port to listen to.

	Example(s): [::1]:3306

``handler - HandlerRef``
	Handler used to parse, transform and send data.



TRANSFORMERS
============

Transformers are the only tools that can actively modify a metric. See the
"HANDLERS" section for more discussion. Note that the "templater" transformer
does not need to be defined - if a handler lists "templater", one will be
created behind the scenes. The available transformers are:

data
----

Enforces custom-rules for data fields of metrics.

Settings:

``ban - []string``
	Fail if any of these data fields are present

``flatten - [][]string``
	Flatten nested structures down to the root level

``remove - []string``
	Remove these data fields.

``require - []string``
	Require the pressence of these data fields.

``set - map[string]interface {}``
	Set data fields to specific values.

metadata
--------

Enforces custom-rules on metadata of metrics.

Settings:

``ban - []string``
	Fail if any of these fields are present

``extractfromdata - []string``
	Extract a set of fields from Data and add it to Metadata.

``remove - []string``
	Remove these metadata fields.

``require - []string``
	Require the pressence of these fields.

``set - map[string]interface {}``
	Set metadata fields to specific values.

replace
-------

Uses a regular expression to replace the content of a metadata key, storing it to either a different metadata key, or overwriting the original.

Settings:

``destination - string``
	Metadata key to write to. Defaults to overwriting the source-key if left blank. Destination key will always be overwritten, e.g., even if the source key is missing, the key located at the destination will be removed.

``regex - string``
	Regular expression to match.

``replacement - string``
	Replacement text. Can also use $1, $2, etc to reference sub-matches. Defaults to empty string - remove matching items.

``source - string``
	Metadata key to read from.

split
-----

Splits a metric into multiple metrics based on a field.

Settings:

``fail - bool``
	Fail the transformer entirely if split is unsuccsessful on a metric container. This will prevent successive transformers from working.

``field - []string``
	Split into multiple metrics based on this field (each field denotes the path to a nested object element).

templater
---------

Executes metric templating. See separate documentationf or how skogul templating works.

Aliases: templating template 


HANDLERS
========

There is only one type of handler. It accepts three arguments: A parser to
parse data, a list of optional transformers, and the first sender that will
receive the parsed container(s).

The only valid parsers are "json" and "protobuf". Only three transformers
exist. The "templating" transformer does not need to be explicitly defined
to be referenced, since it has no settings.

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
		"handler": "protobuf"
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
		"handler": "protobuf"
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
        "sender": "mysender"
      }
    },
    "senders": {
      "mysender": {
        "type": "debug"
      }
    }
  }

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

Configuration parsing doesn't provide very helpful errors, and silently
ignores keys/variables that are not used in a specific context.

Workaround: Use the "-show" option to display the parsed configuration.

COPYRIGHT
=========

This document is licensed under the same license as Skogul itself, which
happens to be GPLv2 (or later). See LICENSE for details.

* Copyright (c) 2019 - Telenor Norge AS

