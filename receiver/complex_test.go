/*
 * skogul, complex receiver tests
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
	"testing"
	"time"
)

var validContainer = skogul.Container{}

func init() {

	now := time.Now()

	m := skogul.Metric{}
	m.Time = &now
	m.Metadata = make(map[string]interface{})
	m.Data = make(map[string]interface{})
	m.Metadata["foo"] = "bar"
	m.Data["tall"] = 5
	validContainer.Metrics = []*skogul.Metric{&m}
}

// Tests http receiver, sender and JSON parser implicitly
func TestHttp_stack(t *testing.T) {
	one := &(sender.Test{})

	h := skogul.Handler{Sender: one, Parser: parser.JSON{}}
	rcv := receiver.HTTP{Address: "[::1]:1339", Handlers: map[string]*skogul.HandlerRef{"/foo": &skogul.HandlerRef{H: &h}}}
	go rcv.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
	hs := sender.HTTP{URL: "http://[::1]:1339/foo"}
	one.TestQuick(t, &hs, &validContainer, 1)
	blank := skogul.Container{}
	one.TestNegative(t, &hs, &blank)

	hs2 := sender.HTTP{URL: "http://[::1]:1339/bar"}
	one.TestNegative(t, &hs2, &validContainer)
}
