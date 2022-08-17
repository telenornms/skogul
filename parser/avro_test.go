/*
 * skogul, avro parser
 *
 * Copyright (c) 2022 Telenor Norge AS
 * Author(s):
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
package parser_test

import (
	"io/ioutil"
	"testing"

	"github.com/telenornms/skogul/parser"
)

func TestAVROParser(t *testing.T) {
	parseAVRO(t, "./testdata/avro_testdata.json")
}

func parseAVRO(t *testing.T, file string) {
	t.Helper()

	b, err := ioutil.ReadFile(file)

	if err != nil {
		t.Logf("Failed to read test data file: %v", err)
		t.FailNow()

	}

	p := parser.AVRO{
		Schema: "../docs/examples/avro/avro_schema",
	}
	container, err := p.Parse(b)

	if err != nil {
		t.Logf("Failed to parse AVRO data: %v", err)
		t.FailNow()

	}

	if container == nil || container.Metrics == nil || len(container.Metrics) == 0 {
		t.Logf("Expected parsed AVRO to return a container with at least 1 metric")
		t.FailNow()

	}

}
