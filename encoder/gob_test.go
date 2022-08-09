/*
 * skogul, test gob encoder
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
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/encoder"
	"github.com/telenornms/skogul/parser"
)

// TestGOBEncode tests encoding of a simple GOB format from a skogul
// container
func TestGOBEncode(t *testing.T) {
	testGOB(t, "./testdata/testdata.gob", true)
}

func testGOB(t *testing.T, file string, match bool) {
	t.Helper()
	c, orig := parseGOB(t, file)
	b, err := encoder.GOB{}.Encode(c)

	if err != nil {
		t.Errorf("Encoding %s failed: %v", file, err)
		return
	}
	if len(b) <= 0 {
		t.Errorf("Encoding %s failed: zero length data", file)
		return
	}
	if !match {
		return
	}

	//The test gob data is generated using a separate program and there are differences in 3 locations compared to the encoded data.
	//These differences are fixed by replacing the differences with exact matches before the comparison.
	sorig := string(orig)
	var origstring1 string
	var old1 string
	var old2 string
	var old3 string
	var new1 string
	var new2 string
	var new3 string
	var result1, result2, result3 []byte
	var n int
	origstring1 = sorig
	old1 = "\x1d"
	old2 = "\x0e"
	old3 = "main"
	n = 1
	new1 = "\x1f"
	new2 = "\x10"
	new3 = "skogul"
	result1 = bytes.Replace([]byte(origstring1), []byte(old1), []byte(new1), n)
	result2 = bytes.Replace([]byte(result1), []byte(old2), []byte(new2), n)
	result3 = bytes.Replace([]byte(result2), []byte(old3), []byte(new3), n)

	sorig_trim := string(result3)
	snew := string(b)

	if len(sorig) < 2 {
		t.Logf("Encoding %s failed: original pre-encoded length way too short. Shouldn't happen.", file)
		t.FailNow()
	}

	result := strings.Compare(sorig_trim, snew)
	if result != 0 {
		t.Errorf("Encoding %s failed: original and newly encoded container doesn't match", file)
		t.Logf("orig:\n'%s'", sorig_trim)
		t.Logf("new:\n'%s'", snew)
		t.Logf("result\n %d", result)
		return
	}

}

func parseGOB(t *testing.T, file string) (*skogul.Container, []byte) {
	t.Helper()

	b, err := ioutil.ReadFile(file)

	if err != nil {
		t.Logf("Failed to read test data file: %v", err)
		t.FailNow()
		return nil, nil
	}
	container, err := parser.GOB{}.Parse(b)

	if err != nil {
		t.Logf("Failed to parse GOB data: %v", err)
		t.FailNow()
		return nil, nil
	}

	if container == nil || container.Metrics == nil || len(container.Metrics) == 0 {
		t.Logf("Expected parsed GOB to return a container with at least 1 metric")
		t.FailNow()
		return nil, nil
	}
	return container, b

}
