package parser_test

import (
	"io/ioutil"
	"testing"

	"github.com/telenornms/skogul/parser"
)

func TestGOBParser(t *testing.T) {
	parseGOB(t, "./testdata/testdata.gob")
}

func parseGOB(t *testing.T, file string) {
	t.Helper()

	b, err := ioutil.ReadFile(file)

	if err != nil {
		t.Logf("Failed to read test data file: %v", err)
		t.FailNow()

	}
	container, err := parser.GOB{}.Parse(b)

	if err != nil {
		t.Logf("Failed to parse GOB data: %v", err)
		t.FailNow()

	}

	if container == nil || container.Metrics == nil || len(container.Metrics) == 0 {
		t.Logf("Expected parsed GOB to return a container with at least 1 metric")
		t.FailNow()

	}

}
