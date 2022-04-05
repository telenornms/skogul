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

// Enrich transformer adds additional metadatainformation to metrics.
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

func (e *Enrich) hash(m skogul.Metric) keyType {
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
			eLog.Tracef("Field %v is %T", meta, m.Metadata[meta])
			err = binary.Write(buf, binary.LittleEndian, m.Metadata[meta])
		}
		if err != nil {
			eLog.WithError(err).Warnf("Failed to write binary representation to buffer")
			continue
		}
	}
	h.Write(buf.Bytes())
	ret := h.Sum64()
	eLog.Tracef("Calculated sum: %d - seed %v\n", ret, h.Seed())
	return keyType(ret)
}

func (e *Enrich) find(m skogul.Metric) *entry {
	nm := e.store[e.hash(m)]
	if nm != nil {
		eLog.Tracef("Enrichment hit")
	} else {
		eLog.Tracef("Enrichment miss")
	}
	return nm
}

func (e *Enrich) load() error {
	ms := []skogul.Metric{}
	content, err := ioutil.ReadFile(e.Source)
	if err != nil {
		return fmt.Errorf("unable to load enrichment from %s: %w", e.Source, err)
	}

	err = json.Unmarshal(content, &ms)
	if err != nil {
		return fmt.Errorf("unable to load enrichment data from %s, parsing JSON failed: %w", e.Source, err)
	}
	for _, m := range ms {
		eLog.Tracef("adding m %v", m)
		en := entry(m.Data)
		h := e.hash(m)
		if e.store[h] != nil {
			eLog.Warnf("Hash collision while adding item %v! Overwriting!", m)
		}
		e.store[h] = &en
	}
	return nil
}

func (e *Enrich) Transform(c *skogul.Container) error {
	e.once.Do(func() {
		eLog.Warnf("The enrichment transformer is in use. This transformer is highly experimental and not considered production ready. Functionality and configuration parameters will change as it matures. If you use this, PLEASE provide feedback on what your use cases require.")
		e.ok = false
		e.seed = maphash.MakeSeed()
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
