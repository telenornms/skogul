/*
 * skogul, batch sender - tests
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

func TestBatch(t *testing.T) {
	c := skogul.Container{}
	m := skogul.Metric{}
	// batcher doesn't really worry about the internals of metrics,
	// so we just leave them blank and reuse the same metric.

	c.Metrics = []*skogul.Metric{&m}
	one := &(testSender{})
	batch := &(sender.Batch{Next: one})

	// Test that sending 9 metrics doesn't pass anything on
	for i := 0; i < 9; i++ {
		one.testQuick(t, batch, &c, 0)
	}

	// but the 10th does....
	one.testQuick(t, batch, &c, 1)

	// Rinse and repeat to ensure state is reset
	for i := 0; i < 9; i++ {
		one.testQuick(t, batch, &c, 0)
	}
	one.testQuick(t, batch, &c, 1)

	// Test that a single metric wont be passed on instantly...
	one.testQuick(t, batch, &c, 0)

	// but that it will after the timer expires
	time.Sleep(time.Duration(1 * time.Second))
	if one.received != 1 {
		t.Errorf("batch.Send(), no data sent after timeout expired. Expected %d, got %d", 1, one.received)
	}

	// Ensure we don't botch the resize - send containers with multiple
	// metrics and see how that works.
	c.Metrics = []*skogul.Metric{&m, &m, &m, &m, &m, &m, &m, &m, &m}
	one.testQuick(t, batch, &c, 0)
	one.testQuick(t, batch, &c, 1)
}
