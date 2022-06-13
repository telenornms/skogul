package parser

import (
	"bytes"
	"encoding/gob"

	"github.com/telenornms/skogul"
)

// parses the bytes.buffer to skogul container
func Parse(b bytes.Buffer) (*skogul.Container, error) {
	container := skogul.Container{}
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&container)
	return &container, err
}

// parses the bytes.buffer to skogul metrics and wraps in a container.
func ParseMetric(b bytes.Buffer) (*skogul.Container, error) {
	container := skogul.Container{}
	metric := skogul.Metric{}
	dec := gob.NewDecoder(&b)
	err := dec.Decode(&metric)
	metrics := []*skogul.Metric{&metric}
	container.Metrics = metrics
	return &container, err
}
