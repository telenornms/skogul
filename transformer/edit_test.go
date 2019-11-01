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

func replaceMkContainer(foo, bar interface{}) *skogul.Container {
	metric := skogul.Metric{}
	metric.Metadata = make(map[string]interface{})
	metric.Metadata["foo"] = foo
	metric.Metadata["bar"] = bar
	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric}
	return &c
}

func TestReplace(t *testing.T) {
	c := replaceMkContainer("original", "no matter")

	replace := transformer.Replace{
		Source:      "foo",
		Destination: "bar",
		Regex:       "igi",
		Replacement: "KJELL MAGNE BONDEVIK MED SMÅ BOKSTAVER UTEN MELLOMNAVN",
	}

	t.Logf("Container before transform:\n%v", c)
	err := replace.Transform(c)
	if err != nil {
		t.Errorf("Metadata() returned non-nil err: %v", err)
	}
	t.Logf("Container after transform:\n%v", c)
	checkReplace(t, c.Metrics[0], "foo", "original")
	checkReplace(t, c.Metrics[0], "bar", "orKJELL MAGNE BONDEVIK MED SMÅ BOKSTAVER UTEN MELLOMNAVNnal")
}

func TestReplace_config1(t *testing.T) {
	conf := testConfOk(t, `
	{
		"transformers": {
			"ok": {
				"type": "replace",
				"source": "foo",
				"regex": "^(.*)(o+)$",
				"replacement": "$2² - $1"
			}
		}
	}`)

	container := replaceMkContainer("originalo", "unchanged")
	conf.Transformers["ok"].Transformer.Transform(container)
	checkReplace(t, container.Metrics[0], "foo", "o² - original")
	checkReplace(t, container.Metrics[0], "bar", "unchanged")
}

func TestReplace_config2(t *testing.T) {
	conf := testConfOk(t, `
	{
		"transformers": {
			"ok": {
				"type": "replace",
				"source": "foo",
				"destination": "bar",
				"regex": "kjeks",
				"replacement": "sjokoladekake"
			}
		}
	}`)

	container := replaceMkContainer("vaffelkjeks", "unchanged")
	conf.Transformers["ok"].Transformer.Transform(container)
	checkReplace(t, container.Metrics[0], "bar", "vaffelsjokoladekake")
	checkReplace(t, container.Metrics[0], "foo", "vaffelkjeks")

	container = replaceMkContainer("ikke match", "unchanged")
	conf.Transformers["ok"].Transformer.Transform(container)
	checkReplace(t, container.Metrics[0], "bar", "ikke match")
	checkReplace(t, container.Metrics[0], "foo", "ikke match")

	container = replaceMkContainer(nil, "unchanged")
	conf.Transformers["ok"].Transformer.Transform(container)
	checkReplace(t, container.Metrics[0], "bar", nil)
	checkReplace(t, container.Metrics[0], "foo", nil)

	x := make(map[string]interface{}, 1)
	x["kek"] = "lol"
	container = replaceMkContainer(x, "original")
	conf.Transformers["ok"].Transformer.Transform(container)
	checkReplace(t, container.Metrics[0], "bar", nil)

	container = replaceMkContainer("vaffelkjeks", "unchanged")
	container2 := replaceMkContainer("ikkematch", "unchanged")
	container2.Metrics[0].Metadata = nil
	container.Metrics = append(container.Metrics, container2.Metrics[0])
	conf.Transformers["ok"].Transformer.Transform(container)
	checkReplace(t, container.Metrics[0], "bar", "vaffelsjokoladekake")
	checkReplace(t, container.Metrics[0], "foo", "vaffelkjeks")
}

func TestReplace_BadConfig(t *testing.T) {
	testConfBad(t, `
	{
		"transformers": {
			"ok": {
				"type": "replace",
				"regex": "^(.*)(o+)$",
				"replacement": "$2² - $1"
			}
		}
	}`)
	testConfBad(t, `
	{
		"transformers": {
			"ok": {
				"type": "replace",
				"source": "foo",
				"regex": "^(.*(o+)$",
				"replacement": "$2² - $1"
			}
		}
	}`)
	testConfBad(t, `
	{
		"transformers": {
			"ok": {
				"type": "replace",
				"source": null,
				"regex": "^(.*)(o+)$",
				"replacement": "$2² - $1"
			}
		}
	}`)
	testConfBad(t, `
	{
		"transformers": {
			"ok": {
				"type": "replace",
				"source": "foo",
				"replacement": "$2² - $1"
			}
		}
	}`)
	testConfBad(t, `
	{
		"transformers": {
			"ok": {
				"type": "replace",
				"regex": "^(.*)(o+)$",
				"replacement": "$2² - $1"
			}
		}
	}`)
}
