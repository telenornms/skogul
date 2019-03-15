/*
 * gollector, common trivialities
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
	"log"
)

type Handler struct {
	Transformers []Transformer
	Senders      []Sender
}

type Sender interface {
	Send(c *GollectorContainer) error
}

type Transformer interface {
	Transform(c *GollectorContainer) error
}

type Receiver interface {
	Start() error
}

type Gerror struct {
	Reason string
}

func (e Gerror) Error() string {
	log.Printf("Error: %v", e.Reason)
	return e.Reason
}
