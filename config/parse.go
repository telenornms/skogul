/*
 * skogul, configuration parsing
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
Package config handles Skogul configuration parsing.
*/
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"github.com/KristianLyng/skogul/transformer"
)

// Sender wraps the skogul.Sender for configuration parsing.
type Sender struct {
	Type   string
	Sender skogul.Sender `json:"-"`
}

// Receiver wraps the skogul.Receiver for configuration parsing.
type Receiver struct {
	Type     string
	Receiver skogul.Receiver `json:"-"`
}

// Handler wraps skogul.Handler for configuration parsing.
type Handler struct {
	Parser       string
	Transformers []string
	Sender       skogul.SenderRef
	Handler      skogul.Handler `json:"-"`
}

// Config encapsulates all configuration for Skogul, and represent the
// top-level configuration object.
type Config struct {
	Handlers  map[string]*Handler
	Receivers map[string]*Receiver
	Senders   map[string]*Sender
}

// MarshalJSON for a receiver marshals the actual instantiated receiver,
// then merges it to add "type". Probably not the most efficient
// implementation, since it does marshal-unmarshal-merge-marshal, but since
// this isn't really performance sensitive, that's ok.
func (r *Receiver) MarshalJSON() ([]byte, error) {
	nest, err := json.Marshal(r.Receiver)
	if err != nil {
		return nil, err
	}
	var merged map[string]interface{}
	if err := json.Unmarshal(nest, &merged); err != nil {
		return nil, err
	}
	merged["type"] = r.Type
	return json.Marshal(merged)
}

// UnmarshalJSON picks up the type of the Receiver, instantiates a copy of
// that receiver, than unmarshals the remaining configuration onto that.
func (r *Receiver) UnmarshalJSON(b []byte) error {
	type tType struct {
		Type string
	}
	var t tType
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	r.Type = t.Type
	if receiver.Auto[r.Type] == nil {
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Unknown receiver %v", r.Type)}
	}
	if receiver.Auto[r.Type].Alloc == nil {
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Bad receiver %v", r.Type)}
	}
	r.Receiver = receiver.Auto[r.Type].Alloc()
	if err := json.Unmarshal(b, &r.Receiver); err != nil {
		return skogul.Error{Source: "config parser", Reason: "Failed marshalling", Next: err}
	}
	return nil
}

// MarshalJSON marshals Sender config. See MarshalJSON for receiver - same
// same.
func (s *Sender) MarshalJSON() ([]byte, error) {
	nest, err := json.Marshal(s.Sender)
	if err != nil {
		return nil, err
	}
	var merged map[string]interface{}
	if err := json.Unmarshal(nest, &merged); err != nil {
		return nil, err
	}
	merged["type"] = s.Type
	return json.Marshal(merged)
}

// UnmarshalJSON for Sender. See UnmarshalJSON for Receiver - same same.
func (s *Sender) UnmarshalJSON(b []byte) error {
	type tType struct {
		Type string
	}
	var t tType
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	s.Type = t.Type
	if sender.Auto[s.Type] == nil {
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Unknown sender %v", s.Type)}
	}
	if sender.Auto[s.Type].Alloc == nil {
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Bad sender %v", s.Type)}
	}
	s.Sender = sender.Auto[s.Type].Alloc()
	if err := json.Unmarshal(b, &s.Sender); err != nil {
		return skogul.Error{Source: "config parser", Reason: "Failed marshalling", Next: err}
	}
	return nil
}

// File opens a config file and parses it, then returns the valid
// configuration.
func File(f string) (*Config, error) {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, skogul.Error{Source: "config parser", Reason: "Failed to read config file", Next: err}
	}
	var c *Config
	c, err = Bytes(dat)
	return c, err
}

// Bytes parses json in the provided byte array and returns a
// configuration. To avoid dependency loops, the config parsing is,
// unfortunately, dependent on global variables in the skogul top-level
// package - this means Bytes() will change state of skogul.SenderMap and
// skogul.HandlerMap
func Bytes(b []byte) (*Config, error) {
	c := Config{}
	err := json.Unmarshal(b, &c)
	if err != nil {
		return nil, skogul.Error{Source: "config parser", Reason: "Unable to parse JSON config", Next: err}
	}

	for _, s := range skogul.SenderMap {
		if c.Senders[s.Name] == nil {
			return nil, skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Unresolved sender reference %s", s.Name)}
		}
		s.S = c.Senders[s.Name].Sender
	}
	skogul.SenderMap = skogul.SenderMap[0:0]
	for _, h := range c.Handlers {
		h.Handler.Sender = h.Sender.S
		h.Handler.Transformers = make([]skogul.Transformer, 0)
		if h.Parser == "json" {
			h.Handler.Parser = parser.JSON{}
		} else {
			return nil, skogul.Error{Source: "config", Reason: fmt.Sprintf("Unknown parser %s", h.Parser)}
		}
		for _, t := range h.Transformers {
			if t == "templater" {
				h.Handler.Transformers = append(h.Handler.Transformers, transformer.Templater{})
			} else {
				return nil, skogul.Error{Source: "config", Reason: fmt.Sprintf("Unknown transformer %s", t)}
			}
		}
	}
	for _, h := range skogul.HandlerMap {
		if c.Handlers[h.Name] == nil {
			return nil, skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Unresolved handler reference %s", h.Name)}
		}
		h.H = &(c.Handlers[h.Name].Handler)
	}
	skogul.HandlerMap = skogul.HandlerMap[0:0]
	for _, h := range c.Handlers {
		e := h.Handler.Verify()
		if e != nil {
			return nil, skogul.Error{Source: "config", Reason: "Handler corrupt", Next: e}
		}
	}
	for idx, s := range c.Senders {
		sv, ok := s.Sender.(skogul.Verifier)
		if ok {
			err := sv.Verify()
			if err != nil {
				return nil, skogul.Error{Source: "config parser", Reason: fmt.Sprintf("sender %s isn't valid", idx), Next: err}
			}
		}
	}
	for idx, r := range c.Receivers {
		rv, ok := r.Receiver.(skogul.Verifier)
		if ok {
			err := rv.Verify()
			if err != nil {
				return nil, skogul.Error{Source: "config parser", Reason: fmt.Sprintf("receiver %s isn't valid", idx), Next: err}
			}
		}
	}
	return &c, nil
}