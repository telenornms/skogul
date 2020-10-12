/*
 * skogul, json parser
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
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

package parser

import (
	"encoding/json"

	"github.com/telenornms/skogul"
)

// JSON parses a byte string-representation of a Container
type JSON struct{}

// Parse accepts a byte slice of JSON data and marshals it into a container
func (x JSON) Parse(b []byte) (*skogul.Container, error) {
	container := skogul.Container{}
	err := json.Unmarshal(b, &container)
	return &container, err
}

// RawJSON can be used when the JSON format does not conform to the final JSON format of skogul,
// e.g. when it is used as the first step of parsing from a third party source where modifying
// the source data structure might be hard/impossible
type RawJSON struct{}

// Parse accepts a byte slice of JSON data and marshals it into an empty skogul.Container
func (data RawJSON) Parse(b []byte) (*skogul.Container, error) {

	// The Validate() func of a container expects a timestamp to be valid.
	// Better way to fix?
	time := skogul.Now()
	metric := skogul.Metric{
		Metadata: make(map[string]interface{}),
		Time:     &time,
	}

	err := json.Unmarshal(b, &metric.Data)

	if err != nil {
		return nil, err
	}

	container := skogul.Container{
		Metrics: []*skogul.Metric{&metric},
	}

	return &container, nil
}
