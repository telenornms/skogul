/*
 * skogul, counter tests
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
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/sender"
	"testing"
	"time"
)

func TestCounter(t *testing.T) {
	one := &(sender.Test{})
	two := &(sender.Test{})

	h := skogul.Handler{Sender: one, Parser: parser.JSON{}}
	cnt := sender.Counter{Next: two, Stats: h, Period: time.Duration(50 * time.Millisecond)}
	two.TestQuick(t, &cnt, &validContainer, 1)
	if one.Received() != 0 {
		t.Errorf("Stats received too early.")
	}
	two.TestQuick(t, &cnt, &validContainer, 1)
	two.TestQuick(t, &cnt, &validContainer, 1)
	two.TestQuick(t, &cnt, &validContainer, 1)
	if one.Received() != 0 {
		t.Errorf("Stats received too early 2.")
	}
	time.Sleep(time.Duration(50 * time.Millisecond))
	if one.Received() != 0 {
		t.Errorf("Stats received too early 2.")
	}
	two.TestQuick(t, &cnt, &validContainer, 1)
	if one.Received() != 1 {
		t.Errorf("Correct stats not received")
	}
	two.TestQuick(t, &cnt, &validContainer, 1)
	two.TestQuick(t, &cnt, &validContainer, 1)
	time.Sleep(time.Duration(50 * time.Millisecond))
	two.TestQuick(t, &cnt, &validContainer, 1)
	if one.Received() != 2 {
		t.Errorf("Correct stats not received")
	}
}
