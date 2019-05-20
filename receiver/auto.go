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
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/KristianLyng/skogul"
)

// New creates a new Receiver based on the url provided. Only receivers that
// participate in the Auto-scheme are applicable, though that SHOULD be
// most of them. However, some senders offer more functionality than what
// can be expressed in a single URL, so review specific senders if
// something is missing.
//
// See cmd/skogul-x2y for usage and a list.
func New(in string, h skogul.Handler) (skogul.Receiver, error) {
	u, err := url.Parse(in)
	if err != nil {
		return nil, skogul.Error{Source: "auto receiver", Reason: "unable to parse URL", Next: err}
	}
	if Auto[u.Scheme] == nil {
		return nil, skogul.Error{Source: "auto receiver", Reason: fmt.Sprintf("no applicable receiver for scheme %s", u.Scheme)}
	}
	x := Auto[u.Scheme].Init(*u, h)
	if x == nil {
		return nil, skogul.Error{Source: "auto receiver", Reason: fmt.Sprintf("failed to initialize receiver for %s", u.Scheme)}
	}
	return x, nil
}

// Auto maps schemas to AutoReceivers to allow skogul-x2y (and others?) to
// automatically support receivers.
var Auto map[string]*AutoReceiver

// AutoReceiver is used to initialize and document a receiver based on URL
type AutoReceiver struct {
	Init  func(url url.URL, h skogul.Handler) skogul.Receiver
	Help  string
	Flags func() *flag.FlagSet
}

func newAutoReceiver(scheme string, r *AutoReceiver) error {
	if Auto == nil {
		Auto = make(map[string]*AutoReceiver)
	}
	if Auto[scheme] != nil {
		log.Panicf("BUG: Attempting to overwrite existing auto-add receiver %v", scheme)
	}
	Auto[scheme] = r
	return nil
}

// addAutoReceiver is used by receiver-implementations to "participate" in
// the Auto-scheme described here.
func addAutoReceiver(scheme string, init func(url url.URL, h skogul.Handler) skogul.Receiver, help string) {
	if Auto == nil {
		Auto = make(map[string]*AutoReceiver)
	}
	if Auto[scheme] != nil {
		log.Panicf("BUG: Attempting to overwrite existing auto-add receiver %v", scheme)
	}
	Auto[scheme] = &AutoReceiver{Init: init, Help: help, Flags: nil}
}
