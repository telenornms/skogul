/*
 * skogul, receiver examples
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

package receiver_test

import (
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/parser"
	"github.com/telenornms/skogul/receiver"
	"github.com/telenornms/skogul/sender"
	"github.com/telenornms/skogul/transformer"
)

/*
HTTP can have different skogul.Handler's for different paths, with potentially different behaviors.
*/
func ExampleHTTP() {
	h := receiver.HTTP{Address: "localhost:8080"}
	template := skogul.Handler{Transformers: []skogul.Transformer{transformer.Templater{}}, Sender: &sender.Debug{}}
	template.SetParser(parser.SkogulJSON{})
	noTemplate := skogul.Handler{Sender: &sender.Debug{}}
	noTemplate.SetParser(parser.SkogulJSON{})
	h.Handlers = map[string]*skogul.HandlerRef{
		"/template":   {H: &template},
		"/notemplate": {H: &noTemplate},
	}
	h.Start()
}
