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
	"log"

	"github.com/KristianLyng/skogul"
)

// Auto maps names to Receivers to allow auto configuration
var Auto map[string]*Receiver

// Receiver is the generic data for all receivers, used for
// auto-configuration and more.
type Receiver struct {
	Name    string
	Aliases []string
	Alloc   func() skogul.Receiver
	Help    string
}

// Add is used to announce a receiver-implementation to the world, so to
// speak. It is exported to allow out-of-package senders to exist.
func Add(r Receiver) error {
	if Auto == nil {
		Auto = make(map[string]*Receiver)
	}
	if Auto[r.Name] != nil {
		log.Panicf("BUG: Attempting to overwrite existing auto-add receiver %v", r.Name)
	}
	if r.Alloc == nil {
		log.Panicf("No alloc function for %s", r.Name)
	}
	for _, alias := range r.Aliases {
		if Auto[alias] != nil {
			log.Panicf("BUG: An alias(%s) for receiver %s overlaps an existing receiver %s", alias, r.Name, Auto[alias].Name)
		}
		Auto[alias] = &r
	}
	Auto[r.Name] = &r
	return nil
}

func init() {
	Add(Receiver{
		Name:    "http",
		Aliases: []string{"https"},
		Alloc:   func() skogul.Receiver { return &HTTP{} },
		Help:    "Listen for metrics on HTTP or HTTPS. Optionally requiring authentication. Each request received is passed to the handler.",
	})
	Add(Receiver{
		Name:  "file",
		Alloc: func() skogul.Receiver { return &File{} },
		Help:  "Reads from a file, then stops. Assumes one collection per line.",
	})
	Add(Receiver{
		Name:  "fifo",
		Alloc: func() skogul.Receiver { return &LineFile{} },
		Help:  "Reads continuously from a file. Can technically read from any file, but since it will re-open and re-read the file upon EOF, it is best suited for reading a fifo. Assumes one collection per line.",
	})
	Add(Receiver{
		Name:  "mqtt",
		Alloc: func() skogul.Receiver { return &MQTT{} },
		Help:  "Listen for Skogul-formatted JSON on a MQTT endpoint",
	})
	Add(Receiver{
		Name:  "stdin",
		Alloc: func() skogul.Receiver { return &Stdin{} },
		Help:  "Reads from standard input, one collection per line, allowing you to pipe collections to Skogul on a command line or similar.",
	})
	Add(Receiver{
		Name:  "test",
		Alloc: func() skogul.Receiver { return &Tester{} },
		Help:  "Generate dummy-data. Useful for testing, including in combination with the http sender to send dummy-data to an other skogul instance.",
	})
	Add(Receiver{
		Name:  "tcp",
		Alloc: func() skogul.Receiver { return &TCPLine{} },
		Help:  "Listen for Skogul-formatted JSON on a tcp socket, reading one collection per line.",
	})
}
