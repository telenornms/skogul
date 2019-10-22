/*
 * skogul, split transformer
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.com>
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

	"github.com/telenornms/skogul"
)

// Split is the configuration for the split transformer
type Split struct {
	Field []string `doc:"Split into multiple metrics based on this field (each field denotes the path to a nested object element)."`
	Fail  bool     `doc:"Fail the transformer entirely if split is unsuccsessful on a metric container. This will prevent successive transformers from working."`
}

// Transform splits the thing
func (split *Split) Transform(c *skogul.Container) error {

	metrics := c.Metrics

	if split.Field != nil {
		splitMetrics, err := split.splitMetricsByObjectKey(&metrics)
		if err == nil {
			c.Metrics = splitMetrics
		} else if split.Fail {
			return fmt.Errorf("failed to split metrics: %v", err)
		}
	}

	return nil
}

// splitMetricsByObjectKey splits the metrics into multiple metrics based on a key in a list of sub-metrics
func (split *Split) splitMetricsByObjectKey(metrics *[]*skogul.Metric) ([]*skogul.Metric, error) {
	origMetrics := *metrics
	var newMetrics []*skogul.Metric

	for mi := range origMetrics {
		splitObj, err := skogul.ExtractNestedObject(origMetrics[mi].Data, split.Field)

		if err != nil {
			if !split.Fail {
				continue
			}
			return nil, fmt.Errorf("Failed to extract nested obj '%v' from '%v' to string/interface map", split.Field, origMetrics[mi].Data)
		}

		metrics, ok := splitObj[split.Field[len(split.Field)-1]].([]interface{})

		if !ok {
			return nil, fmt.Errorf("Failed to cast '%v' to string/interface map on '%s'", origMetrics[mi].Data, split.Field[0])
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
