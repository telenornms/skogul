/*
 * skogul, detach sender - tests
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

package sender_test

import (
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/sender"
	"testing"
	"time"
)

func TestDetacher(t *testing.T) {
	c := skogul.Container{}
	m := skogul.Metric{}

	c.Metrics = []*skogul.Metric{&m}
	tst := &(testSender{})
	delay := &(sender.Sleeper{Base: time.Duration(100 * time.Millisecond), MaxDelay: time.Duration(100 * time.Millisecond), Next: tst})
	detach := &(sender.Detacher{Next: delay})

	start := time.Now()
	tst.testQuick(t, detach, &c, 0)

	diff := time.Since(start)
	if diff > (20 * time.Millisecond) {
		t.Errorf("Took too long sending to the detach-sender. Took more than 20ms (%v). Should be ~instant.", diff)
	}
	time.Sleep(200 * time.Millisecond)
	if tst.received != 1 {
		t.Errorf("Didn't get the event after time expired? Wanted %d containers, got %d", 1, tst.received)
	}
}

func TestFanout(t *testing.T) {
	c := skogul.Container{}
	m := skogul.Metric{}

	c.Metrics = []*skogul.Metric{&m}
	tst := &(testSender{})
	delay := &(sender.Sleeper{Base: time.Duration(100 * time.Millisecond), MaxDelay: time.Duration(100 * time.Millisecond), Next: tst})
	fanout := &(sender.Fanout{Next: delay, Workers: 3})

	start := time.Now()
	// With a work queue of 3, we can send 3 at the same time and not
	// worry.
	tst.testQuick(t, fanout, &c, 0)
	tst.testQuick(t, fanout, &c, 0)
	tst.testQuick(t, fanout, &c, 0)

	diff := time.Since(start)
	if diff > (20 * time.Millisecond) {
		t.Errorf("Took too long sending to the fanout-sender. Took more than 20ms (%v). Should be ~instant.", diff)
	}
	time.Sleep(200 * time.Millisecond)
	if tst.received != 3 {
		t.Errorf("Didn't get the event after time expired? Wanted %d containers, got %d", 1, tst.received)
	}

	start = time.Now()
	tst.testQuick(t, fanout, &c, 0)
	tst.testQuick(t, fanout, &c, 0)
	tst.testQuick(t, fanout, &c, 0)

	diff = time.Since(start)
	if diff > (20 * time.Millisecond) {
		t.Errorf("Took too long sending to the fanout-sender. Took more than 20ms (%v). Should be ~instant.", diff)
	}
	// Should block, have 1 received as one worker clears up.
	tst.testQuick(t, fanout, &c, 1)

	diff = time.Since(start)
	if diff < (100 * time.Millisecond) {
		t.Errorf("Unexpectedly fast. Expected to block.Took less than 100ms (%v).", diff)
	}
	time.Sleep(200 * time.Millisecond)
	// Now all should be done...
	if tst.received != 4 {
		t.Errorf("Fanout: Expected 4 received events after timer(s) expired. Got %d", tst.received)
	}
}
