/*
 * skogul, template transformer
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
Package transformer provides the means to mutate a container as part of a
skogul.Handler, before it is passed on to a Sender.
*/
package transformer

import (
	"github.com/telenornms/skogul"
)

/*
Templater iterates over all Metrics in a Skogul Container and fills in any
missing template variable from the Template in the Container, so noone else
has to.

Skogul Templates are shallow:

        {
                "template": {
                        "metadata": {
                                "foo": "bar",
                                "geo": {
                                        "country": "Norway"
                                }
                        }
                },
                "metrics": [
                {
                        "timestamp": "...",
                        "metadata": {
                                "name": "john",
                                "geo": {
                                        "city": "Oslo"
                                }
                        },
                        "data": ...
                },
                {
                        "timestamp": "...",
                        "metadata": {
                                "name": "fred",
                                "foo": "BANANA"
                        },
                        "data": ...
                }
                ]
        }

Will result in:

        {
                "metrics": [
                {
                        "timestamp": "...",
                        "metadata": {
                                "name": "john",
                                "foo", "bar",
                                "geo": {
                                        "city": "Oslo"
                                }
                        },
                        "data": ...
                },
                {
                        "timestamp": "...",
                        "metadata": {
                                "name": "fred",
                                "foo": "BANANA",
                                "geo": {
                                        "country": "Norway"
                                }
                        },
                        "data": ...
                }
                ]
        }

The template just checked if "geo" was present in metadata or not - it did not merge missing keys.

It is good practice to use a template for any common fields, particularly timestamps.
*/
type Templater struct{}

// Transform compiles/expands the template of a container
func (t Templater) Transform(c *skogul.Container) error {
	if c.Template != nil {
		for mi, m := range c.Metrics {
			if m.Time == nil && c.Template.Time != nil {
				c.Metrics[mi].Time = c.Template.Time
			}
			for key, value := range c.Template.Metadata {
				if c.Metrics[mi].Metadata == nil {
					c.Metrics[mi].Metadata = make(map[string]interface{})
				}
				if c.Metrics[mi].Metadata[key] == nil {
					c.Metrics[mi].Metadata[key] = value
				}
			}
			for key, value := range c.Template.Data {
				if c.Metrics[mi].Data == nil {
					c.Metrics[mi].Data = make(map[string]interface{})
				}
				if c.Metrics[mi].Data[key] == nil {
					c.Metrics[mi].Data[key] = value
				}
			}
		}
		c.Template = nil
	}
	return nil
}
