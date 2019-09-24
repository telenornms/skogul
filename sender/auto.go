/*
 * skogul, sender automation
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

package sender

import (
	"log"

	"github.com/KristianLyng/skogul"
)

// Sender provides a framework that all sender-implementations should
// follow, and allows auto-initialization.
type Sender struct {
	Name  string
	Alloc func() skogul.Sender
	Help  string
}

// Auto maps sender-names to sender implementation, used for auto
// configuration.
var Auto map[string]*Sender

// Add announces the existence of a sender to the world at large.
func Add(s Sender) error {
	if Auto == nil {
		Auto = make(map[string]*Sender)
	}
	if Auto[s.Name] != nil {
		log.Panicf("BUG: Attempting to overwrite existing auto-add sender %v", s.Name)
	}
	if s.Alloc == nil {
		log.Printf("No alloc function for %s", s.Name)
	}
	Auto[s.Name] = &s
	return nil
}
