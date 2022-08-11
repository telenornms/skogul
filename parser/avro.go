/*
 * skogul, test avro parser
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

package parser

import (
	"fmt"
	"os"
	"sync"
	"github.com/hamba/avro"
	"github.com/telenornms/skogul"
)

type AVRO struct {
	Schema string
	once sync.Once
}


func (x AVRO) Parse(b []byte) (*skogul.Container, error) {
	x.once.Do(func() {
	s, error:= os.ReadFile(x.Schema)
	if error != nil {
		err = fmt.Errorf("Schema read error")
	} else {
	schema := avro.MustParse(string(s))
	}})
	container := skogul.Container{}
	return avro.Unmarshal(schema, b, &container)

func (x AVRO) ParseMetric(m *skogul.Metric) ([]byte, error) {
	return nil, fmt.Errorf("Not supported")
}
