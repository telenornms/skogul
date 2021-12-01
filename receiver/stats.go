/*
 * skogul, stats receiver
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
	Interval skogul.Duration
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
	if s.Interval.Duration == 0 {
		statsLog.Debugf("Missing interval for stats reporting, defaulting to every %d seconds", stats.DefaultInterval)
		s.Interval.Duration = stats.DefaultInterval
	}
	if s.ChanSize == 0 {
		s.ChanSize = 100
	}

	// XXX: we shouldn't allow multiple stats receivers instantiated probably
	// because they'll just steal the stats from each other.
	// or we'll have to direct stats to a specific stats instance.

	s.ch = make(chan *skogul.Metric, s.ChanSize)

	s.ticker = time.NewTicker(s.Interval.Duration)

	go s.runner()

	stats.StatsDrainCancel()

	for metric := range stats.StatsChan {
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

// runner is the function listening for stats and emits
// them when there is time for it.
func (s *Stats) runner() {
	for range s.ticker.C {
		statsLog.WithField("stats", len(s.ch)).Trace("Time to send skogul stats")

		metrics := make([]*skogul.Metric, len(s.ch))

		// Drain the current messages on the channel
		for i := range metrics {
			metric, more := <-s.ch
			if !more {
				break
			} else if metric == nil {
				statsLog.Debug("Got nil metric on stats channel with more metrics left. Discarding this.")
				break
			}
			metrics[i] = metric
		}

		if len(metrics) == 0 {
			// We have no metrics so we wait until next tick.
			statsLog.Trace("Skipping sending metrics since we have none")
			continue
		}

		container := skogul.Container{
			Metrics: metrics,
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
