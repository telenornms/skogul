/*
 * skogul, test avro encoder
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
	"fmt"
	"os"
	"sync"

	"github.com/hamba/avro"
	"github.com/telenornms/skogul"
)

type AVRO struct {
	Schema string
	s      avro.Schema
	err    error
	once   sync.Once
}

func (x AVRO) Encode(c *skogul.Container) ([]byte, error) {
	x.once.Do(func() {
		b, err := os.ReadFile(x.Schema)
		x.err = err
		if x.err != nil {
			fmt.Errorf("Schema read error")
		} else {
			x.s = avro.MustParse(string(b))
		}
	})
	if x.err != nil {
		return nil, x.err
	}
	return avro.Marshal(x.s, c)
}

func (x AVRO) EncodeMetric(m *skogul.Metric) ([]byte, error) {
	return nil, fmt.Errorf("not supported")
}
