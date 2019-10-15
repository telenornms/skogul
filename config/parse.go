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
	"reflect"
	"strings"

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
	var jsonData map[string]interface{}

	if err := json.Unmarshal(b, &jsonData); err != nil {
		log.WithError(err).Fatal("The JSON configuration is improperly formatted JSON")
	}

	c := Config{}
	if err := json.Unmarshal(b, &c); err != nil {
		log.WithError(err).Fatal("The JSON configuration is improperly formatted JSON")
		return nil, skogul.Error{Source: "config parser", Reason: "Unable to parse JSON config", Next: err}
	}

	return secondPass(&c, &jsonData)
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
		log.WithField("sender", s.Name).Debug("Resolving sender")

		if c.Senders[s.Name] == nil {
			log.WithField("sender", s.Name).Error("Unresolved sender reference")
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
		logger := log.WithField("parser", h.Parser)

		h.Handler.Sender = h.Sender.S
		h.Handler.Transformers = make([]skogul.Transformer, 0)
		if h.Parser == "protobuf" {
			logger.Debug("Using protobuf parser")
			h.Handler.Parser = parser.ProtoBuf{}
		} else if h.Parser == "json" || h.Parser == "" {
			logger.Debug("Using JSON parser")
			h.Handler.Parser = parser.JSON{}
		} else {
			logger.Error("Unknown parser")
			return skogul.Error{Source: "config", Reason: fmt.Sprintf("Unknown parser %s", h.Parser)}
		}
		for _, t := range h.Transformers {
			logger = logger.WithField("transformer", t)

			var nextT skogul.Transformer
			if c.Transformers[t] != nil {
				logger.Debug("Using predefined transformer")
				nextT = c.Transformers[t].Transformer
			} else if t == "templater" {
				logger.Debug("Using templating transformer")
				nextT = transformer.Templater{}
			} else {
				logger.Error("Unknown transformer")
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
func secondPass(c *Config, jsonData *map[string]interface{}) (*Config, error) {
	if err := resolveSenders(c); err != nil {
		return nil, err
	}
	if err := resolveHandlers(c); err != nil {
		return nil, err
	}
	for idx, h := range c.Handlers {
		log.WithField("handler", idx).Debug("Verifying handler configuration")
		if err := verifyItem("handler", idx, h.Handler); err != nil {
			return nil, err
		}
	}
	for idx, s := range c.Senders {
		log.WithField("sender", idx).Debug("Verifying sender configuration")
		if err := verifyItem("sender", idx, s.Sender); err != nil {
			return nil, err
		}
	}
	for idx, r := range c.Receivers {
		log.WithField("receiver", idx).Debug("Verifying receiver configuration")
		if err := verifyItem("receiver", idx, r.Receiver); err != nil {
			return nil, err
		}

		receiverConfigStruct := reflect.ValueOf(r.Receiver).Elem().Type()
		_ = findMissingRequiredConfigProps(jsonData, "receivers", idx, receiverConfigStruct)
		verifyOnlyRequiredConfigProps(jsonData, "receivers", idx, receiverConfigStruct)
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
		log.WithFields(log.Fields{"family": family, "name": name}).Error("Invalid item configuration")
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("%s %s isn't valid", family, name), Next: err}
	}
	return nil
}

func findFieldsOfStruct(T reflect.Type) []string {
	log.WithField("type", T.Name()).Debug("Finding defined fields for type")
	requiredProps := make([]string, 0)
	switch T.Kind() {
	case reflect.Struct:
		for i := 0; i < T.NumField(); i++ {
			field := T.Field(i)
			jsonTag := field.Tag.Get("json")

			property := field.Name
			if jsonTag != "" {
				property = jsonTag
			}
			requiredProps = append(requiredProps, property)
		}
	}

	return requiredProps
}

func getRelevantRawConfigSection(rawConfig *map[string]interface{}, family, section string) map[string]interface{} {
	// log.Debugf("Fetching config section '%s' for '%s' from %v", section, family, rawConfig)
	return (*rawConfig)[family].(map[string]interface{})[strings.ToLower(section)].(map[string]interface{})
}

func findMissingRequiredConfigProps(rawConfig *map[string]interface{}, family, handler string, T reflect.Type) []string {
	requiredProps := findFieldsOfStruct(T)
	log.Debugf("Required fields: %v", requiredProps)

	relevantConfig := getRelevantRawConfigSection(rawConfig, family, handler)

	missingProps := make([]string, 0)

	for _, requiredProp := range requiredProps {
		if relevantConfig[strings.ToLower(requiredProp)] == nil {
			log.WithField("property", strings.ToLower(requiredProp)).Error("Missing required configuration property")
			missingProps = append(missingProps, requiredProp)
		}
	}

	return missingProps
}

func verifyOnlyRequiredConfigProps(rawConfig *map[string]interface{}, family, handler string, T reflect.Type) {
	requiredProps := findFieldsOfStruct(T)
	log.Debugf("Required fields: %v", requiredProps)

	relevantConfig := getRelevantRawConfigSection(rawConfig, family, handler)

	for prop := range relevantConfig {
		propertyDefined := false

		if prop == "type" {
			// Skip the type specifying what type this is
			continue
		}

		for _, requiredProp := range requiredProps {
			if strings.ToLower(prop) == strings.ToLower(requiredProp) {
				propertyDefined = true
				break
			}
		}
		if !propertyDefined {
			log.WithField("property", prop).Warn("Property configured but not defined in code (this property won't change anything, is it wrongly defined?)")
		}
	}
}
