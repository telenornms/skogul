package parser_test

import (
	"os"
	"testing"

	"github.com/telenornms/skogul/parser"
)

func readFile(file string, t *testing.T) []byte {
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

func TestParseFile(t *testing.T) {
	d := readFile("testdata/usp.bin", t)

	x := parser.ProtoBuffer{}
	x.Parse(d)
}
