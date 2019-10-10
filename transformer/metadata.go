/*
 * skogul, metadata transformer
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
 * 02110-1301  USA
 */

package transformer

import (
	"fmt"

	"github.com/getlantern/deepcopy"

	"github.com/KristianLyng/skogul"
)

// Metadata enforces a set of rules on metadata in all metrics, potentially
// changing the metric metadata.
type Metadata struct {
	Set     map[string]interface{} `doc:"Set metadata fields to specific values."`
	Split   []string               `doc:"Split into multiple metrics based on this field (dot '.' denotes nested object element)."`
	Require []string               `doc:"Require the pressence of these fields."`
	Remove  []string               `doc:"Remove these metadata fields."`
	Ban     []string               `doc:"Fail if any of these fields are present"`
}

// Transform enforces the Metadata rules
func (meta *Metadata) Transform(c *skogul.Container) error {
	metrics := c.Metrics

	if meta.Split != nil {
		splitMetrics, err := splitMetricsByObjectKey(&metrics, meta)
		if err != nil {
			return fmt.Errorf("failed to split metrics: %v", err)
		}

		c.Metrics = splitMetrics
	}

	for mi := range c.Metrics {
		for key, value := range meta.Set {
			if c.Metrics[mi].Metadata == nil {
				c.Metrics[mi].Metadata = make(map[string]interface{})
			}
			c.Metrics[mi].Metadata[key] = value
		}
		for _, value := range meta.Require {
			if c.Metrics[mi].Metadata == nil || c.Metrics[mi].Metadata[value] == nil {
				return skogul.Error{Source: "metadata transformer", Reason: fmt.Sprintf("missing required metadata field %s", value)}
			}
		}
		for _, value := range meta.Remove {
			if c.Metrics[mi].Metadata == nil {
				continue
			}
			delete(c.Metrics[mi].Metadata, value)
		}
		for _, value := range meta.Ban {
			if c.Metrics[mi].Metadata == nil {
				continue
			}
			if c.Metrics[mi].Metadata[value] != nil {
				return skogul.Error{Source: "metadata transformer", Reason: fmt.Sprintf("illegal/banned metadata field %s present", value)}
			}
		}
	}
	return nil
}

// splitMetricsByObjectKey splits the metrics into multiple metrics based on a key in a list of sub-metrics
func splitMetricsByObjectKey(metrics *[]*skogul.Metric, metadata *Metadata) ([]*skogul.Metric, error) {
	origMetrics := *metrics
	var newMetrics []*skogul.Metric

	// @ToDo: Parse this to a path if more than one element deep (i.e. my.nested.json.structure should unfurl)
	splitPath := metadata.Split[0]

	for mi := range origMetrics {
		metricObj := origMetrics[mi].Data[splitPath]

		for _, obj := range metricObj.([]interface{}) {
			// Create a new metrics object as a copy of the original one, then reassign the data field
			newMetric := skogul.Metric{}
			deepcopy.Copy(&newMetric, &origMetrics[mi])
			newMetric.Data = obj.(map[string]interface{})

			newMetrics = append(newMetrics, &newMetric)
		}
	}

	return newMetrics, nil
}

// Data enforces a set of rules on data in all metrics, potentially
// changing the metric data.
type Data struct {
	Set     map[string]interface{} `doc:"Set data fields to specific values."`
	Require []string               `doc:"Require the pressence of these data fields."`
	Remove  []string               `doc:"Remove these data fields."`
	Ban     []string               `doc:"Fail if any of these data fields are present"`
}

// Transform enforces the Metadata rules
func (data *Data) Transform(c *skogul.Container) error {
	for mi := range c.Metrics {
		for key, value := range data.Set {
			if c.Metrics[mi].Data == nil {
				c.Metrics[mi].Data = make(map[string]interface{})
			}
			c.Metrics[mi].Data[key] = value
		}
		for _, value := range data.Require {
			if c.Metrics[mi].Data == nil || c.Metrics[mi].Data[value] == nil {
				return skogul.Error{Source: "datadata transformer", Reason: fmt.Sprintf("missing required datadata field %s", value)}
			}
		}
		for _, value := range data.Remove {
			if c.Metrics[mi].Data == nil {
				continue
			}
			delete(c.Metrics[mi].Data, value)
		}
		for _, value := range data.Ban {
			if c.Metrics[mi].Data == nil {
				continue
			}
			if c.Metrics[mi].Data[value] != nil {
				return skogul.Error{Source: "datadata transformer", Reason: fmt.Sprintf("illegal/banned datadata field %s present", value)}
			}
		}
	}
	return nil
}
