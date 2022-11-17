/*
 * skogul, switch transformer
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
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
	"github.com/dolmen-go/jsonptr"

	"github.com/telenornms/skogul"
)

// Case requires the path to a field ("when") and a value ("is") to match
// for the set of transformers to run
type Case struct {
	When         string                   `doc:"Used as a conditional statement on a field"`
	Is           interface{}              `doc:"Used for the specific value of the stated metadata field"`
	Transformers []*skogul.TransformerRef `doc:"The transformers to run when the defined conditional is true"`
}

// Switch is a wrapper for a list of cases
type Switch struct {
	Cases []Case `doc:"A list of switch cases "`
}

var switchLogger = skogul.Logger("transformer", "switch")

// Transform checks the cases and applies the matching transformers
func (sw *Switch) Transform(c *skogul.Container) error {
	for _, cas := range sw.Cases {

		field := cas.When
		condition := cas.Is

		for _, metric := range c.Metrics {
			var fieldValue interface{}
			// If Case.When starts with a '/', we use it as a JSON pointer.
			if cas.When[0] == '/' {
				var err error
				fieldValue, err = jsonptr.Get(metric.Metadata, cas.When)
				if err != nil {
					switchLogger.WithField("field", field).Warn("Failed to get field value from JSON pointer")
					continue
				}
			} else if metric.Metadata[field] == nil || metric.Metadata[field] == "" {
				continue
			} else {
				fieldValue = metric.Metadata[field]
			}

			if fieldValue != condition {
				continue
			}

			for _, wantedTransformerName := range cas.Transformers {
				switchLogger.WithField("wantedTransformer", wantedTransformerName).Tracef("Transformer: %v", wantedTransformerName)
				wantedTransformerName.T.Transform(c)
			}
		}
	}

	return nil
}
