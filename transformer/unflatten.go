package transformer

import (
	"sort"
	"strings"

	"github.com/telenornms/skogul"
)

type Unflatten struct {

}

// Transform created a container of metrics
func (u *Unflatten) Transform(c *skogul.Container) error {
	metrics := c.Metrics
	newMetric := []*skogul.Metric{}

	for mi := range metrics {
		k := u.convertValues(metrics[mi])
		newMetric = append(newMetric, k)
	}

	c.Metrics = newMetric

	return nil
}

func (u *Unflatten) convertValues(d *skogul.Metric) *skogul.Metric {
	tmp := map[string]interface{}{}
	newMetric := skogul.Metric{}
	keys := make([]string, 0, len(d.Data))

	for k := range d.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		spl := strings.Split(k, ".")

		if len(spl) == 1 {
			tmp[spl[0]] = d.Data[k]
		} else if len(spl) == 2 {
			if _, ok := tmp[spl[0]]; !ok {
				tmp[spl[0]] = make(map[string]interface{})
			}
			t := tmp[spl[0]].(map[string]interface{})
			t[spl[1]] = d.Data[k]
			tmp[spl[0]] = t
		} else if len(spl) == 3 {
			if _, ok := tmp[spl[0]]; !ok {
				tmp[spl[0]] = make(map[string]map[string]interface{})
			}

			x := tmp[spl[0]].(map[string]map[string]interface{})

			if _, ok := x[spl[1]]; !ok {
				x[spl[1]] = map[string]interface{}{
					spl[2]: d.Data[k],
				}
			} else {
				x[spl[1]][spl[2]] = d.Data[k]
			}

			tmp[spl[0]] = x
		}
	}

	newMetric.Data = tmp

	return &newMetric
}
