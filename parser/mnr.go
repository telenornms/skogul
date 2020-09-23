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
package parser

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/telenornms/skogul"
)

var mnrLog = skogul.Logger("parser", "mnr")

// MNR supports parsing MNR data through the Parse() function
type MNR struct {
	ExtractFieldName bool     `doc:"Extract field name from 'name' parameter in MNR properties."`
	DefaultFieldName string   `doc:"Name of field to store variable value in, in the case that 'name' is not present."`
	ParseAsString    []string `doc:"Parse these properties as string (do not try to parse their value)."`
	StoreVariable    bool     `doc:"Store the uniquely generated variable as a data field."`
	once             sync.Once
}

const mnrSeparator byte = 9 // tab

// init prints information about parser configuration
func (mnr *MNR) init() error {
	if len(mnr.ParseAsString) > 0 {
		mnrLog.Debugf("MNR parser configured with parsing fields as string: %v", mnr.ParseAsString)
	}
	return nil
}

// Parse converts data from MNR into a skogul Container
func (mnr *MNR) Parse(bytes []byte) (*skogul.Container, error) {
	mnr.once.Do(func() {
		err := mnr.init()
		if err != nil {
			mnrLog.WithError(err).Error("Failed to initialise MNR parser")
		}
	})

	metrics, err := mnr.mnrParseData(bytes)
	if err != nil {
		return nil, err
	}

	return &skogul.Container{
		Metrics: metrics,
	}, nil
}

// mnrParseData takes the raw input and parses it
// this takes care of splitting input on newlines etc
func (mnr *MNR) mnrParseData(data []byte) ([]*skogul.Metric, error) {
	lines := bytes.Split(data, []byte("\n"))

	metrics := make([]*skogul.Metric, 0)

	for _, l := range lines {
		line := bytes.TrimSpace(l)
		if len(line) == 0 {
			// Skip empty lines
			continue
		}

		metric, err := mnr.mnrParseLine(line)
		if err != nil {
			// If we get an error on this line we will continue
			// in hopes of having other lines working successfully.
			mnrLog.WithError(err).Error("Failed to parse MNR line")
			continue
		}
		metrics = append(metrics, metric)
	}
	if len(metrics) == 0 {
		mnrLog.WithField("lines", len(lines)).Warnf("MNR parser failed to parse any of the %d lines", len(lines))
		return nil, skogul.Error{Reason: "Failed to parse MNR lines", Source: "mnr-parser"}
	}
	return metrics, nil
}

func (mnr *MNR) mnrParseLine(data []byte) (*skogul.Metric, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(mnrSplitLineFunc)

	metric := skogul.Metric{}
	metric.Metadata = make(map[string]interface{})
	metric.Metadata["mnr_changed"] = false
	metric.Metadata["mnr_deleted"] = false

	// First scan should give timestamp or changed symbol
	ok := scanner.Scan()
	if !ok {
		return nil, skogul.Error{Reason: "Failed to extract first value from a MnR line"}
	}
	timestamp := scanner.Text()

	// MNR might prepend the whole line with +r or +d, so let's handle these
	if timestamp[:1] == "+" {
		changed := timestamp
		metric.Metadata["mnr_tag"] = changed

		if changed[1:] == "r" {
			metric.Metadata["mnr_changed"] = true
		}
		if changed[1:] == "d" {
			metric.Metadata["mnr_deleted"] = true
		}

		// The next value should be a timestamp
		ok := scanner.Scan()
		if !ok {
			return nil, skogul.Error{Reason: "Failed to extract timestamp as second value from a MnR line"}
		}
		timestamp = scanner.Text()
	}

	tint, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, skogul.Error{Reason: "Failed to convert string to integer for timestamp", Source: "mnr-parser"}
	}
	ts := time.Unix(tint, 0)
	metric.Time = &ts

	// Fetch MNR group
	ok = scanner.Scan()
	if !ok {
		return nil, skogul.Error{Reason: "Failed to extract MNR group", Source: "mnr-parser"}
	}
	metric.Metadata["group"] = scanner.Text()

	// Time to fetch the actual variable and its value
	metric.Data = make(map[string]interface{})
	ok = scanner.Scan()
	if !ok {
		return nil, skogul.Error{Reason: "Failed to extract MNR variable name", Source: "mnr-parser"}
	}
	variable := scanner.Text()

	if mnr.StoreVariable {
		metric.Data["variable"] = variable
	}

	ok = scanner.Scan()
	if !ok {
		return nil, skogul.Error{Reason: "Failed to extract MNR variable value", Source: "mnr-parser"}
	}
	val := parseMNRFieldValue(scanner.Text())

	// Parsing properties
	for scanner.Scan() {
		pair := strings.SplitN(scanner.Text(), "=", 2)
		if len(pair) != 2 {
			// Skip because we didn't get a 'key=value' pair as we expected
			continue
		}
		key := pair[0]
		parseType := true
		for _, field := range mnr.ParseAsString {
			if key == field {
				parseType = false
				break
			}
		}
		if parseType {
			metric.Data[key] = parseMNRFieldValue(pair[1])
		} else {
			metric.Data[key] = string(pair[1])
		}
	}

	// If we have extracted a 'name' property, that tells us the actual
	// name of the variable. We'll store this as the field for the value.
	if mnr.ExtractFieldName && metric.Data["name"] != nil {
		metric.Data[metric.Data["name"].(string)] = val
	} else {
		// If the 'name' property didn't exist, write the value to this key.
		if mnr.DefaultFieldName != "" {
			// The variable name is by default pretty hard to identify
			// programmatically, so we can write the value here too for convenience.
			metric.Data[mnr.DefaultFieldName] = val
		}

		// Notify that we weren't able to extract the name.
		metric.Data["failed_name_extract"] = true
	}

	return &metric, nil
}

// parseMNRFieldValue tries to convert a value into a non-string type
// such as integers or floats. If no parse succeeds, the same string-value
// is returned as the value.
func parseMNRFieldValue(value string) interface{} {
	if i, err := strconv.ParseInt(value, 0, 64); err == nil {
		return i
	}

	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	return value
}

func mnrSplitLineFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, nil
	}

	advance = 0
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width := utf8.DecodeRune(data[start:])
		// const separator as a rune instead ?
		if byte(r) == mnrSeparator {
			break
		}
		start += width
	}
	advance = start

	// Cleaner way?
	// If we're at the end of the data, return the length of the data
	if start+1 > len(data) {
		return start, data[:advance], nil
	}

	// Otherwise, we skip the next byte (which is a separator)
	return start + 1, data[:advance], nil
}
