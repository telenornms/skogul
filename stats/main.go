/*
 * skogul, internal stats
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

package stats

import (
	"context"
	"time"

	"github.com/telenornms/skogul"
)

var statsLog = skogul.Logger("stats", "main")

// StatsChan is a channel which accepts skogul statistic as a skogul.Metric
// By configuring the stats receiver, this channel is drained and sent on to
// the specified handler.
var StatsChan chan *skogul.Metric

// StatsDrainCtx and StatsDrainCancel are the context and cancel functions
// for the automatically created stats.StatsChan.
// If a skogul stats receiver is configured, StatsDrainCancel MUST be called
// so that statistics are not discarded.
var StatsDrainCtx, StatsDrainCancel = context.WithCancel(context.Background())

// init makes sure that the skogul stats channel exists at all times.
// Furthermore, it starts a goroutine to empty the channel in the case
// that the stats receiver is not configured, in which case the chan
// would end up blocking after it is filled.
func init() {
	// Create stats.StatsChan so we don't have components blocking on it
	if StatsChan == nil {
		StatsChan = make(chan *skogul.Metric, 100)
	}
	go DrainStats(StatsDrainCtx)
}

// drainStats drains all statistics on the stats channel.
// If the passed context is cancelled it will stop draining the channel
// so that a configured stats-receiver can listen on the channel.
func DrainStats(ctx context.Context) {
	statsLog.Debug("Starting stats drain. All stats are being dropped.")
	for {
		select {
		case <-StatsChan:
			continue
		case <-ctx.Done():
			statsLog.Debug("Stopping stats drain. Stats are being consumed.")
			return
		}
	}
}

// DefaultInterval is the default interval used for sending stats.
var DefaultInterval = time.Second * 10
