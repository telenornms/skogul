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

package transformers

import (
	"github.com/KristianLyng/skogul/pkg"
	"log"
)

/*
The template transformer is responsbile for consuming the template of a
container and copying it to the individual metrics.

Note that templates are SHALLOW. If you have a template that goes

"metadata": { "foo": { "key1": "value1" }

then a metric provides

"metadata": { "foo": { "key2": "value2" } }

then the "foo.key1" will NOT be part of the end result, as the template logic
only sees that the metric has a "foo" object, and thus does not try to do
further templating of that key.

Implementation-wise, this isn't really perfect - ideally it should be done
on read, but because this provides a degree of flexibility that is hard to
achieve with a getter, this seems like a reasonable compromise.
*/
type Templater struct{}

func (t Templater) Transform(c *skogul.Container) error {
	for mi, m := range c.Metrics {
		if m.Time == nil && c.Template.Time != nil {
			c.Metrics[mi].Time = c.Template.Time
		}
		for key, value := range c.Template.Metadata {
			if c.Metrics[mi].Metadata[key] == nil {
				c.Metrics[mi].Metadata[key] = value
			}
		}
		for key, value := range c.Template.Data {
			if c.Metrics[mi].Data[key] == nil {
				c.Metrics[mi].Data[key] = value
			}
		}
	}
	return nil
}
