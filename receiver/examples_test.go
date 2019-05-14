package receiver_test

import (
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"github.com/KristianLyng/skogul/transformer"
)

/*
HTTP can have different skogul.Handler's for different paths, with potentially different behaviors.
*/
func ExampleHTTP() {
	h := receiver.HTTP{Address: "localhost:8080"}
	template := skogul.Handler{Parser: parser.JSON{}, Transformers: []skogul.Transformer{transformer.Templater{}}, Sender: sender.Debug{}}
	noTemplate := skogul.Handler{Parser: parser.JSON{}, Sender: sender.Debug{}}
	h.Handle("/template", &template)
	h.Handle("/notemplate", &noTemplate)
	h.Start()
}

/*
Using New() sets up a single handler on the specified path. This is the same as
*/
func ExampleHTTP_new() {
	handler := skogul.Handler{Parser: parser.JSON{}, Transformers: []skogul.Transformer{transformer.Templater{}}, Sender: sender.Debug{}}
	h, err := receiver.New("http://localhost:8080/foobar", handler)
	if err != nil {
		panic(err)
	}
	h.Start()
}
