/*
 * skogul, influxdb parser
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
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/telenornms/skogul"
)

// InfluxDB provides a byte sequence parser for the InfluxDB Line Protocol
// https://docs.influxdata.com/influxdb/v1.7/write_protocols/line_protocol_tutorial/
type InfluxDB struct{}

// InfluxDBLineProtocol is a struct with the same data types as defined in the InfluxDB
// Line Protocol; namely the measurement name, a set of tags, a set of fields and a timestamp.
type InfluxDBLineProtocol struct {
	measurement string
	tags        map[string]interface{}
	fields      map[string]interface{}
	timestamp   time.Time
}

var influxLogger = skogul.Logger("parser", "influxdb")

// Parse marshals a byte sequence of InfluxDB line protocol into a skogul container
func (influxdb InfluxDB) Parse(bytes []byte) (*skogul.Container, error) {
	s := string(bytes)

	// Do we receive data with \r\n?
	lines := strings.Split(s, "\n")

	container := skogul.Container{
		Metrics: make([]*skogul.Metric, 0),
	}

	errors := make([]skogul.Error, 0)

	for i, l := range lines {
		line := strings.TrimSpace(l)
		if len(strings.TrimSpace(line)) == 0 {
			// Skip empty lines
			continue
		}
		influxLine := InfluxDBLineProtocol{}
		if err := influxLine.ParseLine(line); err != nil {
			// skogul.Error.Source shows index in list and actual line that failed parsing.
			// Skip either? Both?
			errors = append(errors, skogul.Error{Source: fmt.Sprintf("%d-'%s'", i, line), Reason: "Failed to parse influx line", Next: err})
			influxLogger.WithError(err).Error("Failed to parse influx line protocol")
			continue
		}

		container.Metrics = append(container.Metrics, influxLine.Metric())
	}

	if len(errors) > 0 {
		return &container, skogul.Error{
			Source: "parser-influxdb",
			Reason: fmt.Sprintf("One or more influxdb line protocol parse failures. Returning %d successful parses and skipping %d errors.", len(container.Metrics), len(errors)),
		}
	}

	return &container, nil
}

// ParseLine parses a single line into an internal InfluxDBLineProtocol
func (line *InfluxDBLineProtocol) ParseLine(s string) error {
	// Let's see if we can find a , which is not escaped
	// That'll be our measurement name.
	prev := ""
	for i, c := range string(s) {
		if (c == ',' || c == ' ') && prev != "\\" {
			line.measurement = s[:i]
			break
		}
		prev = fmt.Sprint(c)
	}

	if line.measurement == "" {
		return skogul.Error{Source: "parser-influxdb", Reason: "Could not find a measurement name"}
	}

	// skip the comma trailing measurement name
	sections := s[len(line.measurement)+1:]

	scanner := bufio.NewScanner(strings.NewReader(sections))
	scanner.Split(splitSections)

	canContinue := scanner.Scan()

	if !canContinue && scanner.Err() != nil {
		return skogul.Error{Source: "parser-influxdb", Reason: "Scanner cannot continue after first scan, received an error", Next: scanner.Err()}
	}

	tags := scanner.Text()

	canContinue = scanner.Scan()

	if !canContinue && scanner.Err() != nil {
		return skogul.Error{Source: "parser-influxdb", Reason: "Scanner cannot continue after second scan, received an error", Next: scanner.Err()}
	}

	fields := scanner.Text()

	canContinue = scanner.Scan()

	if !canContinue && scanner.Err() != nil {
		return skogul.Error{Source: "parser-influxdb", Reason: "Scanner cannot continue after third scan, received an error", Next: scanner.Err()}
	}

	// If we get a valid length here we have a value in the timestamp section
	hasTimestamp := len(scanner.Text()) > 0

	if hasTimestamp {
		timestamp := scanner.Text()
		nsTime, err := strconv.ParseInt(timestamp, 0, 64)
		if err != nil {
			return skogul.Error{
				Source: "parserinfluxdb",
				Reason: "Failed to parse time for influxdb line",
				Next:   scanner.Err()}
		}
		line.timestamp = time.Unix(0, nsTime)
	} else {
		// Create own timestamp if it doesn't exist in the source line
		line.timestamp = skogul.Now()
	}

	// Parse tags and fields

	line.tags = make(map[string]interface{})

	tagScanner := bufio.NewScanner(strings.NewReader(tags))
	tagScanner.Split(splitInfluxKeyValuePairs)
	for {
		canContinue := tagScanner.Scan()

		tag := strings.Trim(tagScanner.Text(), "\u0000")
		tagValue := strings.SplitN(tag, "=", 2)

		if len(tagValue) != 2 {
			break
		}

		line.tags[tagValue[0]] = tagValue[1]

		if !canContinue {
			break
		}
	}

	line.fields = make(map[string]interface{})

	fieldScanner := bufio.NewScanner(strings.NewReader(fields))
	fieldScanner.Split(splitInfluxKeyValuePairs)
	for {
		canContinue := fieldScanner.Scan()

		a := fieldScanner.Text()

		field := strings.Trim(a, "\u0000")
		fieldValue := strings.Split(field, "=")

		if len(fieldValue) != 2 {
			break
		}

		line.fields[fieldValue[0]] = parseFieldValue(fieldValue[1])

		if !canContinue {
			break
		}
	}

	return nil
}

// Metric converts an internal InfluxDBLineProtocol struct to a skogul.Metric
func (line *InfluxDBLineProtocol) Metric() *skogul.Metric {
	line.tags["measurement"] = line.measurement

	metric := skogul.Metric{
		Time:     &line.timestamp,
		Metadata: line.tags,
		Data:     line.fields,
	}

	return &metric
}

func parseFieldValue(value string) interface{} {

	// If the last char is an 'i' and the rest is numeric, this is an integer
	if value[len(value)-1:] == "i" {
		if i, err := strconv.ParseInt(value[0:len(value)-1], 0, 64); err == nil {
			return i
		}
	}

	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	if value == "t" || value == "T" || value == "true" || value == "True" || value == "TRUE" {
		return true
	}

	if value == "f" || value == "F" || value == "false" || value == "False" || value == "FALSE" {
		return false
	}

	if value[0] == '"' && value[len(value)-1] == '"' {
		return value[1 : len(value)-1]
	}

	return value
}

// splitFieldFunc is a SplitFunc for Scanner which splits a string into
// influxdb line protocol sections for tag key=value pairs and field key=value pairs.
// Sections are split on a non-escaped space character, and we retain all escaped
// characters and let the next splitter handle them.
func splitSections(data []byte, atEOF bool) (advance int, token []byte, err error) {

	fieldWidth, newData := influxLineParser(data, ' ', false)

	returnChars := len(newData)

	if returnChars == len(data) {
		// EOF, return with what we have left
		return returnChars, newData[:returnChars], nil
	}

	// Skip the trailing comma between each key=value pair, but still advance counter
	return fieldWidth, newData[:returnChars], nil
}

// splitInfluxKeyValuePairs splits a section (tag key=value pairs or field key=value pairs)
// into key=value pairs, honoring escape rules as per the influx line protocol.
// A key=value pair is split on a non-escaped comma.
func splitInfluxKeyValuePairs(data []byte, atEOF bool) (advance int, token []byte, err error) {

	fieldWidth, newData := influxLineParser(data, ',', true)

	returnChars := len(newData)

	if returnChars == len(data) {
		// EOF, return with what we have left
		return returnChars, newData[:returnChars], nil
	}

	// Skip the trailing comma between each key=value pair, but still advance counter
	return fieldWidth, newData[:returnChars], nil
}

// influxLineParser parses part of an influxdb line protocol line and tells the
// calling scanner how far it should advance (pretty similar to the splitFunc API).
// The character to split on is passed to the function, and would usually be
// a space or a comma character, as those are what's used to split
// the influx line protocol section or key=value pair from each other.
// A boolean flag decides whether or not escape characters should remain in the output
// or have their prepending escape character removed.
func influxLineParser(data []byte, sectionBreak rune, removeEscapedCharsFromResult bool) (int, []byte) {
	openQuote := false
	escape := false
	escapeChars := make([]int, 0)
	escapeCharsWidth := make([]int, 0)

	start := 0
	previousWidth := 0
	for width := 0; start < len(data); start += width {
		var c rune
		c, width = utf8.DecodeRune(data[start:])

		if escape {
			escape = false
			if c != 'x' && c != 'X' && c != '0' && c != 'u' && c != 'U' {
				// \x is hex, so let's keep the \ and the x so that a consumer can
				// parse the value themselves. Let's also do the same for decimals (\0) and unicode (\u).
				if removeEscapedCharsFromResult {
					escapeChars = append([]int{start - previousWidth}, escapeChars...)
					escapeCharsWidth = append([]int{previousWidth}, escapeCharsWidth...)
				}
			}
			continue
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
			previousWidth = width
			continue
		}

		if c == sectionBreak {
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

	return len(data[:start]) + 1, data[:start]
}
