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
		Help:    "Passes data on, but will retry up to Retries times, with an exponential delay between retries.",
	})
	Add(Sender{
		Name:    "batch",
		Aliases: []string{"batcher"},
		Alloc:   func() skogul.Sender { return &Batch{} },
		Help:    "Collects multiple metrics into a single container before sending them on in a batch. Data is sent when either a treshold of metrics or a timeout is reached, whichever comes first.",
	})
	Add(Sender{
		Name:    "counter",
		Aliases: []string{"count"},
		Alloc:   func() skogul.Sender { return &Counter{} },
		Help:    "Passes the metrics on to the Next sender, but every Period it will send statistics on how much data it has seen to Stats",
	})
	Add(Sender{
		Name:  "debug",
		Alloc: func() skogul.Sender { return &Debug{} },
		Help:  "Prints received metrics to stdout",
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
		Help:    "Calles the next sender, but if that fails, the error itself will be converted to a Container and sent to Err, allowing error handling or stats of errors",
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
		Help:  "Forwards metrics, but always returns failure. Useful in complex failure handling involving e.g. fallback sender, where it might be used to write log or stats on failure while still propogating a failure upward",
	})
	Add(Sender{
		Name:    "http",
		Aliases: []string{"https"},
		Alloc:   func() skogul.Sender { return &HTTP{} },
		Help:    "Sends Skogul-formatted JSON-data to a HTTP endpoint (e.g.: an other Skogul instance?)",
	})
	Add(Sender{
		Name:    "influx",
		Aliases: []string{"influxdb"},
		Alloc:   func() skogul.Sender { return &InfluxDB{} },
		Help:    "Send to a InfluxDB HTTP endpoint",
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
		Help:    "Sends M&R line format to a TCP endpoint",
	})
	Add(Sender{
		Name:  "mqtt",
		Alloc: func() skogul.Sender { return &MQTT{} },
		Help:  "Publishes received metrics to an MQTT broker/topic",
	})
	Add(Sender{
		Name:  "mysql",
		Alloc: func() skogul.Sender { return &Mysql{} },
		Help:  "Execute a MySQL query for each received metric, using a template.",
	})
	Add(Sender{
		Name:  "null",
		Alloc: func() skogul.Sender { return &Null{} },
		Help:  "Discards all data. Mainly useful for testing.",
	})
	Add(Sender{
		Name:  "sleep",
		Alloc: func() skogul.Sender { return &Sleeper{} },
		Help:  "Injects a random delay before passing data on.",
	})
	Add(Sender{
		Name:  "test",
		Alloc: func() skogul.Sender { return &Test{} },
		Help:  "Used for internal testing. Basically just discards data but provides an internal counter of received data",
	})

}
