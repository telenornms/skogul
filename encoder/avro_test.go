/*
 * skogul, test avro encoder
 *
 * Copyright (c) 2022 Telenor Norge AS
 * Author:
 *  - Roshini Narasimha Raghavan <roshiragavi@gmail.com>
 *  - Kristian Lyngstøl <kly@kly.no>
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

package encoder_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hamba/avro"
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/encoder"
	"github.com/telenornms/skogul/parser"
)

// TestAVROEncode tests encoding of a simple avro format from a skogul
// container
func TestAVROEncode(t *testing.T) {
	by := []byte("{\"metrics\":[{\"timestamp\":\"0001-01-01T00:00:00Z\",\"metadata\":{\"key\":\"value\"}}]}")
	testAVRO(t, by, true)
}

func testAVRO(t *testing.T, by []byte, match bool) {
	t.Helper()
	c, orig := parseAVRO(t, by)
	e := encoder.AVRO{
		Schema: "../docs/examples/avro/avro_schema",
	}
	b, err := e.Encode(c)

	if err != nil {
		t.Errorf("Encoding failed: %v", err)
		return
	}
	if len(b) <= 0 {
		t.Errorf("Encoding failed: zero length data")
		return
	}
	if !match {
		return
	}

	sorig := string(orig)
	snew := string(b)

	if len(sorig) < 2 {
		t.Logf("Encoding failed: original pre-encoded length way too short. Shouldn't happen.")
		t.FailNow()
	}

	result := strings.Compare(sorig, snew)
	if result != 0 {
		t.Errorf("Encoding failed: original and newly encoded container doesn't match")
		t.Logf("orig:\n'%s'", sorig)
		t.Logf("new:\n'%s'", snew)
		t.Logf("result\n %d", result)
		return
	}

}

func parseAVRO(t *testing.T, by []byte) (*skogul.Container, []byte) {
	t.Helper()

	data_container := skogul.Container{}
	json.Unmarshal(by, &data_container)
	schema, err := avro.Parse(`{
	    "type": "record",
	    "name": "simple",
	    "namespace": "org.hamba.avro",
	    "values": "string",
	    "fields": [
	    	{
			"name": "Metrics", "type": {
				"type": "array",
				"items": {
					"type": "record",
					"name": "metrics",
					"fields": [
						{
							"name": "Metadata",
							"type": {
								"type": "map",
								"values": ["string", "int", "long", "float", "double" ]
							}
						},
						{
							"name": "Data",
							"type": {
								"type": "map",
								"values": ["string", "int", "long", "float", "double" ]
							}
						}
					]
				}
			}
		}
	]
}
`)
	if err != nil {
		t.Logf("Failed to parse AVRO schema and test data preparation error: %v", err)
		t.FailNow()

	}

	b, err := avro.Marshal(schema, data_container)

	if err != nil {
		t.Logf("Failed to read test data file: %v", err)
		t.FailNow()
		return nil, nil
	}

	p := parser.AVRO{
		Schema: "../docs/examples/avro/avro_schema",
	}
	container, err := p.Parse(b)

	if err != nil {
		t.Logf("Failed to parse AVRO data: %v", err)
		t.FailNow()
		return nil, nil
	}

	if container == nil || container.Metrics == nil || len(container.Metrics) == 0 {
		t.Logf("Expected parsed AVRO to return a container with at least 1 metric")
		t.FailNow()
		return nil, nil
	}
	return container, b

}
