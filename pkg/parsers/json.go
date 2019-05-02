/*
 * skogul, json parser
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

package parsers

import (
	"encoding/json"
	"github.com/KristianLyng/skogul/pkg"
)

// JSON parses a byte string-representation of a Container
type JSON struct{}

// Parse accepts a byte slice of JSON data and marshals it into a container
func (x JSON) Parse(b []byte) (skogul.Container, error) {
	container := skogul.Container{}
	err := json.Unmarshal(b, &container)
	return container, err
}
