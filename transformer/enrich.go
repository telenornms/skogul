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
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/telenornms/skogul"
	"hash/maphash"
	"io/ioutil"
	"sync"
)

// keyType identifies the type the hash returns
type keyType uint64

// entry represents what is supposed to be added. An entry is added to a
// metric's metadatafields on a hit.
type entry map[string]interface{}

/*
Enrich transformer adds additional metadatainformation to metrics. This
is a work in progress, see #224 for on-going discussions.

Some items that need to be resolved:
- Need benchmarks
- This hash thing is just weird, maybe just use sha512 instead.
- I think we _have_ to stringify. A weird example is that float 0 and
  int 0 is binary equivalent, which means that if the source is parsed as
  float and the metric as int, most numbers wont match, but 0 will, which
  is just confusing.
- I don't like the buffer thing at all, but it might be needed I suppose.
  Need to double check how performant that is for the typical scenario, and
  possibly make it configurable.
- Need to warn if the source has more metadata fields than we care about, and fail if has fewer.
- Need to be able to reload, preferably in the background.
- We might need to add a Init interface here to give the transformer a
  chance to load before we start accepting data. All once.Do() works OK for
  things that are fast, but this could possibly be hundreds of megabytes of
  JSON that needs to be read and parsed and hashed.
- Yeah, benchmarks.
*/
type Enrich struct {
	Keys      []string `doc:"Metadatafields to match"`
	Stringify bool     `doc:"Set to true to convert all metadata to text prior to hashing. Doesn't change actual data, but might simplify matching numbers that might be binary encoded to different formats, at the cost of a slight performance hit."`
	Source    string   `doc:"Path to enrichment-file."`
	ok        bool
	once      sync.Once
	store     map[keyType]*entry
	seed      maphash.Seed
}

var eLog = skogul.Logger("transformer", "enrich")

// Hash calculates the hash for a metric using relevant settings (e.g.:
// e.Keys). Only exported for the sake of tests and benchmarks.
//
// FIXME: What to do if stringify fails? Logging isn't really a solution.
// We risk matching the wrong thing.
func (e *Enrich) Hash(m skogul.Metric) keyType {
	var h maphash.Hash
	h.SetSeed(e.seed)
	buf := new(bytes.Buffer)
	for _, meta := range e.Keys {
		if m.Metadata[meta] == nil {
			continue
		}
		var err error
		if e.Stringify {
			err = binary.Write(buf, binary.LittleEndian, []byte(fmt.Sprintf("%v", m.Metadata[meta])))
		} else {
			err = binary.Write(buf, binary.LittleEndian, m.Metadata[meta])
		}
		if err != nil {
			eLog.WithError(err).Warnf("Failed to write binary representation to buffer for %s", meta)
			continue
		}
	}
	h.Write(buf.Bytes())
	ret := h.Sum64()
	return keyType(ret)
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

func (e *Enrich) load() error {
	ms := []skogul.Metric{}
	eLog.Debugf("Loading: Reading file")
	content, err := ioutil.ReadFile(e.Source)
	if err != nil {
		return fmt.Errorf("unable to load enrichment from %s: %w", e.Source, err)
	}

	eLog.Debugf("Loading: Umarshaling")
	err = json.Unmarshal(content, &ms)
	if err != nil {
		return fmt.Errorf("unable to load enrichment data from %s, parsing JSON failed: %w", e.Source, err)
	}
	eLog.Debugf("Loading: Populating")
	for _, m := range ms {
		en := entry(m.Data)
		h := e.Hash(m)
		if e.store[h] != nil {
			eLog.Warnf("Hash collision while adding item %v! Overwriting!", m)
		}
		e.store[h] = &en
	}
	eLog.Debugf("Loading: Done")
	return nil
}

// MakeSeed() sets up the seed for an enricher, internal, only exported for
// tests.
func (e *Enrich) MakeSeed() {
	e.seed = maphash.MakeSeed()
}

func (e *Enrich) Transform(c *skogul.Container) error {
	e.once.Do(func() {
		eLog.Warnf("The enrichment transformer is in use. This transformer is highly experimental and not considered production ready. Functionality and configuration parameters will change as it matures. If you use this, PLEASE provide feedback on what your use cases require.")
		e.ok = false
		e.MakeSeed()
		e.store = make(map[keyType]*entry)
		err := e.load()
		if err != nil {
			eLog.WithError(err).Error()
		} else {
			e.ok = true
		}
	})

	for _, m := range c.Metrics {
		hit := e.find(*m)
		if hit == nil {
			continue
		}
		for idx, f := range *hit {
			m.Metadata[idx] = f
		}
	}
	return nil
}
