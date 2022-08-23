package encoder

import (
	"bytes"
	"encoding/gob"

	"github.com/telenornms/skogul"
)

type GOB struct{}

// encode the content in the skogul container as a gob format
func (x GOB) Encode(c *skogul.Container) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*c)
	b := buf.Bytes()
	return b, err
}

// encode the metrics in the skogul container as a gob format
func (x GOB) EncodeMetric(m *skogul.Metric) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*m)
	b := buf.Bytes()
	return b, err
}
