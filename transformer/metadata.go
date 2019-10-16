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

	"github.com/KristianLyng/skogul"
)

// Metadata enforces a set of rules on metadata in all metrics, potentially
// changing the metric metadata.
type Metadata struct {
	Set             map[string]interface{} `doc:"Set metadata fields to specific values."`
	Split           []string               `doc:"Split into multiple metrics based on this field (each field denotes the path to a nested object element)."`
	Require         []string               `doc:"Require the pressence of these fields."`
	ExtractFromData []string               `doc:"Extract a set of fields from Data and add it to Metadata."`
	Remove          []string               `doc:"Remove these metadata fields."`
	Ban             []string               `doc:"Fail if any of these fields are present"`
}

// Transform enforces the Metadata rules
func (meta *Metadata) Transform(c *skogul.Container) error {
	metrics := c.Metrics

	if meta.Split != nil {
		splitMetrics, err := splitMetricsByObjectKey(&metrics, meta)
		if err == nil {
			c.Metrics = splitMetrics
		} else {
			fmt.Println("failed to split metrics")

			// dont hard fail metrics unless we really want to
			if false {
				return fmt.Errorf("failed to split metrics: %v", err)
			}
		}
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
		for _, extract := range meta.ExtractFromData {

			c.Metrics[mi].Metadata[extract] = c.Metrics[mi].Data[extract]
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

// extractNestedObject extracts an object from a nested object structure. All intermediate objects has to map[string]interface{}
func extractNestedObject(object map[string]interface{}, keys []string) (map[string]interface{}, error) {
	if len(keys) == 1 {
		return object, nil
	}

	next, ok := object[keys[0]].(map[string]interface{})

	if !ok {
		return nil, skogul.Error{Reason: "Failed to cast nested object to map[string]interface{}"}
	}

	return extractNestedObject(next, keys[1:])
}

// splitMetricsByObjectKey splits the metrics into multiple metrics based on a key in a list of sub-metrics
func splitMetricsByObjectKey(metrics *[]*skogul.Metric, metadata *Metadata) ([]*skogul.Metric, error) {
	origMetrics := *metrics
	var newMetrics []*skogul.Metric

	for mi := range origMetrics {
		splitObj, err := extractNestedObject(origMetrics[mi].Data, metadata.Split)

		if err != nil {
			return nil, fmt.Errorf("Failed to extract nested obj '%v' from '%v' to string/interface map", metadata.Split, origMetrics[mi].Data)
		}

		metrics, ok := splitObj[metadata.Split[len(metadata.Split)-1]].([]interface{})

		if !ok {
			return nil, fmt.Errorf("Failed to cast '%v' to string/interface map on '%s'", origMetrics[mi].Data, metadata.Split[0])
		}

		for _, obj := range metrics {
			// Create a new metrics object as a copy of the original one, then reassign the data field
			metricsData, ok := obj.(map[string]interface{})

			if !ok {
				return nil, fmt.Errorf("Failed to cast '%v' to string/interface map", obj)
			}

			newMetric := *origMetrics[mi]
			newMetric.Data = metricsData
			newMetric.Metadata = make(map[string]interface{})

			for key, val := range origMetrics[mi].Metadata {
				newMetric.Metadata[key] = val
			}

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
