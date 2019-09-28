package transformer

import (
	"github.com/KristianLyng/skogul"
	"log"
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
	if Auto[r.Name] != nil {
		log.Panicf("BUG: Attempting to overwrite existing auto-add transformer %v", r.Name)
	}
	if r.Alloc == nil {
		log.Panicf("No alloc function for %s", r.Name)
	}
	for _, alias := range r.Aliases {
		if Auto[alias] != nil {
			log.Panicf("BUG: An alias(%s) for transformer %s overlaps an existing transformer %s", alias, r.Name, Auto[alias].Name)
		}
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
}
