/*
 * skogul, receiver examples
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

package receiver_test

import (
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"github.com/KristianLyng/skogul/transformer"
)

/*
HTTP can have different skogul.Handler's for different paths, with potentially different behaviors.
*/
func ExampleHTTP() {
	h := receiver.HTTP{Address: "localhost:8080"}
	template := skogul.Handler{Parser: parser.JSON{}, Transformers: []skogul.Transformer{transformer.Templater{}}, Sender: sender.Debug{}}
	noTemplate := skogul.Handler{Parser: parser.JSON{}, Sender: sender.Debug{}}
	h.Handle("/template", &template)
	h.Handle("/notemplate", &noTemplate)
	h.Start()
}

/*
Using New() sets up a single handler on the specified path. This is the same as
*/
func ExampleHTTP_new() {
	handler := skogul.Handler{Parser: parser.JSON{}, Transformers: []skogul.Transformer{transformer.Templater{}}, Sender: sender.Debug{}}
	h, err := receiver.New("http://localhost:8080/foobar", handler)
	if err != nil {
		panic(err)
	}
	h.Start()
}
