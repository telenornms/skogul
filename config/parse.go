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

	log "github.com/sirupsen/logrus"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/parser"
	"github.com/telenornms/skogul/receiver"
	"github.com/telenornms/skogul/sender"
	"github.com/telenornms/skogul/transformer"
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

// Transformer wraps skogul.Transformer
type Transformer struct {
	Type        string
	Transformer skogul.Transformer `json:"-"`
}

// Config encapsulates all configuration for Skogul, and represent the
// top-level configuration object.
type Config struct {
	Handlers     map[string]*Handler
	Receivers    map[string]*Receiver
	Senders      map[string]*Sender
	Transformers map[string]*Transformer
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
func (t *Transformer) UnmarshalJSON(b []byte) error {
	type tType struct {
		Type string
	}
	var myt tType
	if err := json.Unmarshal(b, &myt); err != nil {
		return err
	}
	t.Type = myt.Type
	if transformer.Auto[t.Type] == nil {
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Unknown transformer %v", t.Type)}
	}
	if transformer.Auto[t.Type].Alloc == nil {
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Bad transformer %v", t.Type)}
	}
	t.Transformer = transformer.Auto[t.Type].Alloc()
	if err := json.Unmarshal(b, &t.Transformer); err != nil {
		return skogul.Error{Source: "config parser", Reason: "Failed marshalling", Next: err}
	}
	return nil
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

// Bytes parses json in the provided byte array and returns a
// configuration.
//
// It does this by first doing a pass where it just does JSON
// unmarshalling, which also updates sender and handler reference tables
// globally (unfortunately...), then calling secondPass(), which resolves
// references and does a final validation.
func Bytes(b []byte) (*Config, error) {
	c := Config{}
	if err := json.Unmarshal(b, &c); err != nil {
		log.WithError(err).Fatal("The JSON configuration is improperly formatted JSON")
		return nil, skogul.Error{Source: "config parser", Reason: "Unable to parse JSON config", Next: err}
	}

	return secondPass(&c)
}

// File opens a config file and parses it, then returns the valid
// configuration, using Bytes()
func File(f string) (*Config, error) {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		log.WithError(err).Fatal("Failed to read config file")
		return nil, skogul.Error{Source: "config parser", Reason: "Failed to read config file", Next: err}
	}
	return Bytes(dat)
}

// resolveSenders iterates over the skogul.SenderMap and resolves senders,
// using the provided configuration. It zeroes the senderMap upon
// completion.
func resolveSenders(c *Config) error {
	for _, s := range skogul.SenderMap {
		if c.Senders[s.Name] == nil {
			return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Unresolved sender reference %s", s.Name)}
		}
		s.S = c.Senders[s.Name].Sender
	}
	skogul.SenderMap = skogul.SenderMap[0:0]
	return nil
}

// resolveHandlers iterates over handlers and instantiates them, since
// there is no unmarshaller (or need for one) that does this. It then
// iterates of the skogul.HandlerMap and resolves the handler references to
// the actual handlers.
//
// It then zeroes the skogul.HandlerMap
func resolveHandlers(c *Config) error {
	for _, h := range c.Handlers {
		h.Handler.Sender = h.Sender.S
		h.Handler.Transformers = make([]skogul.Transformer, 0)
		if h.Parser == "protobuf" {
			h.Handler.Parser = parser.ProtoBuf{}
		} else if h.Parser == "json" || h.Parser == "" {
			h.Handler.Parser = parser.JSON{}
		} else {
			return skogul.Error{Source: "config", Reason: fmt.Sprintf("Unknown parser %s", h.Parser)}
		}
		for _, t := range h.Transformers {
			var nextT skogul.Transformer
			if c.Transformers[t] != nil {
				nextT = c.Transformers[t].Transformer
			} else if t == "templater" {
				nextT = transformer.Templater{}
			} else {
				return skogul.Error{Source: "config", Reason: fmt.Sprintf("Unknown transformer %s", t)}
			}
			h.Handler.Transformers = append(h.Handler.Transformers, nextT)
		}
	}
	for _, h := range skogul.HandlerMap {
		if c.Handlers[h.Name] == nil {
			return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Unresolved handler reference %s", h.Name)}
		}
		h.H = &(c.Handlers[h.Name].Handler)
	}
	skogul.HandlerMap = skogul.HandlerMap[0:0]
	return nil
}

// secondPass accepts a parsed configuration as input and resolves the
// references in it, and verifies basic integrity.
func secondPass(c *Config) (*Config, error) {
	if err := resolveSenders(c); err != nil {
		return nil, err
	}
	if err := resolveHandlers(c); err != nil {
		return nil, err
	}
	for idx, h := range c.Handlers {
		if err := verifyItem("handler", idx, h.Handler); err != nil {
			return nil, err
		}
	}
	for idx, s := range c.Senders {
		if err := verifyItem("sender", idx, s.Sender); err != nil {
			return nil, err
		}
	}
	for idx, r := range c.Receivers {
		if err := verifyItem("receiver", idx, r.Receiver); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// verifyItem checks if the item implements Verifier and if so, verifies
// the item. Otherwise, returns nil.
func verifyItem(family string, name string, item interface{}) error {
	i, ok := item.(skogul.Verifier)
	if !ok {
		return nil
	}
	err := i.Verify()
	if err != nil {
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("%s %s isn't valid", family, name), Next: err}
	}
	return nil
}
