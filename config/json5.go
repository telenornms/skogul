/*
 * skogul, support for JSON5 in configuration files
 *
 * Copyright (c) 2026 Telenor Norge AS
 * Author(s):
 *  - Aslak Bakkeland <aslak.bakkeland@telenor.no>
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

package config

import (
	"github.com/titanous/json5"
)

// configUnmarshal unmarshals configuration data using JSON5.
//
// JSON5 is used for all configuration files because it is a superset of JSON,
// meaning all valid JSON files parse correctly, but it also support useful
// (for config files) syntax like comments and trailing commas
//
// Note: This is only used for configuration parsing. Data parsing should not
// use this.
func configUnmarshal(data []byte, v any) error {
	return json5.Unmarshal(data, v)
}

// configSyntaxError attempts to extract position information from a parsing
// error. Returns the offset if available, or -1 if not.
func configSyntaxError(err error) (offset int, message string) {
	if err == nil {
		return -1, ""
	}

	// Check if it's a json5.SyntaxError
	if jerr, ok := err.(*json5.SyntaxError); ok {
		return int(jerr.Offset), jerr.Error()
	}

	// For other errors, return -1 to indicate no offset available
	return -1, err.Error()
}
