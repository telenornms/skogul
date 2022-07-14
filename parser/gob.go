/*
 * skogul, test gob parser
 *
 * Copyright (c) 2022 Telenor Norge AS
 * Author(s):
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
package parser

import (
	"bytes"
	"encoding/gob"

	"github.com/telenornms/skogul"
)

type GOB struct{}

// Parser accepts the byte buffer of GOB
func (x GOB) Parse(b []byte) (*skogul.Container, error) {
	z := bytes.NewBuffer(b)
	container := skogul.Container{}
	dec := gob.NewDecoder(z)
	err := dec.Decode(&container)
	return &container, err
}

type GOBMetric struct{}

// parses the bytes.buffer to skogul metrics and wraps in a container.
func (x GOBMetric) ParseMetric(b []byte) (*skogul.Container, error) {
	container := skogul.Container{}
	metric := skogul.Metric{}
	z := bytes.NewBuffer(b)
	dec := gob.NewDecoder(z)
	err := dec.Decode(&metric)
	metrics := []*skogul.Metric{&metric}
	container.Metrics = metrics
	return &container, err
}
