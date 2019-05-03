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
	"fmt"
	"log"
	"net/url"

	"github.com/KristianLyng/skogul"
)

// New creates a new Sender based on the url provided. Only senders that
// participate in the Auto-scheme are applicable, mostly senders that
// actually store data. See cmd/skogul-x2y for usage and a list.
func New(in string) (skogul.Sender, error) {
	u, err := url.Parse(in)
	if err != nil {
		return nil, skogul.Error{Source: "auto sender", Reason: "unable to parse URL", Next: err}
	}
	if Auto[u.Scheme] == nil {
		return nil, skogul.Error{Source: "auto sender", Reason: fmt.Sprintf("no applicable sender for scheme %s", u.Scheme)}
	}
	x := Auto[u.Scheme].Init(*u)
	if x == nil {
		return nil, skogul.Error{Source: "auto sender", Reason: fmt.Sprintf("failed to initialize sender for %s", u.Scheme)}
	}
	return x, nil
}

// AutoSender is used to provide generic constructors by URL/Scheme.
type AutoSender struct {
	Scheme string
	Init   func(url url.URL) skogul.Sender
	Help   string
}

// Auto maps schemas to senders and help text to make appropriate senders.
var Auto map[string]*AutoSender

func addAutoSender(scheme string, init func(url url.URL) skogul.Sender, help string) {
	if Auto == nil {
		Auto = make(map[string]*AutoSender)
	}
	if Auto[scheme] != nil {
		log.Fatalf("BUG: Attempting to overwrite existing auto-add sender %v", scheme)
	}
	Auto[scheme] = &AutoSender{scheme, init, help}
}
