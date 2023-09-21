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
	tmp := []*skogul.Metric{}
	for pathKey, pathValue := range b.LookupData {
		for _, mi := range c.Metrics {
			var ptr interface{}

			ptr, _ = jsonptr.Get(mi.Data, pathKey)

			if ptr == pathValue {
				continue
			}
			tmp = append(tmp, mi)
		}
	}	

	c.Metrics = tmp

	tmp = []*skogul.Metric{}
	for pathKey, pathValue := range b.LookupMetadata {
		for _, mi := range c.Metrics {
			var ptr interface{}

			ptr, _ = jsonptr.Get(mi.Metadata, pathKey)

			if ptr == pathValue {
				continue
			}
			tmp = append(tmp, mi)
		}
	}

	c.Metrics = tmp

	return nil
}

