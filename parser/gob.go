package parser

import (
	"bytes"
	"encoding/gob"

	"github.com/telenornms/skogul"
)

type GOB struct{}

// Parser accepts the byte buffer of GOB
func (x GOB) Parse(b []byte) (*skogul.Container, error) {
	z := bytes.NewBuffer(b)
	container := skogul.Container{}
	dec := gob.NewDecoder(z)
	err := dec.Decode(&container)
	return &container, err
}

type GOBMetric struct{}

// parses the bytes.buffer to skogul metrics and wraps in a container.
func (x GOBMetric) ParseMetric(b []byte) (*skogul.Container, error) {
	container := skogul.Container{}
	metric := skogul.Metric{}
	z := bytes.NewBuffer(b)
	dec := gob.NewDecoder(z)
	err := dec.Decode(&metric)
	metrics := []*skogul.Metric{&metric}
	container.Metrics = metrics
	return &container, err
}
