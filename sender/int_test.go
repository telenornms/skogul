/*
 * skogul, internal sender tests
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
)

type testSender struct {
	received int
}

func (ts *testSender) Send(c *skogul.Container) error {
	ts.received++
	return nil
}

func TestNull(t *testing.T) {
	c := skogul.Container{}
	x := sender.Null{}

	err := x.Send(&c)
	if err != nil {
		t.Errorf("Debug.Send returned non-nil: %v", err)
	}
}

func TestDupe(t *testing.T) {
	c := skogul.Container{}
	one := &(testSender{})
	two := &(testSender{})
	dupe := sender.Dupe{Next: []skogul.Sender{one, two}}

	err := dupe.Send(&c)
	if err != nil {
		t.Errorf("dupe.Send() failed: %v", err)
	}
	if one.received != 1 {
		t.Errorf("dupe.Send(), sender 1 expected %d recevied, got %d", 1, one.received)
	}
	if two.received != 1 {
		t.Errorf("dupe.Send(), sender 2 expected %d recevied, got %d", 1, two.received)
	}
}

func TestFallback(t *testing.T) {
	c := skogul.Container{}
	one := &(testSender{})
	two := &(testSender{})
	fb := sender.Fallback{Next: []skogul.Sender{one, two}}

	err := fb.Send(&c)
	if err != nil {
		t.Errorf("fallback.Send() failed: %v", err)
	}
	if one.received != 1 {
		t.Errorf("fallback.Send(), sender 1 expected %d recevied, got %d", 1, one.received)
	}
	if two.received != 0 {
		t.Errorf("fallback.Send(), sender 2 expected %d recevied, got %d", 0, two.received)
	}
}

func TestFallback_fail(t *testing.T) {
	c := skogul.Container{}
	one := &(testSender{})
	two := &(testSender{})
	three := &(testSender{})
	faf := &(sender.ForwardAndFail{Next: one})
	fb := sender.Fallback{Next: []skogul.Sender{faf, two}}

	fb.Add(three)

	err := fb.Send(&c)
	if err != nil {
		t.Errorf("fallback.Send() failed: %v", err)
	}
	if one.received != 1 {
		t.Errorf("fallback.Send(), sender 1 expected %d recevied, got %d", 1, one.received)
	}
	if two.received != 1 {
		t.Errorf("fallback.Send(), sender 2 expected %d recevied, got %d", 1, two.received)
	}
	if three.received != 0 {
		t.Errorf("fallback.Send(), sender 3 expected %d recevied, got %d", 0, two.received)
	}
}

func TestForwardAndFail(t *testing.T) {
	c := skogul.Container{}
	one := &(testSender{})
	faf := sender.ForwardAndFail{Next: one}

	err := faf.Send(&c)
	if err == nil {
		t.Errorf("forwardandfail.Send() .... failed to fail (returned true)")
	}
	if one.received != 1 {
		t.Errorf("forwardandfail.Send(), sender 1 expected %d recevied, got %d", 1, one.received)
	}
}
