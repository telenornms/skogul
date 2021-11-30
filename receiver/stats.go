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
	"time"

	"github.com/telenornms/skogul"
)

var statsLog = skogul.Logger("receiver", "stats")

// Stats sends metrics to a HTTP endpoint
type Stats struct {
	Handler           *skogul.HandlerRef
	Interval          skogul.Duration
	SendEveryInterval bool `doc:"Send stats for every configured interval, even if no new stats are to be reported."` //FIXME: Skogul crashes on sending empty metrics
	ch                chan *skogul.Metric
	ticker            *time.Ticker
	//ch       chan *skogul.Container // FIXME: no *  (?), we should own this // global chan ?
}

// Starts starts listening for Skogul stats and
// emits them on the configured interval.
func (s *Stats) Start() error {
	if s.Interval.Duration == 0 {
		statsLog.Debug("Missing interval for stats reporting, defaulting to every 3 seconds")
		s.Interval.Duration = 3 * time.Second
	}

	if skogul.StatsChan == nil {
		skogul.StatsChan = make(chan *skogul.Metric, 1000)
	}

	s.ch = make(chan *skogul.Metric, 100)

	s.ticker = time.NewTicker(s.Interval.Duration)

	go s.runner()

	for metric := range skogul.StatsChan {
		s.ch <- metric
	}
	return nil
}

// runner is the function listening for stats and emits
// them when there is time for it.
func (s *Stats) runner() {
	for range s.ticker.C {
		statsLog.Trace("Time to send skogul stats")
		// Create a new channel,
		// consume all on the existing, and then close it
		ch := s.ch
		s.ch = make(chan *skogul.Metric, 100)
		close(ch)

		metrics := make([]*skogul.Metric, 0)

		// Drain the old channel
		for {
			statsLog.Trace("Consuming entry from channel..")
			metric, more := <-ch
			if !more {
				break
			} else if metric == nil {
				statsLog.Error("Got nil metric on stats channel with more metrics left. Discarding this.")
				break
			}
			metrics = append(metrics, metric)
		}

		if len(metrics) == 0 && !s.SendEveryInterval {
			// We have no metrics and we're configured
			// to *not* ship the metrics if we have none,
			// so we wait until next tick.
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
