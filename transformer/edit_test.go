/*
 * skogul, edit transformer tests
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

package transformer_test

import (
	"testing"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/transformer"
)

func checkReplace(t *testing.T, m *skogul.Metric, field string, want interface{}) {
	t.Helper()
	if m.Metadata[field] != want {
		t.Errorf("Edit transformer failed to enforce rule for field \"%s\". Wanted \"%v\", got \"%v\"", field, want, m.Metadata[field])
	}
}

func TestReplace(t *testing.T) {
	metric := skogul.Metric{}
	metric.Metadata = make(map[string]interface{})
	metric.Metadata["foo"] = "original"
	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric}

	replace := transformer.Replace{
		Source:      "foo",
		Destination: "bar",
		Regex:       "igi",
		Replacement: "KJELL MAGNE BONDEVIK MED SMÅ BOKSTAVER UTEN MELLOMNAVN",
	}

	t.Logf("Container before transform:\n%v", c)
	err := replace.Transform(&c)
	if err != nil {
		t.Errorf("Metadata() returned non-nil err: %v", err)
	}
	t.Logf("Container after transform:\n%v", c)

	checkReplace(t, c.Metrics[0], "foo", "original")
	checkReplace(t, c.Metrics[0], "bar", "orKJELL MAGNE BONDEVIK MED SMÅ BOKSTAVER UTEN MELLOMNAVNnal")
}
