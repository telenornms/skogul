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
	"strings"

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
		if len(line) <= 2 {
			// Skip empty lines and lines without wrapping []
			continue
		}

		kvScanner := bufio.NewScanner(bytes.NewReader(line[1 : len(line)-1])) // trim leading [ and trailing ]
		// ToDo: Support Example 2 from https://tools.ietf.org/html/rfc5424#section-6.3.5
		kvScanner.Split(splitKeyValuePairs)

		// First entry in the line is a single value, not a key/value pair
		// Extract it, and store it.
		ok := kvScanner.Scan()
		if !ok {
			sdLog.Warning("Failed to parse structured data line")
			continue
		}
		structuredDataID := kvScanner.Text()

		metadata := make(map[string]interface{})
		metadata["sd-id"] = structuredDataID

		data := make(map[string]interface{})
		for {
			canContinue := kvScanner.Scan()

			tag := strings.Trim(kvScanner.Text(), "\u0000")
			tagValue := strings.SplitN(tag, "=", 2)

			if len(tagValue) != 2 {
				break
			}

			paramName := tagValue[0]
			paramValue := tagValue[1][1 : len(tagValue[1])-1] // remove leading and trailing "s

			// @ToDo: Support multiple paramName with different paramValue
			// if the value already exists, replace it with an array ?

			data[paramName] = paramValue
			sdLog.WithField("tag", paramName).WithField("val", paramValue).Debug("Got :)")

			if !canContinue {
				break
			}
		}
		metrics = append(metrics, &skogul.Metric{
			Time:     &timestamp,
			Data:     data,
			Metadata: metadata,
		})
	}
	if len(metrics) == 0 {
		sdLog.WithField("lines", len(lines)).Warnf("RFC5424/Structured Data parser failed to parse any of the %d lines", len(lines))
		return nil, skogul.Error{Reason: "Failed to parse RFC5424 lines", Source: "structured_data-parser"}
	}
	return metrics, nil
}

// splitKeyValuePairs splits a section (tag key=value pairs or field key=value pairs)
// into key=value pairs, honoring escape rules as per the influx line protocol.
// A key=value pair is split on a non-escaped space.
func splitKeyValuePairs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	fieldWidth, newData := influxLineParser(data, ' ', true)

	returnChars := len(newData)

	if returnChars == len(data) {
		// EOF, return with what we have left
		return returnChars, newData[:returnChars], nil
	}

	// Skip the trailing comma between each key=value pair, but still advance counter
	return fieldWidth, newData[:returnChars], nil
}
