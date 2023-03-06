package transformer

import (
	"github.com/dolmen-go/jsonptr"
	"github.com/telenornms/skogul"
)

type Ban struct {
	LookupData     map[string]interface{} `doc:"Map of key value pairs to lookup in metrics. Looks in data fields. Key is json pointer, value any. E.g. /foo/bar: \"bar\""`
	LookupMetadata map[string]interface{} `doc:"Map of key value pairs to lookup in metrics. Looks in metadata fields. Key is json pointer, value any. E.g. /foo/bar: \"bar\""`
}

func (b *Ban) Transform(c *skogul.Container) error {
	for pathKey, pathValue := range b.LookupData {
		for metricKey, mi := range c.Metrics {
			var ptr interface{}

			ptr, _ = jsonptr.Get(mi.Data, pathKey)

			if ptr == pathValue {
				c.Metrics = append(c.Metrics[:metricKey], c.Metrics[metricKey+1:]...)
			}
		}
	}

	for pathKey, pathValue := range b.LookupMetadata {
		for metricKey, mi := range c.Metrics {
			var ptr interface{}

			ptr, _ = jsonptr.Get(mi.Metadata, pathKey)

			if ptr == pathValue {
				c.Metrics = append(c.Metrics[:metricKey], c.Metrics[metricKey+1:]...)
			}
		}
	}

	return nil
}
