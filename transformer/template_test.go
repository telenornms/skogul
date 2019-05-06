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

package transformer_test

import (
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/transformer"
	"testing"
	"time"
)

func TestTemplate(t *testing.T) {
	now := time.Now()

	metric := skogul.Metric{}
	metric.Data = make(map[string]interface{})
	metric.Data["test"] = "foo"

	template := skogul.Metric{}
	template.Time = &now
	template.Metadata = make(map[string]interface{})
	template.Metadata["test"] = "foo"

	c := skogul.Container{}
	c.Metrics = []skogul.Metric{metric}
	c.Template = &template

	templater := transformer.Templater{}

	err := templater.Transform(&c)

	if err != nil {
		t.Errorf("Templater() returned non-nil err: %v", err)
	}
	if c.Metrics[0].Time != &now {
		t.Errorf("Templater() did not expand time. Wanted %v got %v", &now, c.Metrics[0].Time)
	}
	if c.Metrics[0].Metadata["test"] != "foo" {
		t.Errorf("Templater() did not expand Metadata correctly. Wanted Metadata[\"test\"]=\"foo\", got Metadata[\"test\"]=%v", c.Metrics[0].Metadata["test"])
	}
	if c.Template != nil {
		t.Errorf("Templater() left template intact: %v", c.Template)
	}
}
