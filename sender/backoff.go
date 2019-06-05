/*
 * skogul, backoff sender
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

package sender

import (
	"github.com/KristianLyng/skogul"
	"time"
)

// Backoff sender will send to Next, but retry up to Retries times, with
// exponential backoff, starting with time.Duration
type Backoff struct {
	Next    skogul.Sender
	Base    time.Duration
	Retries int
}

// Send with a delay
func (bo *Backoff) Send(c *skogul.Container) error {
	var err error
	delay := bo.Base
	for i := 1; i <= bo.Retries; i++ {
		err = bo.Next.Send(c)
		if err == nil {
			return nil
		}
		time.Sleep(delay)
		delay = delay * 2
	}
	return err
}
