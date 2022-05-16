/*
 * skogul, metadata transformer
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
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

	"github.com/telenornms/skogul"
)

// SourceDestination provides a source and destination key, and the option
// to delete the source. At this writing, it is only used to copy from the
// data-section to metadata, but it's left intentionally generic.
type SourceDestination struct {
	Source      string `doc:"Name of the source field"`
	Destination string `doc:"The destination name/field. If left blank/undefined, the source name will be used as a destination name."`
	Delete      bool   `doc:"Set to true to delete the original. Default is to leave the original."`
}

// Metadata enforces a set of rules on metadata in all metrics, potentially
// changing the metric metadata.
type Metadata struct {
	Set             map[string]interface{} `doc:"Set metadata fields to specific values."`
	Require         []string               `doc:"Require the pressence of these fields."`
	ExtractFromData []string               `doc:"Extract a set of fields from Data and add it to Metadata. Removes the original. Obsolete, will be removed. Use CopyFromData instead."`
	CopyFromData    []SourceDestination    `doc:"Copy and potentially rename keys from the data section to the metadata section." example:"[{\"source\": \"datakey\", \"destination\": \"destkey\"},{\"source\":\"otherkey\"}]" `
	Remove          []string               `doc:"Remove these metadata fields."`
	Ban             []string               `doc:"Fail if any of these fields are present"`
}

// Transform enforces the Metadata rules
func (meta *Metadata) Transform(c *skogul.Container) error {
	for mi := range c.Metrics {
		for key, value := range meta.Set {
			if c.Metrics[mi].Metadata == nil {
				c.Metrics[mi].Metadata = make(map[string]interface{})
			}
			c.Metrics[mi].Metadata[key] = value
		}
		for _, value := range meta.Require {
			if c.Metrics[mi].Metadata == nil || c.Metrics[mi].Metadata[value] == nil {
				return skogul.Error{Source: "metadata transformer", Reason: fmt.Sprintf("missing required metadata field %s", value)}
			}
		}
		for _, extract := range meta.ExtractFromData {
			if _, ok := c.Metrics[mi].Data[extract]; !ok {
				continue
			}
			if c.Metrics[mi].Metadata == nil {
				c.Metrics[mi].Metadata = make(map[string]interface{})
			}
			c.Metrics[mi].Metadata[extract] = c.Metrics[mi].Data[extract]
			delete(c.Metrics[mi].Data, extract)
		}
		for _, cpy := range meta.CopyFromData {
			if _, ok := c.Metrics[mi].Data[cpy.Source]; !ok {
				continue
			}
			if c.Metrics[mi].Metadata == nil {
				c.Metrics[mi].Metadata = make(map[string]interface{})
			}
			// XXX: So ideally we do this only once, but
			// yeah...
			if cpy.Destination == "" {
				cpy.Destination = cpy.Source
			}
			c.Metrics[mi].Metadata[cpy.Destination] = c.Metrics[mi].Data[cpy.Source]
			if cpy.Delete {
				delete(c.Metrics[mi].Data, cpy.Source)
			}
		}
		for _, value := range meta.Remove {
			if c.Metrics[mi].Metadata == nil {
				continue
			}
			delete(c.Metrics[mi].Metadata, value)
		}
		for _, value := range meta.Ban {
			if c.Metrics[mi].Metadata == nil {
				continue
			}
			if c.Metrics[mi].Metadata[value] != nil {
				return skogul.Error{Source: "metadata transformer", Reason: fmt.Sprintf("illegal/banned metadata field %s present", value)}
			}
		}
	}
	return nil
}

func (meta *Metadata) Deprecated() error {
	if len(meta.ExtractFromData) > 0 {
		return fmt.Errorf("ExtractFromData is replaced by CopyFromData. ExtractFromData will be removed in future versions.")
	}
	return nil
}

// flattenStructure copies a nested object/array to the root level
func flattenStructure(nestedPath []string, separator string, metric *skogul.Metric) error {
	nestedObjectPath := nestedPath[0]

	// Create a nested path unless configuration says not to
	if separator != "drop" && len(nestedPath) > 1 {
		for _, p := range nestedPath[1:] {
			nestedObjectPath = fmt.Sprintf("%s%s%s", nestedObjectPath, separator, p)
		}
	}

	if separator == "drop" {
		separator = ""
		nestedObjectPath = ""
	}

	obj, err := skogul.ExtractNestedObject(metric.Data, nestedPath)

	if err == nil {
		nestedObj, ok := obj[nestedPath[len(nestedPath)-1]].(map[string]interface{})

		if !ok {

			nestedObjArray, ok := obj[nestedPath[len(nestedPath)-1]].([]interface{})
			if !ok {
				return skogul.Error{Reason: "Failed cast"}
			}

			nestedObj = make(map[string]interface{})
			for i, val := range nestedObjArray {

				obj, isMap := val.(map[string]interface{})

				// If the cast is successful, the array of items is a list of map[string]interface{},
				// and we want to extract each key to its own key in the root (prefixed with the path,
				// which may be removed by using 'drop' as separator. Array keys will still be included)
				// Otherwise, the array is a list of a primitive construct and we
				// simply prefix the key with the array index
				if isMap {
					for key, val := range obj {
						nestedObj[fmt.Sprintf("%s%s%s", fmt.Sprintf("%d", i), separator, key)] = val
					}
				} else {
					nestedObj[fmt.Sprintf("%d", i)] = val
				}
			}
		}

		for key, val := range nestedObj {
			metric.Data[fmt.Sprintf("%s%s%s", nestedObjectPath, separator, key)] = val
		}
	} else {
		return err
	}

	return nil
}

// Data enforces a set of rules on data in all metrics, potentially
// changing the metric data.
type Data struct {
	Set              map[string]interface{} `doc:"Set data fields to specific values."`
	Require          []string               `doc:"Require the pressence of these data fields."`
	Flatten          [][]string             `doc:"Flatten nested structures down to the root level"`
	FlattenSeparator string                 `doc:"Custom separator to use for flattening. Use 'drop' to drop intermediate keys. This will overwrite existing keys with the same name."`
	Remove           []string               `doc:"Remove these data fields."`
	Ban              []string               `doc:"Fail if any of these data fields are present"`
}

// Transform enforces the Metadata rules
func (data *Data) Transform(c *skogul.Container) error {
	// Set flatten separator to default value if not configured
	if data.FlattenSeparator == "" {
		data.FlattenSeparator = "__"
	}

	for mi := range c.Metrics {
		for key, value := range data.Set {
			if c.Metrics[mi].Data == nil {
				c.Metrics[mi].Data = make(map[string]interface{})
			}
			c.Metrics[mi].Data[key] = value
		}
		for _, nestedPath := range data.Flatten {
			_ = flattenStructure(nestedPath, data.FlattenSeparator, c.Metrics[mi])
		}
		for _, value := range data.Require {
			if c.Metrics[mi].Data == nil || c.Metrics[mi].Data[value] == nil {
				return skogul.Error{Source: "datadata transformer", Reason: fmt.Sprintf("missing required datadata field %s", value)}
			}
		}
		for _, value := range data.Remove {
			if c.Metrics[mi].Data == nil {
				continue
			}
			delete(c.Metrics[mi].Data, value)
		}
		for _, value := range data.Ban {
			if c.Metrics[mi].Data == nil {
				continue
			}
			if c.Metrics[mi].Data[value] != nil {
				return skogul.Error{Source: "datadata transformer", Reason: fmt.Sprintf("illegal/banned datadata field %s present", value)}
			}
		}
	}
	return nil
}
