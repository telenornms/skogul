/*
 * skogul, Encoder automation
 *
 * Copyright (c) 2022 Telenor Norge AS
 * Author(s):
 *  - Roshini NarasimhaRaghavan(roshiragavi@gmail.com)
 *  - Kristian Lyngstøl <kly@kly.no>
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
Package encoder provides a generic method of encoding Skogul containers
into a byte stream. It is meant to reduce the number of semi-identical
sender modules that mainly differ in encoding.

At present time, it is not extensively used, but efforts to unify
particularly the HTTP senders are expected.
*/
package encoder

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
		Help:     "Encodes the standard Skogul JSON format.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "prettyskogul",
		Aliases:  []string{"prettyjson"},
		Alloc:    func() interface{} { x := JSON{}; x.Pretty = true; return &x },
		Help:     "Encodes the standard Skogul JSON format with indentation.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "gob",
		Alloc:    func() interface{} { return &GOB{} },
		Help:     "Encodes the GOB format.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:     "blob",
		Alloc:    func() interface{} { return &Blob{} },
		Help:     "Use data[\"data\"] as the raw message, unaltered. Optionally with a delimiter between metrics. Useful for transparently moving data in conjunction with the blob parser.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:  "avro",
		Alloc: func() interface{} { return &AVRO{} },
		Help:  "Encodes the avro format.",
	})
}
