/*
 * skogul, stats receiver
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

package receiver

import (
	"context"
	"time"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/stats"
)

var statsLog = skogul.Logger("receiver", "stats")

// Stats receives metrics from skogul and forwards it to a handler.
type Stats struct {
	Handler  *skogul.HandlerRef
	ChanSize uint64
	ch       chan *skogul.Metric
	ticker   *time.Ticker
}

// Start starts listening for Skogul stats and
// emits them on the configured interval.
func (s *Stats) Start() error {
	return s.StartC(context.Background())
}

// StartC allows starting Stats with a context.
func (s *Stats) StartC(ctx context.Context) error {
	if s.ChanSize == 0 {
		s.ChanSize = 100
	}

	// XXX: we shouldn't allow multiple stats receivers instantiated probably
	// because they'll just steal the stats from each other.
	// or we'll have to direct stats to a specific stats instance.

	s.ch = make(chan *skogul.Metric, s.ChanSize)

	go s.runner()

	stats.CancelDrain()

	for metric := range stats.Chan {
		if len(s.ch) >= cap(s.ch) {
			statsLog.Debug("Dropping stats because the channel is full")
			continue
		}
		select {
		case s.ch <- metric:
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// runner is the function listening for stats and
// emits them through the configured handler
func (s *Stats) runner() {
	for metric := range s.ch {
		statsLog.WithField("stats", len(s.ch)).Trace("Time to send skogul stats")
		container := skogul.Container{
			Metrics: []*skogul.Metric{metric},
		}

		if err := s.Handler.H.Send(&container); err != nil {
			statsLog.WithError(err).Error("Failed to send skogul stats")
		}
	}
}

// Verify makes sure all required parameters are set
func (s *Stats) Verify() error {
	return nil
}
