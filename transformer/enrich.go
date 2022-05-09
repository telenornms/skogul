/*
 * skogul, enrichment transformer
 *
 * Copyright (c) 2022 Telenor Norge AS
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

package transformer

import (
	"encoding/json"
	"fmt"
	"github.com/telenornms/skogul"
	"io"
	"os"
	"sync"
)

// keyType identifies the type the hash returns
type keyType string

// entry represents what is supposed to be added. An entry is added to a
// metric's metadatafields on a hit.
type entry map[string]interface{}

/*
Enrich transformer adds additional metadatainformation to metrics. This
is a work in progress, see #224 for on-going discussions.

Some items that need to be resolved:
- Need benchmarks
- I think we _have_ to stringify. A weird example is that float 0 and
  int 0 is binary equivalent, which means that if the source is parsed as
  float and the metric as int, most numbers wont match, but 0 will, which
  is just confusing.
- Need to warn if the source has more metadata fields than we care about,
  and fail if has fewer.
- Need to be able to reload, preferably in the background.
- We might need to add a Init interface here to give the transformer a
  chance to load before we start accepting data. All once.Do() works OK for
  things that are fast, but this could possibly be hundreds of megabytes of
  JSON that needs to be read and parsed and hashed.
- Memory usage is a bit high, and it seems to be because of json decoding.
  Profiling and common sense suggest that the pure storage isn't too
  bloated, but the json decoder state lingers, so there might be some
  references there.
- Loading is way too slow. Need alternatives, but that might be different
  implementations as the simplicity of json is real nice.
*/
type Enrich struct {
	Keys   []string `doc:"Metadatafields to match"`
	Source string   `doc:"Path to enrichment-file."`
	lock   sync.RWMutex
	once   sync.Once
	store  map[keyType]*entry
}

var eLog = skogul.Logger("transformer", "enrich")

// Hash calculates the hash for a metric using relevant settings (e.g.:
// e.Keys). Only exported for the sake of tests and benchmarks.
//
// XXX: This used to be a manually created hash, now we rely on the
// underlying map[] implementation instead.
func (e *Enrich) Hash(m skogul.Metric) keyType {
	ret := ""
	for _, meta := range e.Keys {
		ret = ret + fmt.Sprintf("%v", m.Metadata[meta])
	}
	return keyType(ret)
}

func (e *Enrich) Update(c *skogul.Container) {
	for _, m := range c.Metrics {
		eLog.Infof("Updating metrics %v", m)
		e.save(m)
	}
}

func (e *Enrich) find(m skogul.Metric) *entry {
	nm := e.store[e.Hash(m)]
	if nm != nil {
		eLog.Tracef("Enrichment hit")
	} else {
		eLog.Tracef("Enrichment miss")
	}
	return nm
}
func (e *Enrich) save(m *skogul.Metric) {
	e.lock.Lock()
	h := e.Hash(*m)
	if e.store[h] != nil {
		eLog.Warnf("Hash collision while adding item %v! Overwriting!", m)
	}
	en := entry(m.Data)
	e.store[h] = &en
	e.lock.Unlock()
}

func (e *Enrich) load() error {

	eLog.Debugf("Loading: Reading file")
	f, err := os.Open(e.Source)
	if err != nil {
		return fmt.Errorf("unable to open enrichment from %s: %w", e.Source, err)
	}

	eLog.Debugf("Loading: Umarshaling")
	dec := json.NewDecoder(f)
	for {
		var m skogul.Metric
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Unable to decode JSON message in enrichment: %w", err)
		}
		e.save(&m)
	}
	eLog.Debugf("Loading: Done")
	return nil
}

// Transform looks up a metric in the stored enrichment database (loading
// it on initial run) and adds metadata if a match is found.
func (e *Enrich) Transform(c *skogul.Container) error {
	e.once.Do(func() {
		eLog.Warnf("The enrichment transformer is in use. This transformer is highly experimental and not considered production ready. Functionality and configuration parameters will change as it matures. If you use this, PLEASE provide feedback on what your use cases require.")
		e.store = make(map[keyType]*entry)
		err := e.load()
		if err != nil {
			eLog.WithError(err).Error()
		}
	})
	e.lock.RLock()
	for _, m := range c.Metrics {
		hit := e.find(*m)
		if hit == nil {
			continue
		}
		for idx, field := range *hit {
			m.Metadata[idx] = field
		}
	}
	e.lock.RUnlock()
	return nil
}
