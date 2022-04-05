/*
 * skogul, enrichment transformer tests
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

package transformer_test

import (
	"testing"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/transformer"
)

func TestEnrich_config1(t *testing.T) {
	testConfOk(t, `
	{
		"transformers": {
			"ok": {
				"type": "enrich",
				"source": "docs/examples/payloads/enirch.json",
				"keys": ["key1"]
			}
		}
	}`)
}

func BenchmarkHash(b *testing.B) {
	e := transformer.Enrich{}
	e.Keys = []string{"key1", "key2"}
	e.MakeSeed()
	e.Stringify = true

	m := skogul.Metric{}
	m.Metadata = make(map[string]interface{})
	m.Metadata["key1"] = "a car"
	m.Metadata["key2"] = "lol kek bikes rock"

	for i := 0; i < b.N; i++ {
		e.Hash(m)
	}
}
