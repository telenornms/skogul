package transformer

import (
	"sort"
	"strings"

	"github.com/telenornms/skogul"
)

type Unflatten struct {}

// Transform created a container of metrics
func (u *Unflatten) Transform(c *skogul.Container) error {
	for mi := range c.Metrics {
		c.Metrics[mi] = u.convertValues(c.Metrics[mi])
	}

	return nil
}

func (u *Unflatten) convertValues(d *skogul.Metric) *skogul.Metric {
	tmp := map[string]interface{}{}
	newMetric := &skogul.Metric{}
	keys := make([]string, 0, len(d.Data))

	for k := range d.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	// Populate keys
	for _, k := range keys {
		s := strings.Split(k, ".")
		tmp[s[0]] = map[string]interface{}{}
	}

	for _, k := range keys {
		s := strings.Split(k, ".")
		tmp[s[0]] = u.recursivelyCreateMap(tmp[s[0]].(map[string]interface{}), s[1:], d.Data[k], 0)
	}

	newMetric.Data = tmp

	return newMetric
}

func (u *Unflatten) recursivelyCreateMap(root map[string]interface{}, keys []string, value interface{}, pos int) interface{} {
	if pos == len(keys) {
		return value
	}

	if pos < len(keys) {
		if _, ok := root[keys[pos]]; ok {
			root[keys[pos]] = u.recursivelyCreateMap(root[keys[pos]].(map[string]interface{}), keys, value, pos + 1)
		} else {
			t := map[string]interface{}{
				keys[pos]: root[keys[pos]],
			}
			root[keys[pos]] = u.recursivelyCreateMap(t, keys, value, pos + 1)
		}
	}
	return root
}
