package receiver_test

import (
	"github.com/telenornms/skogul"
	"time"
)

var validContainer = skogul.Container{}

func init() {

	now := time.Now()
	m := skogul.Metric{}
	m.Time = &now
	m.Metadata = make(map[string]interface{})
	m.Data = make(map[string]interface{})
	m.Metadata["foo"] = "bar"
	m.Data["tall"] = 5
	validContainer.Metrics = []*skogul.Metric{&m}
}
