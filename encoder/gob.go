/*
 * skogul, test gob encoder
 *
 * Copyright (c) 2022 Telenor Norge AS
 * Author:
 *  - Roshini Narasimha Raghavan <roshiragavi@gmail.com>
 *  - Kristian Lyngstøl <kly@kly.no>
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

package encoder

import (
	"bytes"
	"encoding/gob"

	"github.com/telenornms/skogul"
)

type GOB struct{}

// encode the content in the skogul container as a gob format
func (x GOB) Encode(c *skogul.Container) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*c)
	b := buf.Bytes()
	return b, err
}

// encode the metrics in the skogul container as a gob format
func (x GOB) EncodeMetric(m *skogul.Metric) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*m)
	b := buf.Bytes()
	return b, err
}
