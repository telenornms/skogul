package encoder

import (
	"os"

	"github.com/hamba/avro"
	"github.com/telenornms/skogul"
)

type AVRO struct {
	Schema avro.Schema
	In     skogul.Container
}

func (x AVRO) Encode(c *skogul.Container) ([]byte, error) {
	var A AVRO
	b, _ := os.ReadFile("./schema/avro_schema")
	A.Schema = avro.MustParse(string(b))
	A.In = *c
	return avro.Marshal(A.Schema, A.In)
}
