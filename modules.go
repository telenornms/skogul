/*
 * skogul, module automation utilities
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
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

package skogul

import (
	log "github.com/sirupsen/logrus"
)

// Module is metadata for a skogul module. It is used by the receiver,
// sender and transformer package. The Alloc() function must return a
// data structure that implements the relevant module interface, which is
// checked primarily in config.Parse.
//
// See */auto.go for how to utilize this, and config/help.go,
// cmd/skogul/main.go for how to extract information/help, and ultimately
// config/parse.go for how it is applied.
type Module struct {
	Name     string             // short name of the module (e.g: "http")
	Aliases  []string           // optional aliases (e.g. "https")
	Alloc    func() interface{} // allocation of a blank module structure
	Extras   []interface{}      // Optional additional custom data structures that should be exposed in documentation.
	Help     string             // Human-readable help description.
	AutoMake bool               // If set, this module is auto-created with default variables if referenced by implementation name without being defined in config.
}

// ModuleMap maps a name of a module to the Module data structure. Each
// type of module has its own module map. E.g.: receiver.Auto, sender.Auto
// and transformer.Auto.
type ModuleMap map[string]*Module

// Identity is used to map instances of modules to their configured name.
// E.g.: If you have 3 influx senders, the module can access
// skogul.Identity[] to distinguish between them. This is meant for logging
// and statistics. Note that this is independent of which type of module
// we're dealing with, since they all have unique addresses, while they
// don't have unique names (e.g.: You can have a receiver named "test" and
// a sender named "test" at the same time - it will still work fine).
var Identity map[interface{}]string

// Lookup will return a module if the name exists AND it should be
// autocreated. It is used during config loading to look up a module which
// is subsequently allocated.
//
// FIXME: This should probably be a replaced by Make() which returns an
// allocated module using Alloc() in the future.
func (mm ModuleMap) Lookup(name string) *Module {
	if mm[name] == nil {
		return nil
	}
	// XXX: https://github.com/telenornms/skogul/issues/182
	if name == "json" {
		log.Warn("Parser 'json' is deprecated, use 'skogul' instead.")
		return mm["skogul"]
	}
	if mm[name].AutoMake {
		return mm[name]
	}
	return nil
}

// Add adds a module to a module map, ensuring basic sanity and announcing
// it to the world, so to speak.
func (mm *ModuleMap) Add(item Module) error {
	if *mm == nil {
		*mm = make(map[string]*Module)
	}
	lm := *mm
	Assert(lm[item.Name] == nil)
	Assert(item.Alloc != nil)
	for _, alias := range item.Aliases {
		Assert(lm[alias] == nil)
		lm[alias] = &item
	}
	lm[item.Name] = &item
	return nil
}
