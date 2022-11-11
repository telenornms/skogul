/*
 * skogul, test receiver tests (... I know)
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
	"testing"
	"time"
)

func TestTester_stack(t *testing.T) {
	one := &(sender.Test{})
	h := skogul.Handler{Sender: one}
	h.SetParser(parser.SkogulJSON{})
	rcv := receiver.Tester{Metrics: 10, Values: 5, Threads: 2, Handler: skogul.HandlerRef{H: &h}}
	go rcv.Start()

	zzz := 300 * time.Millisecond
	time.Sleep(time.Duration(zzz))

	// With atomic covermode and race condition testing, this is pretty
	// slow, hence the modest values
	got := one.Received()
	want := uint64(500)
	if got < want {
		t.Errorf("Less than %d received events after %v of tester running. Got %d", want, zzz, got)
	}
}

func TestTester_auto(t *testing.T) {
	one := &(sender.Test{})
	h := skogul.Handler{Sender: one}
	h.SetParser(parser.SkogulJSON{})
	parsedD, _ := time.ParseDuration("1s")
	x := receiver.Tester{
		Handler: skogul.HandlerRef{H: &h},
		Metrics: 1,
		Values:  1,
		Threads: 2,
		Delay:   skogul.Duration{Duration: parsedD},
	}
	go x.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
	got := one.Received()
	if got < 1 {
		t.Errorf("receiver.Tester{}, x.Start() failed to receive data. Expected some data, got 0.")
	}
	if got > 2 {
		t.Errorf("Started tester with 2 threads and 1s delay, expected 2 items sent after 100ms, got %d", got)
	}
}
