package transformer

import (
	"sort"
	"strconv"
	"strings"

	"github.com/telenornms/skogul"
)

type Unflatten struct {

}

// Transform created a container of metrics
func (u *Unflatten) Transform(c *skogul.Container) (*skogul.Container, error) {
	metrics := c.Metrics
	newMetric := []*skogul.Metric{}

	for mi := range metrics {
		k := u.convertValues(metrics[mi])
		newMetric = append(newMetric, k)
	}

	return c, nil
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
				tmp[spl[0]] = make([]map[string]interface{}, 10)
			}

			x := tmp[spl[0]].([]map[string]interface{})
			index, _ := strconv.Atoi(spl[1])

			if index - 1 < 0 {
				index = 0
			} else {
				index = index - 1
			}

			x[index] = u.createValues(keys, index, d.Data)

			tmp[spl[0]] = x
		}
	}

	newMetric.Data = tmp

	return &newMetric
}

func (u *Unflatten) createValues(keys []string, index int,d map[string]interface{}) map[string]interface{} {
	tmp := make(map[string]interface{})
	for _, k := range keys {
		spl := strings.Split(k, ".")

		if len(spl) == 3 {
			in, _ := strconv.Atoi(spl[1])

			if in - 1 < 0 {
				in = 0
			} else {
				in = in - 1
			}

			if in == index {
				tmp[spl[2]] = d[k]
			}
		}
	}

	return tmp
}
