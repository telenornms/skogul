/*
 * skogul, test protobuf parser
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

import (
	"os"
	"testing"
)

type failer interface {
	Fatalf(format string, args ...interface{})
	Helper()
}

func readProtobufFile(t failer, file string) []byte {
	t.Helper()
	b := make([]byte, 9000)
	f, err := os.Open(file)
	if err != nil {
		t.Fatalf("unable to open protobuf packet file: %v", err)
	}
	defer f.Close()
	n, err := f.Read(b)
	if err != nil {
		t.Fatalf("unable to read protobuf packet file: %v", err)
	}
	if n == 0 {
		t.Fatalf("read 0 bytes from protobuf packet file....")
	}
	return b[0:n]
}

func TestProtoBuf(t *testing.T) {
	b := readProtobufFile(t, "testdata/protobuf-packet.bin")
	x := ProtoBuf{}
	c, err := x.Parse(b)
	if err != nil {
		t.Errorf("ProtoBuf.Parse(b) failed: %s", err)
	}
	if c == nil {
		t.Errorf("ProtoBuf.Parse(b) returned nil-container")
	}
}

func BenchmarkProtoBufParse(b *testing.B) {
	by := readProtobufFile(b, "testdata/protobuf-packet.bin")
	x := ProtoBuf{}
	for i := 0; i < b.N; i++ {
		x.Parse(by)
	}
}
