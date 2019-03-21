/*
 * skogul, debug sender
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

package senders

import (
	"encoding/json"
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"math/rand"
	"time"
)

/* Debug sender simply prints the metrics in json-marshalled format to
 * stdout.
 */
type Debug struct {
}

func (db Debug) Send(c *skogul.Container) error {
	b, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		log.Panic("Unable to marshal json for debug output: %s", err)
		return err
	}
	log.Printf("Debug: \n%s", b)
	return nil
}

/* The Sleeper-sender injects a random delay between 0 and MaxDelay before
 * passing execution over to the Next sender.
 *
 * The purpose is testing.
 */
type Sleeper struct {
	Next     skogul.Sender
	MaxDelay time.Duration
	Verbose  bool
}

func (sl Sleeper) Send(c *skogul.Container) error {
	d := rand.Float64() * float64(sl.MaxDelay)
	if sl.Verbose {
		log.Printf("Sleeping for %v", time.Duration(d))
	}
	time.Sleep(time.Duration(d))
	return sl.Next.Send(c)
}

/* The Counter sender emits, periodically, the flow-rate of metrics through
 * it.
 */
type Counter struct {
	Next    skogul.Sender
	last    time.Time
	metrics int64
}

func (co *Counter) Send(c *skogul.Container) error {
	mets := int64(0)
	for _, m := range c.Metrics {
		for range m.Data {
			mets += 1
		}
	}
	co.metrics += mets
	now := time.Now()
	if now.Sub(co.last) > (5 * time.Second) {
		if !co.last.IsZero() {
			log.Printf("Counted %f metrics/s (duration: %v)", float64(co.metrics)*float64(time.Second)/float64(now.Sub(co.last)), now.Sub(co.last))
		}
		co.metrics = 0
		co.last = now
	}
	return co.Next.Send(c)
}
