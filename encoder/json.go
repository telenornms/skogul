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
