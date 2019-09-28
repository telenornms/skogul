/*
 * skogul, templating tests
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

package transformer

import (
	"github.com/KristianLyng/skogul"
	"testing"
)

func check(t *testing.T, m *skogul.Metric, field string, want interface{}) {
	t.Helper()
	if m.Metadata[field] != want {
		t.Errorf("Metadata transformer failed to enforce rule for field \"%s\". Wanted \"%v\", got \"%v\"", field, want, m.Metadata[field])
	}
}
func TestMetadata(t *testing.T) {
	//now := time.Now()

	metric := skogul.Metric{}
	metric.Metadata = make(map[string]interface{})
	metric.Metadata["set"] = "original"
	metric.Metadata["require"] = "present"
	metric.Metadata["remove"] = "not removed"
	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric}

	metadata := Metadata{
		Set:     map[string]interface{}{"set": "new"},
		Require: []string{"require"},
		Remove:  []string{"remove"},
		Ban:     []string{"ban"},
	}

	err := metadata.Transform(&c)

	if err != nil {
		t.Errorf("Metadata() returned non-nil err: %v", err)
	}

	check(t, c.Metrics[0], "set", "new")
	check(t, c.Metrics[0], "require", "present")
	check(t, c.Metrics[0], "remove", nil)
	check(t, c.Metrics[0], "ban", nil)
}
