/*
 * skogul, test avro parser
 *
 * Copyright (c) 2022 Telenor Norge AS
 * Author:
 *  - Roshini Narasimha Raghavan <roshiragavi@gmail.com>
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
	"os"

	"github.com/hamba/avro"
	"github.com/telenornms/skogul"
)

type AVRO struct {
	Schema avro.Schema
	In     skogul.Container
}

func (x AVRO) Parse(b []byte) (*skogul.Container, error) {
	var A AVRO
	s, _ := os.ReadFile("./schema/avro_schema")
	A.Schema = avro.MustParse(string(s))
	err := avro.Unmarshal(A.Schema, b, &A.In)
	return &A.In, err

}
