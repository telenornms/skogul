package transformer

import (
	"github.com/KristianLyng/skogul"
)

// Auto maps names to Transformers to allow auto configuration
var Auto map[string]*Transformer

// Transformer is the generic data for all transformers, used for
// auto-configuration and more.
type Transformer struct {
	Name    string
	Aliases []string
	Alloc   func() skogul.Transformer
	Help    string
}

// Add is used to announce a transformer-implementation to the world, so to
// speak. It is exported to allow out-of-package senders to exist.
func Add(r Transformer) error {
	if Auto == nil {
		Auto = make(map[string]*Transformer)
	}
	skogul.Assert(Auto[r.Name] == nil)
	skogul.Assert(r.Alloc != nil)
	for _, alias := range r.Aliases {
		skogul.Assert(Auto[alias] == nil)
		Auto[alias] = &r
	}
	Auto[r.Name] = &r
	return nil
}

func init() {
	Add(Transformer{
		Name:    "templater",
		Aliases: []string{"template", "templating"},
		Alloc:   func() skogul.Transformer { return Templater{} },
		Help:    "Executes metric templating. See separate documentationf or how skogul templating works.",
	})
	Add(Transformer{
		Name:    "metadata",
		Aliases: []string{},
		Alloc:   func() skogul.Transformer { return &Metadata{} },
		Help:    "Enforces custom-rules on metadata of metrics.",
	})
	Add(Transformer{
		Name:    "data",
		Aliases: []string{},
		Alloc:   func() skogul.Transformer { return &Data{} },
		Help:    "Enforces custom-rules for data fields of metrics.",
	})
	Add(Transformer{
		Name:    "split",
		Aliases: []string{},
		Alloc:   func() skogul.Transformer { return &Split{} },
		Help:    "Splits a metric into multiple metrics based on a field.",
	})
}
