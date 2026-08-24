/*
 * skogul, backoff sender
 *
 * Copyright (c) 2019-2026 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *
 * This library is free software; you can redistribute it and/or modify it
* under the terms of the GNU Lesser General Public License as published by the
* Free Software Foundation; either version 2.1 of the License, or (at your
* option) any later version.
 *
* This library is distributed in the hope that it will be useful, but WITHOUT
* ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
* FOR A PARTICULAR PURPOSE.  See the GNU Lesser General Public License for more
* details.
 *
* You should have received a copy of the GNU Lesser General Public License
* along with this library; if not, write to the Free Software Foundation, Inc.,
* 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301  USA
*/

package sender

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/telenornms/skogul"
)

var backoffLog = skogul.Logger("sender", "backoff")

// Backoff sender will send to Next, making up to Retries attempts, with
// the delay doubling after each failed attempt. It also keeps a holdoff
// count across calls: while a previous send has failed and not yet been
// made up for by a success, each Send waits Base before its first
// attempt.
//
// Deprecated: The 'backoff' sender wrapper is deprecated and will be removed
// in a future version. Use the built-in retry config on each network sender
// with backoff/retry support instead.
type Backoff struct {
	Next    skogul.SenderRef `doc:"The sender to try"`
	Base    skogul.Duration  `doc:"Delay after a failure. Doubles for each further attempt."`
	Retries uint64           `doc:"Total number of attempts before giving up, not counted on top of the first one. A value of 1 means a single attempt and no retry."`
	holdoff atomic.Uint64
}

// Deprecated returns an error for the Deprecated interface.
func (bo *Backoff) Deprecated() error {
	return fmt.Errorf("the 'backoff' sender wrapper is deprecated and will be removed in a future version - network senders now use backoff by default")
}

// Verify checks the configuration. Retries is a total attempt count, so
// zero would mean "never send at all", which Send has no way to report:
// it would return success without ever reaching Next.
func (bo *Backoff) Verify() error {
	if bo.Retries == 0 {
		return fmt.Errorf("retries is the total number of attempts, not a count on top of the first one, so it must be at least 1 - use 1 for a single attempt with no retry")
	}
	if bo.Base.Duration <= 0 {
		backoffLog.Warn("Base is unset, so a failed send is retried with no delay at all")
	}
	return nil
}

// Send with a delay
func (bo *Backoff) Send(c *skogul.Container) error {
	var err error
	delay := bo.Base.Duration
	t := bo.holdoff.Load()
	if t > 0 {
		time.Sleep(delay)
	}
	// Verify rejects a zero Retries, but a Backoff built in code never
	// passes through Verify, and silently dropping the container would
	// be worse than ignoring the setting.
	attempts := bo.Retries
	if attempts == 0 {
		attempts = 1
	}
	for i := uint64(1); i <= attempts; i++ {
		err = bo.Next.S.Send(c)
		if err == nil {
			if i > 1 {
				bo.holdoff.Add(1 - i)
			}
			return nil
		}
		bo.holdoff.Add(1)
		time.Sleep(delay)
		delay = delay * 2
	}
	return err
}
