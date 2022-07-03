package parser

import (
	"os"

	"github.com/hamba/avro"
	"github.com/telenornms/skogul"
)

type AVRO struct {
	Schema avro.Schema
	In     skogul.Container
}

func (x AVRO) Parse(b []byte) (*skogul.Container, error) {
	var A AVRO
	s, _ := os.ReadFile("./schema/avro_schema")
	A.Schema = avro.MustParse(string(s))
	err := avro.Unmarshal(A.Schema, b, &A.In)
	return &A.In, err

}
