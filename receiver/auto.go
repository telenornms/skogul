/*
 * skogul, receiver automation
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

/*
Package receiver provides various skogul Receivers that accept data and
execute a handler. They are the "inbound" API of Skogul.
*/
package receiver

import (
	"log"

	"github.com/KristianLyng/skogul"
)

// Auto maps names to Receivers to allow auto configuration
var Auto map[string]*Receiver

// Receiver is the generic data for all receivers, used for
// auto-configuration and more.
type Receiver struct {
	Name  string
	Alloc func() skogul.Receiver
	Help  string
}

func Add(r Receiver) error {
	if Auto == nil {
		Auto = make(map[string]*Receiver)
	}
	if Auto[r.Name] != nil {
		log.Panicf("BUG: Attempting to overwrite existing auto-add receiver %v", r.Name)
	}
	if r.Alloc == nil {
		log.Panicf("No alloc function for %s", r.Name)
	}
	Auto[r.Name] = &r
	return nil
}
