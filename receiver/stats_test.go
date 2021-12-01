/*
 * skogul, stats receiver tests
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
 * Author(s):
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.no>
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
	"context"
	"testing"
	"time"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/receiver"
	"github.com/telenornms/skogul/sender"
)

func genStatsHandler(tester *sender.Test) *skogul.HandlerRef {
	return &skogul.HandlerRef{
		H: &skogul.Handler{
			Sender: tester,
		},
	}
}

func generateMetric() *skogul.Metric {
	now := skogul.Now()
	d := make(map[string]interface{})
	m := make(map[string]interface{})
	m["key"] = "example"
	d["value"] = 1
	return &skogul.Metric{
		Time:     &now,
		Data:     d,
		Metadata: m,
	}
}

func TestNoStatsReceived(t *testing.T) {
	tester := sender.Test{}
	h := genStatsHandler(&tester)
	stats := receiver.Stats{
		Interval: skogul.Duration{
			Duration: time.Millisecond * 10,
		},
		Handler: h,
	}

	skogul.StatsChan = make(chan *skogul.Metric, 2)
	defer func() {
		close(skogul.StatsChan)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), stats.Interval.Duration*2)
	defer cancel()
	go stats.StartC(ctx)

	// Allow stats to attempt to send
	time.Sleep(2 * stats.Interval.Duration)

	if tester.Received() != 0 {
		t.Errorf("expected to have gotten 0 stats containers but got %d", tester.Received())
	}
}

func TestStatsReceived(t *testing.T) {
	tester := sender.Test{}
	h := genStatsHandler(&tester)
	stats := receiver.Stats{
		Interval: skogul.Duration{
			Duration: time.Millisecond * 10,
		},
		Handler: h,
	}

	skogul.StatsChan = make(chan *skogul.Metric, 2)
	defer func() {
		close(skogul.StatsChan)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), stats.Interval.Duration*2)
	defer cancel()
	go stats.StartC(ctx)

	skogul.StatsChan <- generateMetric()

	// Allow stats to attempt to send
	time.Sleep(2 * stats.Interval.Duration)

	if tester.Received() != 1 {
		t.Errorf("expected to have gotten 1 stats container but got %d", tester.Received())
	}
}

func TestStatsDoesntBlockChan(t *testing.T) {
	tester := sender.Test{}
	h := genStatsHandler(&tester)
	stats := receiver.Stats{
		Interval: skogul.Duration{
			Duration: time.Millisecond * 10,
		},
		Handler: h,
	}

	skogul.StatsChan = make(chan *skogul.Metric, 2)
	defer func() {
		close(skogul.StatsChan)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), stats.Interval.Duration*2)
	defer cancel()
	go stats.StartC(ctx)

	t0 := time.Now()
	for i := 0; i < 100; i++ {
		skogul.StatsChan <- generateMetric()
	}
	td := time.Since(t0)

	// Allow stats to attempt to send
	time.Sleep(2 * stats.Interval.Duration)

	if tester.Received() != 1 {
		t.Errorf("expected to have gotten 1 stats container but got %d", tester.Received())
	}

	if td > time.Millisecond*1 {
		t.Errorf("expected stats channel to not block noticeably, but had to wait %v", td)
	}
}

func TestStatsDoesntBlockChanWithNoConfiguredReceiver(t *testing.T) {
	skogul.StatsChan = make(chan *skogul.Metric, 2)
	defer func() {
		close(skogul.StatsChan)
	}()

	// This is called by init, but since it has already been cancelled by earlier tests, we
	// have to start it again.
	drainCtx, drainCancel := context.WithCancel(context.Background())
	go receiver.DrainStats(drainCtx)
	defer drainCancel()

	done := make(chan bool)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	go func(ctx context.Context) {
		// looping one more than channel capacity to be blocked
		// if the channel is not being drained
		for i := 0; i < cap(skogul.StatsChan)+1; i++ {
			select {
			case <-ctx.Done():
			case skogul.StatsChan <- generateMetric():
			}
		}

		done <- true
		return
	}(ctx)

	select {
	case <-ctx.Done():
		t.Errorf("expected clean exit from context but got '%v'", ctx.Err())
	case <-done:
		// we got the finished signal, all good
	}
}
