package sender

import (
	"testing"

	"github.com/telenornms/skogul"
)

func TestPrepare(t *testing.T) {
	now := skogul.Now()

	metadata := make(map[string]interface{})
	metadata["source"] = "skogul"

	data := make(map[string]interface{})
	data["log"] = "some log"

	metric := skogul.Metric{
		Time:     &now,
		Metadata: metadata,
		Data:     data,
	}

	c := skogul.Container{
		Metrics: []*skogul.Metric{&metric},
	}

	splunk := Splunk{}

	_, err := splunk.prepare(&c)
	if err != nil {
		t.Errorf("failed to prepare splunk event: %s", err)
	}
}

func TestPrepareSplunkMetadataAsFields(t *testing.T) {
	now := skogul.Now()

	metadata := make(map[string]interface{})
	metadata["foo"] = "bar"

	data := make(map[string]interface{})
	data["log"] = "some log"

	metric := skogul.Metric{
		Time:     &now,
		Metadata: metadata,
		Data:     data,
	}

	c := skogul.Container{
		Metrics: []*skogul.Metric{&metric},
	}

	splunk := Splunk{}

	splunkEvents, err := splunk.prepare(&c)
	if err != nil {
		t.Errorf("failed to prepare splunk event: %s", err)
	}

	if splunkEvents[0].Fields["foo"] != "bar" {
		t.Error("expected to find a 'field' named 'foo' with value 'bar' in splunkEvent")
	}
}
