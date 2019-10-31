/*
 * skogul, edit transformer
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

package transformer

import (
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/telenornms/skogul"
)

// Replace executes a set of regular expressions, in order
type Replace struct {
	Source      string `doc:"Metadata key to read from."`
	Destination string `doc:"Metadata key to write to."`
	Regex       string `doc:"Regular expression to match."`
	Replacement string `doc:"Replacement text. Can also use $1, $2, etc to reference sub-matches."`
	regex       *regexp.Regexp
	once        sync.Once
	err         error // Used to signal persistent error - e.g., if the regular expression didn't compile.
}

// Transform executes the regular expression replacement
func (replace *Replace) Transform(c *skogul.Container) error {
	replace.once.Do(func() {
		replace.regex, replace.err = regexp.Compile(replace.Regex)
	})
	if replace.err != nil {
		return skogul.Error{Source: "replace transformer", Reason: fmt.Sprintf("Failed to compile regular expression: %s", replace.err)}
	}

	for mi := range c.Metrics {
		if c.Metrics[mi].Metadata == nil || c.Metrics[mi].Metadata[replace.Source] == nil {
			continue
		}
		str, ok := c.Metrics[mi].Metadata[replace.Source].(string)
		if !ok {
			// FIXME: What to do?
			log.Printf("Unable to transform non-string field %s with content %v", replace.Source, c.Metrics[mi].Metadata[replace.Source])
			continue
		}
		c.Metrics[mi].Metadata[replace.Destination] = string(replace.regex.ReplaceAll([]byte(str), []byte(replace.Replacement)))
	}
	return nil
}
