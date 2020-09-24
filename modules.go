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
	if mm[name].AutoMake {
		return mm[name]
	}
	return nil
}

// Add adds a module to a module map, ensuring basic sanity and announcing
// it to the world, so to speak.
func (amap *ModuleMap) Add(item Module) error {
	if *amap == nil {
		*amap = make(map[string]*Module)
	}
	lm := *amap
	Assert(lm[item.Name] == nil)
	Assert(item.Alloc != nil)
	for _, alias := range item.Aliases {
		Assert(lm[alias] == nil)
		lm[alias] = &item
	}
	lm[item.Name] = &item
	return nil
}
