/*
 * skogul, parser automation
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
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

package parser

import (
	"github.com/telenornms/skogul"
)

// Auto maps parser-names to parser implementation, used for auto
// configuration.
var Auto skogul.ModuleMap

func init() {
	Auto.Add(skogul.Module{
		Name:     "skogul",
		Aliases:  []string{"json"},
		Alloc:    func() interface{} { return &JSON{} },
		Help:     "Parses the standard Skogul JSON format.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "skogulmetric",
		Aliases:  []string{"jsonmetric", "json1"},
		Alloc:    func() interface{} { return &JSONMetric{} },
		Help:     "Parses the byte stream as a single json-encoded metric.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "rawjson",
		Aliases:  []string{"jsonraw", "custom-json"},
		Alloc:    func() interface{} { return &RawJSON{} },
		Help:     "Parses any generic JSON data into a single metric which can then be potentially transformed into multiple metrics if need be.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "influxdb",
		Aliases:  []string{"influx"},
		Alloc:    func() interface{} { return &InfluxDB{} },
		Help:     "Parse InfluxDB line-protocol data",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "protobuf",
		Aliases:  []string{"telemetry", "juniper"},
		Alloc:    func() interface{} { return &ProtoBuf{} },
		Help:     "Parse Juniper telemetry in the form of protocol buffers. Typicially combined with the UDP receiver.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "protobuf_usp",
		Aliases:  []string{"usp"},
		Alloc:    func() interface{} { return &USP_Parser{} },
		Help:     "Parse Usp message contained within Usp record in the form of protocol buffers.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "mnr",
		Aliases:  []string{"m&r"},
		Alloc:    func() interface{} { return &MNR{} },
		Help:     "Parse M&R internal data",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "structured_data",
		Aliases:  []string{},
		Alloc:    func() interface{} { return &StructuredData{} },
		Help:     "Parse structured data as specified in RFC5424",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "gob",
		Aliases:  []string{},
		Alloc:    func() interface{} { return &GOB{} },
		Help:     "Parse a GOB-encoded Skogul Container. GOB is a go-specific encoding, useful for inter-Skogul communication.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "gobmetric",
		Aliases:  []string{},
		Alloc:    func() interface{} { return &GOBMetric{} },
		Help:     "Parse a single GOB-encoded Skogul Metric",
		AutoMake: true,
	})
}
