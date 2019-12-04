/*
 * skogul, timestamp transformer
 *
 * Copyright (c) 2019 Telenor Norge AS
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
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/telenornms/skogul"
)

var timestampLogger = skogul.Logger("transformer", "timestamp")

// Timestamp is the configuration for extracing a timestamp from inside the data
type Timestamp struct {
	Source []string `doc:"The source field of the timestamp"`
	Format string   `doc:"The format to use (default: RFC3339)"`
	Fail   bool     `doc:"Propagate errors back to the caller. Useful if the timestamp is required for the container."`
}

// Transform sets the timestamp of a set of metrics to the specified field
func (config *Timestamp) Transform(c *skogul.Container) error {
	for i, metric := range c.Metrics {

		obj, err := skogul.ExtractNestedObject(metric.Data, config.Source)
		timestamp, ok := obj[config.Source[len(config.Source)-1]].(string)

		if !ok {
			timestampLogger.Error("Failed to cast timestamp field to a string")
			if config.Fail {
				return skogul.Error{Reason: "Failed to cast timestamp field to a string"}
			}
		}

		format := parseTimestamp(config.Format)

		time, err := time.Parse(format, timestamp)
		if err != nil {
			timestampLogger.WithFields(logrus.Fields{
				"timestamp": timestamp,
				"format":    format,
			}).Error("Failed to parse timestamp")
			if config.Fail {
				return err
			}
		}

		c.Metrics[i].Time = &time
	}
	return nil
}

// parseTimestamp parses a timestamp format name into a timestamp format
// e.g. rfc3339 will be returned as "2006-01-02T15:04:05Z07:00"
func parseTimestamp(format string) string {
	switch strings.ToLower(format) {
	case "rfc3339", "iso8601": // 🙈
		return time.RFC3339
	default:
		return format
	}
}

// Verify will make sure the required fields are set, or set some sane defaults
func (config *Timestamp) Verify() error {
	if config.Source == nil {
		return skogul.Error{Reason: "Missing source field for timestamp transformer", Source: "timestamp transformer"}
	}
	if config.Format == "" {
		config.Format = time.RFC3339
	}
	return nil
}
