/*
 * skogul, core data structures
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

package skogul

import (
	"time"
)

/*
 * The Container is the top-level object that simply contains a collection
 * of metrics.
 *
 * It contains an optional template, and an array of metrics. The idea is
 * that a producer of metrics sends a bulk of metrics in a single request,
 * and we deal with it. To provide flexibility, the producer can provide an
 * (optional) template, which will be the starting point of individual
 * metrics. Example use-cases of the template include providing a timestamp
 * if all the metrics provided are from the same time stamp, and metadata
 * keys that are common, such as origin-server perhaps.
 */
type Container struct {
	Template Metric   `json:"template,omitempty"`
	Metrics  []Metric `json:"metrics"`
}

/*
 * A metric is a single set of measurements and related timestamp and
 * metadata.
 *
 * The difference between Data and Metadata is that the metadata is used to
 * identify the data, along with the timestamp. In database-terms, the
 * indexed parts are timestamp and metadata. Examples are "ifName"
 * (interface name) can be metadata, since it makes sense to search for or
 * graph data related to a single port, while "ifHCInOctets" would be data,
 * as it does NOT make sense to search for or graph data related to exactly
 * 12162 ifHCInOctets.
 */
type Metric struct {
	Time     *time.Time             `json:"timestamp,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

func (m Metric) Validate() error {
	if m.Data == nil {
		return Gerror{"Missing data for metric"}
	}
	return nil
}

func (c Container) Validate() error {
	if c.Metrics == nil {
		return Gerror{"Missing metrics[] data"}
	}
	if len(c.Metrics) <= 0 {
		return Gerror{"Empty metrics[] data"}
	}
	for i := 0; i < len(c.Metrics); i++ {
		if c.Metrics[i].Time == nil && c.Template.Time == nil {
			return Gerror{"Missing timestamp in both metric and container"}
		}
		err := c.Metrics[i].Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
