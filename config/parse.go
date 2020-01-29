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

	"github.com/sirupsen/logrus"
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/parser"
	"github.com/telenornms/skogul/receiver"
	"github.com/telenornms/skogul/sender"
	"github.com/telenornms/skogul/transformer"
)

var confLog = skogul.Logger("core", "config")

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
	Transformers []*skogul.TransformerRef
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
// that receiver, then unmarshals the remaining configuration onto that.
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
	var ok bool
	t.Transformer, ok = (transformer.Auto[t.Type].Alloc()).(skogul.Transformer)
	skogul.Assert(ok)
	if err := json.Unmarshal(b, &t.Transformer); err != nil {
		return skogul.Error{Source: "config parser", Reason: "Failed marshalling", Next: err}
	}
	return nil
}

// UnmarshalJSON picks up the type of the Receiver, instantiates a copy of
// that receiver, then unmarshals the remaining configuration onto that.
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
	var ok bool
	r.Receiver, ok = (receiver.Auto[r.Type].Alloc()).(skogul.Receiver)
	skogul.Assert(ok)
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
	var ok bool
	s.Sender, ok = (sender.Auto[s.Type].Alloc()).(skogul.Sender)
	skogul.Assert(ok)
	if err := json.Unmarshal(b, &s.Sender); err != nil {
		return skogul.Error{Source: "config parser", Reason: "Failed marshalling", Next: err}
	}
	return nil
}

func printSyntaxError(b []byte, offset int, text string) {
	start := offset
	start2 := 0
	lines := 0
	// Start by finding where we want to start. Work from offset and
	// decrement. start will represent "print start", start2 will
	// represent the start of the problematic line
	for i := offset; i >= 0 && lines < 3; i-- {
		start = i
		if len(b) <= i || b[i] == '\n' {
			if lines == 0 || (lines == 1 && start2 == offset) {
				start2 = start
			}
			lines++
		}
	}
	end := offset
	end2 := offset
	lines = 0
	// Next do things the other way around. End will be the actual end
	// of what to display, while end2 will be the first line after
	// "start2", e.g., the beginning of the line _after_ the
	// problematic one.
	for i := offset; i <= len(b) && lines < 3; i++ {
		end = i
		if i == len(b) || b[i] == '\n' {
			if lines == 0 {
				end2 = end
			}
			lines++
		}
	}
	fmt.Printf("Unable to parse JSON configuration at byte offset %d.\nError: %s\nContext:\n", offset, text)
	fmt.Println(string(b[start:end2]))
	for i := start2; i < (offset - 2); i++ {
		if b[i] == '	' {
			fmt.Print("--------")
		} else {
			fmt.Print("-")
		}
	}
	// We found the crappy part!
	fmt.Println("💩")
	end2++
	if end2 > len(b) {
		end2 = len(b)
	}
	fmt.Println(string(b[end2:end]))
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
	skogul.HandlerMap = skogul.HandlerMap[0:0]
	skogul.SenderMap = skogul.SenderMap[0:0]
	skogul.TransformerMap = skogul.TransformerMap[0:0]
	if err := json.Unmarshal(b, &jsonData); err != nil {
		jerr, ok := err.(*json.SyntaxError)
		if ok {
			printSyntaxError(b, int(jerr.Offset), jerr.Error())
		}
		return nil, skogul.Error{Source: "config parser", Reason: "Unable to parse JSON config", Next: err}
	}

	c := Config{}
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, skogul.Error{Source: "config parser", Reason: "Unable to parse JSON config", Next: err}
	}

	return secondPass(&c, &jsonData)
}

// File opens a config file and parses it, then returns the valid
// configuration, using Bytes()
func File(f string) (*Config, error) {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		confLog.WithError(err).Fatal("Failed to read config file")
		return nil, skogul.Error{Source: "config parser", Reason: "Failed to read config file", Next: err}
	}
	return Bytes(dat)
}

