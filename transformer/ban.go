package transformer

import (
	"log"

	"github.com/dolmen-go/jsonptr"
	"github.com/telenornms/skogul"
)

type Ban struct {
	Lookup map[string]interface{} `doc:"Map of key value pairs to lookup in metrics. Looks in data and metadata fields. Key is json pointer, value any. E.g. /foo/bar: \"bar\""`
}

func (b *Ban) Transform(c *skogul.Container) error {
	for pathKey, pathValue := range b.Lookup {
		for metricKey, mi := range c.Metrics {
			var ptr interface{}

			ptr, _ = jsonptr.Get(mi.Data, pathKey)

			if ptr == pathValue {
				c.Metrics[metricKey] = &skogul.Metric{}
			}
		}
	}

	for pathKey, pathValue := range b.Lookup {
		for metricKey, mi := range c.Metrics {
			var ptr interface{}

			ptr, _ = jsonptr.Get(mi.Metadata, pathKey)

			if ptr == pathValue {
				log.Printf("%v %v", ptr, pathValue)
				c.Metrics[metricKey] = &skogul.Metric{}
			}
		}
	}

	return nil
}
