/*
 * skogul, mnr parser
 *
 * Copyright (c) 2020 Telenor Norge AS
 * Author(s):
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.com>
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

package transformer

import (
	"fmt"
	"sync"

	"github.com/telenornms/skogul"
)

var parseLog = skogul.Logger("transformer", "parser")

// Parse contains the configuration for a skogul transformer which can be used to
// transform parts of a metric by applying a parser to it. Examples include
// JSON data with nested JSON stringified fields or JSON data with nested
// syslog/RFC5424/key-value data.
type Parse struct {
	Parser              skogul.ParserRef `doc:"Name of skogul parser to use. Can be auto-initialisable, or can be defined in 'parsers' of the config."`
	Source              string           `doc:"Name of data field to apply parser to"`
	Destination         string           `doc:"Field to create with the parsed data (default: [Source]_data)"`
	DestinationMetadata string           `doc:"Field to create with the parsed metadata (default: [Source]_metadata)"`
	Keep                bool             `doc:"Keep the unparsed value (default: false)"`
	Append              bool             `doc:"Insert the values directly in the existing container, instead of creating new fields (destination and destination_metadata). (default: false)"`
	once                sync.Once
}

// init sets the default values of Destination and DestinationMetadata
// if they are not set by the user.
func (p *Parse) init() {
	if p.Destination == "" {
		p.Destination = fmt.Sprintf("%s_data", p.Source)
	}
	if p.DestinationMetadata == "" {
		p.DestinationMetadata = fmt.Sprintf("%s_metadata", p.Source)
	}
}

// Transform transforms a container by applying the Parse configuration
// to the container metrics.
func (p *Parse) Transform(c *skogul.Container) error {
	p.once.Do(func() {
		p.init()
	})

	for i, metric := range c.Metrics {
		val, ok := metric.Data[p.Source].(string)
		if !ok {
			parseLog.Error("Failed to cast value to string")
			continue
		}
		b := []byte(val)
		parsed, err := p.Parser.P.Parse(b)
		if err != nil {
			return err
		}
		if len(parsed.Metrics) == 0 {
			return skogul.Error{Reason: "Parsed 0 metrics", Source: "transformer-parser"}
		}

		if !p.Keep {
			delete(c.Metrics[i].Data, p.Source)
		}

		if p.Append {
			// Insert directly into existing container
			for field, val := range parsed.Metrics[0].Metadata {
				c.Metrics[i].Metadata[field] = val
			}
			for field, val := range parsed.Metrics[0].Data {
				c.Metrics[i].Data[field] = val
			}
		} else {
			// Add new fields
			data := make(map[string]interface{})
			metadata := make(map[string]interface{})

			for field, val := range parsed.Metrics[0].Metadata {
				metadata[field] = val
			}
			for field, val := range parsed.Metrics[0].Data {
				data[field] = val
			}
			c.Metrics[i].Data[p.DestinationMetadata] = metadata
			c.Metrics[i].Data[p.Destination] = data
		}
	}
	return nil
}

// Verify verifies that the configuration is valid, e.g. disallowing
// configurations which does not make sense, and warns about
// configuration options which override each other (if set).
func (p *Parse) Verify() error {
	if p.Source == "" {
		return skogul.Error{Reason: "Missing source field", Source: "transformer-parser"}
	}
	if (p.Destination != "" || p.DestinationMetadata != "") && p.Append {
		parseLog.Warn("Destination or DestinationMetadata configured at the same time as Append - Append takes precedence, and Destination(Metadata) will be ignored.")
	}
	return nil
}
