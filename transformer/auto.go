package transformer

import (
	"github.com/telenornms/skogul"
)

// Auto maps names to Transformers to allow auto configuration
var Auto skogul.ModuleMap

func init() {
	Auto.Add(skogul.Module{
		Name:    "templater",
		Aliases: []string{"template", "templating"},
		Alloc:   func() interface{} { return &Templater{} },
		Help:    "Executes metric templating. See separate documentationf or how skogul templating works.",
	})
	Auto.Add(skogul.Module{
		Name:    "metadata",
		Aliases: []string{},
		Alloc:   func() interface{} { return &Metadata{} },
		Help:    "Enforces custom-rules on metadata of metrics.",
	})
	Auto.Add(skogul.Module{
		Name:    "data",
		Aliases: []string{},
		Alloc:   func() interface{} { return &Data{} },
		Help:    "Enforces custom-rules for data fields of metrics.",
	})
	Auto.Add(skogul.Module{
		Name:    "split",
		Aliases: []string{},
		Alloc:   func() interface{} { return &Split{} },
		Help:    "Splits a metric into multiple metrics based on a field.",
	})
	Auto.Add(skogul.Module{
		Name:    "replace",
		Aliases: []string{},
		Alloc:   func() interface{} { return &Replace{} },
		Help:    "Uses a regular expression to replace the content of a metadata key, storing it to either a different metadata key, or overwriting the original.",
	})
	Auto.Add(skogul.Module{
		Name:    "switch",
		Aliases: []string{},
		Alloc:   func() interface{} { return &Switch{} },
		Help:    "Conditionally apply transformers",
		Extras:  []interface{}{Case{}},
	})
	Auto.Add(skogul.Module{
		Name:    "timestamp",
		Aliases: []string{},
		Alloc:   func() interface{} { return &Timestamp{} },
		Help:    "Extract a timestamp from the container data",
	})
}
