/*
 * gollector, core data structures
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

package common

import (
	"time"
)

type GollectorContainer struct {
	Template GollectorMetric   `json:"template,omitempty"`
	Metrics  []GollectorMetric `json:"metrics"`
}

type GollectorMetric struct {
	Time     *time.Time             `json:"timestamp,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

func (m GollectorMetric) Validate() error {
	if m.Data == nil {
		return Gerror{"Missing data for metric"}
	}
	return nil
}

func (c GollectorContainer) Validate() error {
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
