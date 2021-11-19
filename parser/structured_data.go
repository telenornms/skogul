/*
 * skogul, rfc 5424 structured data parser
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

package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/telenornms/skogul"
)

var sdLog = skogul.Logger("parser", "structured_data")

// StructuredData supports parsing RFC5424 structured data through the Parse() function
// Note: This does not parse a full syslog message.
type StructuredData struct{}

// Parse converts RFC5424 Structured Data data into a skogul Container
func (sd *StructuredData) Parse(bytes []byte) (*skogul.Container, error) {
	metrics, err := sd.parseStructuredData(bytes)
	if err != nil {
		return nil, err
	}

	return &skogul.Container{
		Metrics: metrics,
	}, nil
}

// mnrParseData takes the raw input and parses it
// this takes care of splitting input on newlines etc
func (sd *StructuredData) parseStructuredData(data []byte) ([]*skogul.Metric, error) {
	lines := bytes.Split(data, []byte("\n"))

	metrics := make([]*skogul.Metric, 0)

	timestamp := skogul.Now()

	for _, l := range lines {
		line := bytes.TrimSpace(l)
		if len(line) == 0 {
			// Skip empty lines
			continue
		}

		kvScanner := bufio.NewScanner(bytes.NewReader(line))
		kvScanner.Split(splitKeyValuePairs)

		has_hostname := false

		var metric *skogul.Metric
		for {
			canContinue := kvScanner.Scan()

			tag := strings.Trim(kvScanner.Text(), "\u0000")
			tagValue := strings.SplitN(tag, "=", 2)

			if len(tagValue) == 1 {
				if strings.TrimSpace(tagValue[0]) == "" {
					return nil, skogul.Error{Reason: "Got invalid data in the middle of a structured data line", Source: "structured_data-parser"}
				}
				if metric != nil {
					if len(metric.Data) > 0 {
						metrics = append(metrics, metric)
					} else {
						sdLog.Tracef("NOT Creating new metric because existing is empty, metric: %v, line:'%s'", metric, string(line))
						break
					}
				}
				metric = &skogul.Metric{
					Time:     &timestamp,
					Metadata: make(map[string]interface{}),
					Data:     make(map[string]interface{}),
				}
				metric.Metadata["sd-id"] = tagValue[0]
				has_hostname = true
				continue
			} else if !has_hostname {
				metric = &skogul.Metric{
					Time:     &timestamp,
					Metadata: make(map[string]interface{}),
					Data:     make(map[string]interface{}),
				}
				has_hostname = true
			} else if len(tagValue) != 2 {
				break
			}

			paramName := tagValue[0]
			paramValue := tagValue[1][1 : len(tagValue[1])-1] // remove leading and trailing "s

			// @ToDo: Support multiple paramName with different paramValue
			// if the value already exists, replace it with an array ?

			metric.Data[paramName] = paramValue

			if !canContinue {
				break
			}
		}
		if metric != nil {
			// Note: We add metrics even if they have no data fields.
			// Intended from sender or misconfigured sender?
			metrics = append(metrics, metric)
		}
	}
	if len(metrics) == 0 {
		sdLog.WithField("lines", len(lines)).Warnf("RFC5424/Structured Data parser failed to parse any of the %d lines", len(lines))
		return nil, skogul.Error{Reason: "Failed to parse RFC5424 lines", Source: "structured_data-parser"}
	}
	return metrics, nil
}

// splitKeyValuePairs splits a section (tag key=value pairs or field key=value pairs)
func splitKeyValuePairs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	fieldWidth, newData := structuredDataParser(data, true)

	returnChars := len(newData)

	if atEOF || returnChars == len(data) {
		// EOF, return with what we have left
		return returnChars + 1, newData[:returnChars], nil
	}

	// Skip the trailing comma between each key=value pair, but still advance counter
	return fieldWidth, newData[:returnChars], nil
}

// struturedDataParser parses a structured data-line.
// A boolean flag decides whether or not escape characters should remain in the output
// or have their prepending escape character removed.
func structuredDataParser(data []byte, removeEscapedCharsFromResult bool) (int, []byte) {
	openQuote := false
	escape := false
	escapeChars := make([]int, 0)
	escapeCharsWidth := make([]int, 0)
	start := 0
	for width := 0; start < len(data); start += width {
		var c rune
		c, width = utf8.DecodeRune(data[start:])

		if escape {
			escape = false
			continue
		}

		// If we receive an un-escaped ] character, this section is done
		// and we'll restart parsing of the rest (if any) as a new section.
		if c == ']' {
			break
		}

		// If there is an open quote, continue until we find the closing quote
		if openQuote {
			if c == '"' {
				// We found the closing quote, mark it and continue regular operations
				openQuote = false
				continue
			}
			// Fast forward loop until we find the closing quote
			continue
		}

		// We found the opening of a quote, continue until we find the closing one
		if c == '"' {
			openQuote = true
			continue
		}

		// Skip next char
		if c == '\\' {
			escape = true
			if removeEscapedCharsFromResult {
				escapeChars = append([]int{start}, escapeChars...)
				escapeCharsWidth = append([]int{width}, escapeCharsWidth...)
			}
			continue
		}

		if c == ' ' {
			break
		}
	}

	skippedWidth := 0
	for i, escapedChar := range escapeChars {
		if removeEscapedCharsFromResult {
			data = []byte(fmt.Sprintf("%s%s", data[0:escapedChar], data[escapedChar+escapeCharsWidth[i]:start]))
		}
		skippedWidth += escapeCharsWidth[i]
	}

	// If we haven't skipped any chars, we need to tell the scanner to advance one position extra
	// to skip over the comma separating the next key=value pair
	if skippedWidth == 0 {
		skippedWidth = 1
	}

	// If the value starts with a [, we remove it from the output
	if len(data) >= 1 && data[0] == '[' {
		return len(data[:start]) + 1, data[1:start]
	}

	return len(data[:start]) + 1, data[:start]
}
