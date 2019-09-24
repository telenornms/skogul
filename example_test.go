/*
 * skogul, examples
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

package skogul_test

import (
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
)

// Example of the simplest Skogul chain possible
func Example() {
	// Create a debug-sender. A Debug-sender just prints the metric to
	// stdout.
	s := &sender.Debug{}

	// A handler is used to inform a receiver how to treat incoming
	// data. This one will parse it using the JSON parser, then send it
	// on to the above sender.
	h := skogul.Handler{Parser: parser.JSON{}, Sender: s}

	// Create a receiver. The receiver.New() will parse a URL to find
	// an underlying receiver that implements the schema. In this case,
	// it will use the HTTP receiver.
	r, err := receiver.New("http://localhost:1234", h)
	if err != nil {
		panic(err)
	}

	// Finally, start the receiver.
	r.Start()
}
