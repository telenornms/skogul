/*
 * skogul, splunk writer
 *
 * Copyright (c) 2020 Telenor Norge AS
 * Author(s):
 *  - Håkon Solbjørg <Hakon.Solbjorg@telenor.com>
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
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/telenornms/skogul"
)

var splunkLog = skogul.Logger("sender", "splunk")

// Splunk contains the configuration parameters for this sender.
type Splunk struct {
	URL           string        `doc:"URL to Splunk HTTP Event Collector (HEC)"`
	Token         skogul.Secret `doc:"Token for HTTP Authorization header for HEC endpoint."`
	Index         string        `doc:"Custom Splunk index to send event to."`
	HostnameField string        `doc:"Name of the metadata field with the hostname. Note, this might have to be transformed into metadata depending on the input data."`
	SourceField   string        `doc:"Name of the metadata field with the source. Will fallback to the value set in Source if not found."`
	Source        string        `doc:"Set the source of the data. If used in conjunction with SourceField, SourceField will take precedence and this will be the fallback."`
	HTTP          *HTTP         `doc:"HTTP sender options. URL is overwritten from this config, the rest will be HTTP sender defaults unless overridden."`
	ok            bool
	once          sync.Once
}

// splunkEvent describes the structure of a Splunk
// HTTP Event Collector event
type splunkEvent struct {
	Time   *time.Time             `json:"time,omitempty"`
	Host   string                 `json:"host,omitempty"`
	Index  string                 `json:"index,omitempty"`
	Source string                 `json:"source,omitempty"`
	Event  map[string]interface{} `json:"event"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

// prepare converts a skogul container into the appropriate
// format expceted by the Splunk HEC collector as defined here
// https://docs.splunk.com/Documentation/Splunk/8.0.6/Data/FormateventsforHTTPEventCollector
func (s *Splunk) prepare(c *skogul.Container) ([]splunkEvent, error) {
	events := make([]splunkEvent, len(c.Metrics))
	for i, metric := range c.Metrics {
		t := metric.Time
		if metric.Time == nil {
			t = c.Template.Time
		}

		host := ""
		source := ""
		if s.HostnameField != "" && metric.Metadata != nil && metric.Metadata[s.HostnameField] != nil {
			host = fmt.Sprintf("%v", metric.Metadata[s.HostnameField])
			delete(metric.Metadata, s.HostnameField)
		}
		if s.SourceField != "" && metric.Metadata != nil && metric.Metadata[s.SourceField] != nil {
			source = fmt.Sprintf("%v", metric.Metadata[s.SourceField])
			if source == "" {
				source = s.Source
			}
			delete(metric.Metadata, s.SourceField)
		}
		events[i] = splunkEvent{
			Time:   t,
			Event:  metric.Data,
			Index:  s.Index,
			Host:   host,
			Source: source,
			Fields: metric.Metadata,
		}
	}
	return events, nil
}

// MarshalJSON overrides the marshalling of the
// 'Time' field on a splunkEvent struct to provide
// the 'seconds.milliseconds' value which HEC expects.
func (e *splunkEvent) MarshalJSON() ([]byte, error) {
	t := 0.0
	if e.Time != nil {
		// Convert the time to the format HEC expects,
		// which is 'sec.ms'. We can achieve this by
		// using UnixNano() and dividing
		// back up to seconds, which gives us
		// the milliseconds as decimals.
		t = float64(e.Time.UnixNano()) / 1e9
	}

	// Type aliasing splunkEvent to change the
	// marshalling of 'Time' but keeping the
	// default marshaller for the rest.
	type splunkEventOutput splunkEvent
	return json.Marshal(&struct {
		Time float64 `json:"time,omitempty"`
		*splunkEventOutput
	}{
		Time:              t,
		splunkEventOutput: (*splunkEventOutput)(e),
	})
}

// init handles initializing the Splunk sender.
// since Splunk uses the HTTP sender under the hood,
// we initialize that one too.
func (s *Splunk) init() {
	s.HTTP.URL = s.URL
	s.HTTP.init()
	s.HTTP.Headers["authorization"] = fmt.Sprintf("Splunk %s", s.Token.Expose())
	s.ok = s.HTTP.ok
}

// Send sends a skogul container to Splunk HEC
func (s *Splunk) Send(c *skogul.Container) error {
	s.once.Do(func() {
		s.init()
	})
	if !s.ok {
		return fmt.Errorf("splunk sender not in OK state")
	}

	events, err := s.prepare(c)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for _, event := range events {
		b, err := json.Marshal(&event)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON data to Splunk: %w", err)
		}
		buffer.Write(b)
	}
	if err := s.HTTP.sendBytes(buffer.Bytes()); err != nil {
		return fmt.Errorf("sendBytes failed: %w", err)
	}

	return nil
}

// Verify verifies that the sender config is valid
func (s *Splunk) Verify() error {
	if s.URL == "" {
		return skogul.MissingArgument("URL")
	}
	if s.Token.Expose() == "" {
		return skogul.MissingArgument("Token")
	}
	if s.Index == "" {
		splunkLog.Info("No Splunk index configured, Splunk will send events to its default index.")
	}
	if s.HostnameField == "" {
		splunkLog.Warning("No HostnameField specified, Splunk events will not be metadata-tagged with hostnames")
	}
	if err := s.HTTP.Verify(); err != nil {
		// Verify HTTP handler, but if it contains an error about
		// missing URL, disregard it, since we will override that
		// during our own init().
		if !strings.Contains(err.Error(), "missing required configuration option `URL'") {
			return fmt.Errorf("failed to verify HTTP sender for Splunk: %w", err)
		}
		if s.HTTP.URL != "" && s.URL != "" && s.HTTP.URL != s.URL {
			return fmt.Errorf("duplicate conflicting URLs specified: URL defined in both HTTP.URL and URL - pick one")
		}
	}
	return nil
}
