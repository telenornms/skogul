/*
 * skogul, mnr parser
 *
 * Copyright (c) 2020 Telenor Norge AS
 * Author(s):
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.com>
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
	"github.com/telenornms/skogul/parser"
	"github.com/telenornms/skogul/transformer"
)

func TestParseTransformer(t *testing.T) {
	b := []byte(`{"metrics":[{
        "time": "2020-01-01T13:33:37Z",
        "metadata": {
            "foo": "bar"
        },
        "data": {
            "influxql": "measurement tag=name field=name 1"
        }
    }]}`)

	j := parser.SkogulJSON{}
	c, err := j.Parse(b)
	if err != nil {
		t.Errorf("Failed to parse container")
		return
	}

	parserTransformer := transformer.Parse{
		Source: "influxql",
		Parser: skogul.ParserRef{
			P:    parser.InfluxDB{},
			Name: "influx",
		},
		Keep: true,
	}

	if err := parserTransformer.Transform(c); err != nil {
		t.Errorf("Failed to run parser transformer on container")
		return
	}

	data := c.Metrics[0].Data["influxql_data"].(map[string]interface{})
	if data["field"] != "name" {
		t.Errorf("Expected to find 'name' for 'influxql'.'field' in container after running parser")
	}
	if c.Metrics[0].Data["influxql"] != "measurement tag=name field=name 1" {
		t.Error("Expected to still find the source field for the data to parse in the data")
	}
}

func TestParseTransformerAppend(t *testing.T) {
	b := []byte(`{"metrics":[{
        "time": "2020-01-01T13:33:37Z",
        "metadata": {
            "foo": "bar"
        },
        "data": {
            "some_field": "foo",
            "influxql": "measurement tag=name field=name 1"
        }
    }]}`)

	j := parser.SkogulJSON{}
	c, err := j.Parse(b)
	if err != nil {
		t.Errorf("Failed to parse container")
		return
	}

	parserTransformer := transformer.Parse{
		Source: "influxql",
		Parser: skogul.ParserRef{
			P:    parser.InfluxDB{},
			Name: "influx",
		},
		Append: true,
		Keep:   false,
	}

	if err := parserTransformer.Transform(c); err != nil {
		t.Errorf("Failed to run parser transformer on container")
		return
	}

	if c.Metrics[0].Data["field"] != "name" {
		t.Errorf("Expected to find 'name' for 'field' in container after running parser")
	}
	if c.Metrics[0].Data["some_field"] != "foo" {
		t.Error("Expected to find 'foo' for 'some_field' in container data after running parser with append")
	}
	if c.Metrics[0].Data["influxql"] == "measurement tag=name field=name 1" {
		t.Error("Expected the 'influxql' field to be removed from the container when Keep: false")
	}
}