// resolveSenders iterates over the skogul.SenderMap and resolves senders,
// using the provided configuration. It zeroes the senderMap upon
// completion.
func resolveSenders(c *Config) error {
	for _, s := range skogul.SenderMap {
		confLog.WithField("sender", s.Name).Debug("Resolving sender")

		if c.Senders[s.Name] == nil {
			confLog.WithField("sender", s.Name).Error("Unresolved sender reference")
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
		logger := confLog.WithField("parser", h.Parser)

		h.Handler.Sender = h.Sender.S
		h.Handler.Transformers = make([]skogul.Transformer, 0)
		if h.Parser == "protobuf" {
			logger.Debug("Using protobuf parser")
			h.Handler.SetParser(parser.ProtoBuf{})
		} else if h.Parser == "custom-json" {
			logger.Debug("Using custom JSON parser")
			h.Handler.SetParser(parser.RawJSON{})
		} else if h.Parser == "json" || h.Parser == "" {
			logger.Debug("Using JSON parser")
			h.Handler.SetParser(parser.JSON{})
		} else {
			logger.Error("Unknown parser")
			return skogul.Error{Source: "config", Reason: fmt.Sprintf("Unknown parser %s", h.Parser)}
		}
		for _, t := range h.Transformers {
			logger = logger.WithField("transformer", t.Name)
			logger.Debug("Using predefined transformer")
			skogul.Assert(t.T != nil)
			h.Handler.Transformers = append(h.Handler.Transformers, t.T)
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

// resolveTransformers looks in the parsed config for transformers and initializes the
// actual transformers. Zeroizes the TransformerMap after if case a new
// config is applied without restarting.
func resolveTransformers(c *Config) error {
	logger := confLog.WithField("method", "resolveTransformers")
	for transformerName, t := range skogul.TransformerMap {
		logger = logger.WithField("transformer", transformerName)

		if c.Transformers[t.Name] != nil {
			logger.Debug("Using predefined transformer")
		} else {
			logger.Error("Unknown transformer")
			return skogul.Error{Source: "config", Reason: fmt.Sprintf("Unknown transformer %s", t.Name)}
		}
		skogul.Assert(c.Transformers[t.Name].Transformer != nil)
		t.T = c.Transformers[t.Name].Transformer
	}
	skogul.TransformerMap = skogul.TransformerMap[0:0]
	return nil
}

// secondPass accepts a parsed configuration as input and resolves the
// references in it, and verifies basic integrity.
func secondPass(c *Config, jsonData *map[string]interface{}) (*Config, error) {
	if err := resolveSenders(c); err != nil {
		return nil, err
	}
	if err := resolveTransformers(c); err != nil {
		return nil, err
	}
	if err := resolveHandlers(c); err != nil {
		return nil, err
	}

	for idx, h := range c.Handlers {
		confLog.WithField("handler", idx).Debug("Verifying handler configuration")
		if err := verifyItem("handler", idx, h.Handler); err != nil {
			return nil, err
		}

		configStruct := reflect.TypeOf(h.Handler)
		verifyOnlyRequiredConfigProps(jsonData, "handlers", idx, configStruct)
	}
	for idx, t := range c.Transformers {
		confLog.WithField("transformer", idx).Debug("Verifying transformer configuration")
		if err := verifyItem("transformer", idx, t.Transformer); err != nil {
			return nil, err
		}

		configStruct := reflect.ValueOf(t.Transformer).Elem().Type()
		verifyOnlyRequiredConfigProps(jsonData, "transformers", idx, configStruct)
	}
	for idx, s := range c.Senders {
		confLog.WithField("sender", idx).Debug("Verifying sender configuration")
		if err := verifyItem("sender", idx, s.Sender); err != nil {
			return nil, err
		}

		configStruct := reflect.ValueOf(s.Sender).Elem().Type()
		verifyOnlyRequiredConfigProps(jsonData, "senders", idx, configStruct)
	}
	for idx, r := range c.Receivers {
		confLog.WithField("receiver", idx).Debug("Verifying receiver configuration")
		if err := verifyItem("receiver", idx, r.Receiver); err != nil {
			return nil, err
		}

		configStruct := reflect.ValueOf(r.Receiver).Elem().Type()
		verifyOnlyRequiredConfigProps(jsonData, "receivers", idx, configStruct)
	}

	return c, nil
}

// verifyItem checks if the item implements Verifier and if so, verifies
// the item. Otherwise, returns nil.
func verifyItem(family string, name string, item interface{}) error {
	i, ok := item.(skogul.Verifier)
	if !ok {
		confLog.WithFields(logrus.Fields{"family": family, "name": name}).Trace("No verifier found")
		return nil
	}
	err := i.Verify()
	if err != nil {
		confLog.WithFields(logrus.Fields{"family": family, "name": name}).Error("Invalid item configuration")
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("%s %s isn't valid", family, name), Next: err}
	}
	confLog.WithFields(logrus.Fields{"family": family, "name": name}).Trace("Verified OK")
	return nil
}

func findFieldsOfStruct(T reflect.Type) []string {
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
	configFamily, ok := (*rawConfig)[family].(map[string]interface{})
	if !ok {
		confLog.WithFields(logrus.Fields{
			"family":  family,
			"section": section,
		}).Warnf("Failed to cast config family to map[string]interface{}")
		return nil
	}

	configSection, ok := configFamily[section].(map[string]interface{})
	if !ok {
		confLog.WithFields(logrus.Fields{
			"family":  family,
			"section": section,
		}).Warnf("Failed to cast config section to map[string]interface{}")
		return nil
	}
	return configSection
}

func verifyOnlyRequiredConfigProps(rawConfig *map[string]interface{}, family, handler string, T reflect.Type) []string {
	requiredProps := findFieldsOfStruct(T)

	relevantConfig := getRelevantRawConfigSection(rawConfig, family, handler)

	superfluousProperties := make([]string, 0)

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
			superfluousProperties = append(superfluousProperties, prop)
			confLog.WithField("property", prop).Warn("Configuration property configured but not defined in code (this property won't change anything, is it wrongly defined?)")
		}
	}

	return superfluousProperties
}
