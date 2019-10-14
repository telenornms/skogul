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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/config"
	"github.com/KristianLyng/skogul/transformer"
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

	metadata := transformer.Metadata{
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

func TestMetadata_config(t *testing.T) {
	testConfOk(t, `
	{
		"transformers": {
			"ok": {
				"type": "metadata",
				"set": {
					"this": "to that",
					"foo": "is bar"
				},
				"require": [ "reqthis" ],
				"remove": [ "gruff" ],
				"ban": [ "trash" ]
			}
		}
	}`)
	testConfOk(t, `
	{
		"transformers": {
			"ok": {
				"type": "metadata",
				"set": { },
				"require": [],
				"remove": [ ],
				"ban": [  ]
			}
		}
	}`)
	testConfBad(t, `
	{
		"transformers": {
			"ok": {
				"type": "metadata",
				"set": 5
				"require": [ "reqthis" ],
				"remove": [ "gruff" ],
				"ban": [ "trash" ]
			}
		}
	}`)
}

func testConfOk(t *testing.T, rawconf string) {
	t.Helper()
	conf, err := config.Bytes([]byte(rawconf))
	if err != nil {
		t.Errorf("failed to parse valid transformer config: %v", err)
	}
	if conf == nil {
		t.Errorf("failed to get valid config for transformer")
	}
}

func testConfBad(t *testing.T, rawconf string) {
	t.Helper()
	conf, err := config.Bytes([]byte(rawconf))
	if err == nil {
		t.Errorf("Didn't catch invalid config")
	}
	if conf != nil {
		t.Errorf("got config from invalid source config")
	}
}

func TestExtract(t *testing.T) {
	extracted_value_key := "extract-this"
	extracted_value := "the value"

	metric := skogul.Metric{}
	metric.Metadata = make(map[string]interface{})
	metric.Metadata["extract"] = extracted_value_key

	metric.Data = make(map[string]interface{})
	testData := fmt.Sprintf(`{"%s": "%s"}`, extracted_value_key, extracted_value)
	json.Unmarshal([]byte(testData), &metric.Data)

	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric}

	metadata := transformer.Metadata{
		Extract: [][]string{[]string{extracted_value_key}},
	}

	err := metadata.Transform(&c)

	if err != nil {
		t.Error(err)
	}

	if c.Metrics[0].Metadata[extracted_value_key] != extracted_value {
		t.Errorf(`Expected %s but got %s`, extracted_value, c.Metrics[0].Metadata[extracted_value_key])
	}
}
