/*
 * skogul, receiver boilerplate
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
Receivers accept data and execute a handler. They are the "inbound"
API of Skogul.
*/
package receivers

import (
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"net/url"
)

type AutoReceiver struct {
	Scheme string
	Init   func(url url.URL, h skogul.Handler) skogul.Receiver
	Help   string
}

var Auto map[string]*AutoReceiver

func addAutoReceiver(scheme string, init func(url url.URL, h skogul.Handler) skogul.Receiver, help string) {
	if Auto == nil {
		Auto = make(map[string]*AutoReceiver)
	}
	if Auto[scheme] != nil {
		log.Fatalf("BUG: Attempting to overwrite existing auto-add receiver %v", scheme)
	}
	Auto[scheme] = &AutoReceiver{scheme, init, help}
}
