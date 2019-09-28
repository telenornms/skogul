/*
 * skogul, sender automation
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

package sender

import (
	"log"

	"github.com/KristianLyng/skogul"
)

// Sender provides a framework that all sender-implementations should
// follow, and allows auto-initialization.
type Sender struct {
	Name    string
	Aliases []string
	Alloc   func() skogul.Sender
	Help    string
}

// Auto maps sender-names to sender implementation, used for auto
// configuration.
var Auto map[string]*Sender

// Add announces the existence of a sender to the world at large.
func Add(s Sender) error {
	if Auto == nil {
		Auto = make(map[string]*Sender)
	}
	if Auto[s.Name] != nil {
		log.Panicf("BUG: Attempting to overwrite existing auto-add sender %v", s.Name)
	}
	if s.Alloc == nil {
		log.Printf("No alloc function for %s", s.Name)
	}
	Auto[s.Name] = &s
	for _, alias := range s.Aliases {
		if Auto[alias] != nil {
			log.Panicf("BUG: An alias(%s) for sender %s overlaps an existing sender %s", alias, s.Name, Auto[alias].Name)
		}
		Auto[alias] = &s
	}
	return nil
}

func init() {
	Add(Sender{
		Name:    "backoff",
		Aliases: []string{"retry"},
		Alloc:   func() skogul.Sender { return &Backoff{} },
		Help:    "Forwards data to the next sender, retrying after a delay upon failure. For each retry, the delay is doubled. Gives up after the set number of retries.",
	})
	Add(Sender{
		Name:    "batch",
		Aliases: []string{"batcher"},
		Alloc:   func() skogul.Sender { return &Batch{} },
		Help:    "Accepts metrics and puts them in a shared container. When the container either has a set number of metrics (Threshold), or a timeout occurs, the entire container is forwarded. This allows down-stream senders to work with larger batches of metrics at a time, which is frequently more efficient. A side effect of this is that down-stream errors are not propogated upstream. That means any errors need to be dealt with down stream, or they will be ignored.",
	})
	Add(Sender{
		Name:    "counter",
		Aliases: []string{"count"},
		Alloc:   func() skogul.Sender { return &Counter{} },
		Help:    "Accepts metrics, counts them and passes them on. Then emits statistics to the Stats-handler on an interval.",
	})
	Add(Sender{
		Name:  "debug",
		Alloc: func() skogul.Sender { return &Debug{} },
		Help:  "Prints received metrics to stdout.",
	})
	Add(Sender{
		Name:    "detacher",
		Aliases: []string{"detach"},
		Alloc:   func() skogul.Sender { return &Detacher{} },
		Help:    "Returns OK without waiting for the next sender to finish.",
	})
	Add(Sender{
		Name:    "dupe",
		Aliases: []string{"dup", "duplicate"},
		Alloc:   func() skogul.Sender { return &Dupe{} },
		Help:    "Sends the same metrics to all senders listed in Next.",
	})
	Add(Sender{
		Name:    "errdiverter",
		Aliases: []string{"errordiverter", "errdivert", "errordivert"},
		Alloc:   func() skogul.Sender { return &ErrDiverter{} },
		Help:    "Forwards data to next sender. If an error is returned, the error is converted into a Skogul container and sent to the err-handler. This provides the means of logging errors through regular skogul-chains.",
	})
	Add(Sender{
		Name:  "fanout",
		Alloc: func() skogul.Sender { return &Fanout{} },
		Help:  "Fanout to a fixed number of threads before passing data on. This is rarely needed, as receivers should do this.",
	})
	Add(Sender{
		Name:  "fallback",
		Alloc: func() skogul.Sender { return &Fallback{} },
		Help:  "Tries the senders provided in Next, in order. E.g.: if the first responds OK, the second will never get data. Useful for diverting traffic to alternate paths upon failure.",
	})
	Add(Sender{
		Name:  "forwardfail",
		Alloc: func() skogul.Sender { return &ForwardAndFail{} },
		Help:  "Forwards metrics, but always returns failure. Useful in complex failure handling involving e.g. fallback sender, where it might be used to write log or stats on failure while still propogating a failure upward.",
	})
	Add(Sender{
		Name:    "http",
		Aliases: []string{"https"},
		Alloc:   func() skogul.Sender { return &HTTP{} },
		Help:    "Sends Skogul-formatted JSON-data to a HTTP endpoint (e.g.: an other Skogul instance?). Highly useful in scenarios with multiple data collection methods spread over several servers.",
	})
	Add(Sender{
		Name:    "influx",
		Aliases: []string{"influxdb"},
		Alloc:   func() skogul.Sender { return &InfluxDB{} },
		Help:    "Send to a InfluxDB HTTP endpoint.",
	})
	Add(Sender{
		Name:  "log",
		Alloc: func() skogul.Sender { return &Log{} },
		Help:  "Logs a message, mainly useful for enriching debug information in conjunction with, for example, dupe and debug.",
	})
	Add(Sender{
		Name:    "mnr",
		Aliases: []string{"m&r"},
		Alloc:   func() skogul.Sender { return &MnR{} },
		Help:    "Sends M&R line format to a TCP endpoint.",
	})
	Add(Sender{
		Name:  "mqtt",
		Alloc: func() skogul.Sender { return &MQTT{} },
		Help:  "Publishes received metrics to an MQTT broker/topic.",
	})
	Add(Sender{
		Name:  "sql",
		Alloc: func() skogul.Sender { return &SQL{} },
		Help:  "Execute a SQL query for each received metric, using a template. Any query can be run, and if multiple metrics are present in the same container, they are all executed in a single transaction, which means the batch-sender will greatly increase performance. Supported engines are MySQL/MariaDB and Postgres.",
	})
	Add(Sender{
		Name:  "null",
		Alloc: func() skogul.Sender { return &Null{} },
		Help:  "Discards all data. Mainly useful for testing.",
	})
	Add(Sender{
		Name:  "sleep",
		Alloc: func() skogul.Sender { return &Sleeper{} },
		Help:  "Injects a random delay before passing data on. Mainly for testing.",
	})
	Add(Sender{
		Name:  "test",
		Alloc: func() skogul.Sender { return &Test{} },
		Help:  "Used for internal testing. Basically just discards data but provides an internal counter of received data",
	})

}
