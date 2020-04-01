package parser_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/telenornms/skogul/parser"
)

func TestInfluxDBLineParse(t *testing.T) {
	b := []byte("system,host=testhost uptime=5464i 1585737340000000000")

	container, err := parser.InfluxDB{}.Parse(b)

	if err != nil {
		t.Errorf("Failed to parse data as influx line protocl: %v", err)
		return
	}

	if container == nil || container.Metrics == nil || len(container.Metrics) == 0 {
		t.Errorf("Expected parsed JSON to return a container with 1 metric")
		return
	}

	if container.Metrics[0].Metadata["host"] != "testhost" {
		t.Error("Expected parsed data to contain metadata field 'host'='testhost'")
	}

	uptime, castOk := container.Metrics[0].Data["uptime"].(int64)

	if !castOk {
		t.Errorf("Failed to cast value in 'uptime' data field to int64")
		return
	}

	if uptime != 5464 {
		t.Error("Expected parsed data to contain data field 'uptime'='5464'")
	}

	correctTime := time.Unix(0, 1585737340000000000)

	if err != nil {
		t.Errorf("Parsing correct time for verification failed: %s", err)
		return
	}

	if *container.Metrics[0].Time != correctTime {
		t.Errorf("Time parse failure: expected '%s' but got '%s'", correctTime, *&container.Metrics[0].Time)
	}
}

func TestInfluxDBLineParseWithoutTimestamp(t *testing.T) {
	b := []byte("system,host=testhost uptime=5464i")

	container, err := parser.InfluxDB{}.Parse(b)

	if err != nil {
		t.Errorf("Failed to parse data as influx line protocl: %v", err)
		return
	}

	if container == nil || container.Metrics == nil || len(container.Metrics) == 0 {
		t.Errorf("Expected parsed JSON to return a container with 1 metric")
		return
	}

	if container.Metrics[0].Time == nil {
		t.Errorf("Expected container to add own timestamp")
	}

	isNowish := container.Metrics[0].Time.UnixNano() - time.Now().UnixNano()

	// Arbitrary value for difference between when timestamp was created in test and the
	// one that should have been added in the parser
	if isNowish > 100 {
		t.Errorf("Expected container time to be reasonably close to timestamp generated in test, expected <=100 but got '%d'", isNowish)
	}
}

func TestInfluxDBParseFile(t *testing.T) {
	b, err := ioutil.ReadFile("./testdata/influxdb.txt")

	if err != nil {
		t.Errorf("Failed to read test data file: %v", err)
		return
	}

	container, err := parser.InfluxDB{}.Parse(b)

	if err != nil {
		t.Errorf("Failed to parse data as influx line protocl: %v", err)
		return
	}

	if container == nil || container.Metrics == nil || len(container.Metrics) == 0 {
		t.Errorf("Expected parsed JSON to return a container with at least 1 metric")
		return
	}
}