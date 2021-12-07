/*
 * skogul, stats receiver tests
 *
 * Copyright (c) 2021 Telenor Norge AS
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
	"github.com/telenornms/skogul/stats"
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
	statsReceiver := receiver.Stats{
		Handler: h,
	}

	stats.Chan = make(chan *skogul.Metric, 2)
	defer func() {
		close(stats.Chan)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*20)
	defer cancel()
	go statsReceiver.StartC(ctx)

	// Allow stats to attempt to send
	time.Sleep(time.Millisecond * 20)

	if tester.Received() != 0 {
		t.Errorf("expected to have gotten 0 stats containers but got %d", tester.Received())
	}
}

func TestStatsReceived(t *testing.T) {
	tester := sender.Test{}
	h := genStatsHandler(&tester)
	statsReceiver := receiver.Stats{
		Handler: h,
	}

	stats.Chan = make(chan *skogul.Metric, 2)
	defer func() {
		close(stats.Chan)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*20)
	defer cancel()
	go statsReceiver.StartC(ctx)

	stats.Chan <- generateMetric()

	// Allow stats to attempt to send
	time.Sleep(time.Millisecond * 20)

	if tester.Received() != 1 {
		t.Errorf("expected to have gotten 1 stats container but got %d", tester.Received())
	}
}

func TestStatsDoesntBlockChan(t *testing.T) {
	tester := sender.Test{}
	h := genStatsHandler(&tester)
	statsReceiver := receiver.Stats{
		Handler: h,
	}

	stats.Chan = make(chan *skogul.Metric, 2)
	defer func() {
		close(stats.Chan)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*20)
	defer cancel()
	go statsReceiver.StartC(ctx)

	t0 := time.Now()
	for i := 0; i < 100; i++ {
		stats.Chan <- generateMetric()
	}
	td := time.Since(t0)

	// Allow stats to attempt to send
	time.Sleep(time.Millisecond * 20)

	if tester.Received() != 100 {
		t.Errorf("expected to have gotten 100 stats container but got %d", tester.Received())
	}

	if td > time.Millisecond*1 {
		t.Errorf("expected stats channel to not block noticeably, but had to wait %v", td)
	}
}

func TestStatsDoesntBlockChanWithNoConfiguredReceiver(t *testing.T) {
	stats.Chan = make(chan *skogul.Metric, 2)
	defer func() {
		close(stats.Chan)
	}()

	// This is called by init, but since it has already been cancelled by earlier tests, we
	// have to start it again.
	drainCtx, drainCancel := context.WithCancel(context.Background())
	go stats.DrainStats(drainCtx)
	defer drainCancel()

	done := make(chan bool)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	go func(ctx context.Context) {
		// looping one more than channel capacity to be blocked
		// if the channel is not being drained
		for i := 0; i < cap(stats.Chan)+1; i++ {
			select {
			case <-ctx.Done():
			case stats.Chan <- generateMetric():
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
