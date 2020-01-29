/*
 * skogul, test json parser
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

package parser

import "testing"

import "io/ioutil"

// TestJSONParse tests parsing of a simple JSON document to skogul
// container
func TestJSONParse(t *testing.T) {
	b := []byte("{\"metrics\":[{\"timestamp\":\"2019-03-15T11:08:02+01:00\",\"metadata\":{\"key\":\"value\"},\"data\":{\"string\":\"text\",\"float\":1.11,\"integer\":5}}]}")
	x := JSON{}
	_, err := x.Parse(b)
	if err != nil {
		t.Errorf("JSON.Parse(b) failed: %s", err)
	}
}

func BenchmarkJSONParse(b *testing.B) {
	by := []byte("{\"metrics\":[{\"timestamp\":\"2019-03-15T11:08:02+01:00\",\"metadata\":{\"key\":\"value\"},\"data\":{\"string\":\"text\",\"float\":1.11,\"integer\":5}}]}")
	x := JSON{}
	for i := 0; i < b.N; i++ {
		x.Parse(by)
	}
}

func TestRawJSONParse(t *testing.T) {
	b, err := ioutil.ReadFile("./testdata/raw.json")

	if err != nil {
		t.Errorf("Failed to read test data file: %v", err)
		return
	}

	container, err := RawJSON{}.Parse(b)

	if err != nil {
		t.Errorf("Failed to parse JSON data: %v", err)
		return
	}

	if container == nil || container.Metrics == nil || len(container.Metrics) == 0 {
		t.Errorf("Expected parsed JSON to return a container with at least 1 metric")
		return
	}
}

func BenchmarkRawJSONParse(b *testing.B) {
	by := []byte(`{"string":"text","float":1.11,"integer":5,"timestamp":"2019-03-15T11:08:02+01:00","key":"value"}`)
	x := RawJSON{}
	for i := 0; i < b.N; i++ {
		x.Parse(by)
	}
}
