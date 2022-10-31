package transformer_test

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/telenornms/skogul/parser"
	"github.com/telenornms/skogul/transformer"
	"os"
	"testing"
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

func TestTransformData(t *testing.T) {
	d := readFile("../parser/testdata/usp.bin", t)

	x := parser.USP_Parser{}
	container, err := x.Parse(d)

	if err != nil {

	}

	u := transformer.Unflatten{}
	err = u.Transform(container)

	jsonData, err := json.MarshalIndent(container.Metrics[0].Data["Report"], "", "	")

	if err != nil {
		t.Error("Could not json marshal container")
	}

	log.Println(string(jsonData))
}
