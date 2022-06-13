package encoder

import (
	"bytes"
	"encoding/gob"

	"github.com/telenornms/skogul"
)

type GOB struct{}

// encode the content in the skogul container as a gob format
func (x GOB) Encode(c *skogul.Container) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	return enc.Encode(*c)
}

// encode the metrics in the skogul container as a gob format
func (x GOB) EncodeMetric(m *skogul.Metric) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	return enc.Encode(*m)
}
