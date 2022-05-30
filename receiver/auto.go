/*
 * skogul, receiver automation
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
Package receiver provides various skogul Receivers that accept data and
execute a handler. They are the "inbound" API of Skogul.
*/
package receiver

import (
	"github.com/telenornms/skogul"
)

// Auto maps names to Receivers to allow auto configuration
var Auto skogul.ModuleMap

func init() {
	Auto.Add(skogul.Module{
		Name:    "http",
		Aliases: []string{"https"},
		Alloc:   func() interface{} { return &HTTP{} },
		Help:    "Listen for metrics on HTTP or HTTPS. Optionally requiring authentication. Each request received is passed to a handler, and a single HTTP receiver can listen for multiple formats depending on URL used.",
		Extras:  []interface{}{HTTPAuth{}},
	})
	Auto.Add(skogul.Module{
		Name:  "file",
		Alloc: func() interface{} { return &File{} },
		Help:  "Reads from a file, then stops. Assumes one collection per line. E.g.: If the file has json data, the each line has to be a self-contained document/container.",
	})
	Auto.Add(skogul.Module{
		Name:    "wholefile",
		Aliases: []string{"wfile"},
		Alloc:   func() interface{} { return &WholeFile{} },
		Help:    "Reads an entire file and parses it as a single container, optionally repeatedly.",
	})
	Auto.Add(skogul.Module{
		Name:  "fifo",
		Alloc: func() interface{} { return &LineFile{} },
		Help:  "Reads continuously from a file. Can technically read from any file, but since it will re-open and re-read the file upon EOF, it is best suited for reading a fifo. Assumes one collection per line.",
	})
	Auto.Add(skogul.Module{
		Name:    "logrus",
		Aliases: []string{"log"},
		Alloc:   func() interface{} { return &LogrusLog{} },
		Help:    "Attaches to the internal logging of Skogul and diverts log messages.",
	})
	Auto.Add(skogul.Module{
		Name:  "mqtt",
		Alloc: func() interface{} { return &MQTT{} },
		Help:  "Listen for Skogul-formatted JSON on a MQTT endpoint.",
	})
	Auto.Add(skogul.Module{
		Name:  "stats",
		Alloc: func() interface{} { return &Stats{} },
		Help:  "Gather internal Skogul metrics and send them on to the specified handler. Metrics gathered depends on modules used, and verbosity and completeness also depends on the modules. Examples of metrics gathered are: parse errors, send errors, number of received messages.",
	})
	Auto.Add(skogul.Module{
		Name:  "stdin",
		Alloc: func() interface{} { return &Stdin{} },
		Help:  "Reads from standard input, one collection per line, allowing you to pipe collections to Skogul on a command line or similar.",
	})
	Auto.Add(skogul.Module{
		Name:  "test",
		Alloc: func() interface{} { return &Tester{} },
		Help:  "Generate dummy-data. Useful for testing, including in combination with the http sender to send dummy-data to an other skogul instance.",
	})
	Auto.Add(skogul.Module{
		Name:  "tcp",
		Alloc: func() interface{} { return &TCPLine{} },
		Help:  "Listen for data on a tcp socket, reading one collection per line.",
	})
	Auto.Add(skogul.Module{
		Name:  "sql",
		Alloc: func() interface{} { return &SQL{} },
		Help:  "Periodically poll a database for information. Single threaded.",
	})
	Auto.Add(skogul.Module{
		Name:  "udp",
		Alloc: func() interface{} { return &UDP{} },
		Help:  "Accept UDP messages, one UDP message is one container. Combine with protobuf parser to receive Juniper telemetry.",
	})
	Auto.Add(skogul.Module{
		Name:  "kafka",
		Alloc: func() interface{} { return &Kafka{} },
		Help:  "Connect to a Kafka topic and consume messages.",
	})
}
