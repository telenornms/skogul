package encoder

import (
	"encoding/json"
	"github.com/telenornms/skogul"
)

type JSON struct{}

func (x JSON) Encode(c *skogul.Container) ([]byte, error) {
	b, err := json.Marshal(*c)
	return b, err
}
func (x JSON) EncodeMetric(m *skogul.Metric) ([]byte, error) {
	b, err := json.Marshal(*m)
	return b, err
}
