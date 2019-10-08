/*
 * skogul, log receiver
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

package receiver

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/KristianLyng/skogul"
)

// Log redirects Skoguls log buffer to a handler
type Log struct {
	Echo    bool              `doc:"Logs are also echoed to stdout."`
	Handler skogul.HandlerRef `doc:"Reference to a handler where the data is sent. Parser will be overwritten."`
}

// logContainer returns a container representing the log message in s.
func logContainer(s string) (*skogul.Metric, error) {
	now := time.Now()
	m := skogul.Metric{}
	m.Metadata = make(map[string]interface{})
	m.Data = make(map[string]interface{})
	m.Time = &now
	splat := strings.Split(s, ":")
	if len(splat) < 3 {
		return nil, skogul.Error{Source: "log receiver", Reason: "Log message does not follow expected format"}
	}
	m.Metadata["file"] = splat[0]
	m.Metadata["line"] = splat[1]
	m.Data["message"] = strings.Join(splat[2:], ":")[1:]
	return &m, nil
}

// Parse implements a skogul.Parser logic by parsing the byte array as
// received by log. Each non-empty line results in a single metric.
func (lg *Log) Parse(b []byte) (*skogul.Container, error) {
	c := skogul.Container{}
	c.Metrics = make([]*skogul.Metric, 0)
	cpy := string(b)
	for _, line := range strings.Split(cpy, "\n") {
		if line == "" {
			continue
		}
		if lg.Echo {
			fmt.Println(line)
		}
		m, err := logContainer(line)
		if err != nil {
			fmt.Printf("Failed to parse log line, error: %v, log line: %s\n", err, line)
			return nil, err
		}
		c.Metrics = append(c.Metrics, m)
	}
	return &c, nil
}

// Write splits the input on line-shift, assumes it
// follows the log format and parses it into a container to be sent to a
// sender.
//
// One issue we have is how to report errors. If we get errors, we do not
// want to trigger a feedback loop, so we can't use log.P... internally.
func (lg *Log) Write(p []byte) (n int, err error) {
	lg.Handler.H.Handle(p)
	return len(p), nil
}

// Start acquires the standard log writer, sets appropriate flags and never
// returns.
func (lg *Log) Start() error {
	lg.Handler.H.Parser = lg
	log.SetFlags(log.Lshortfile)
	log.SetOutput(lg)
	for {
		time.Sleep(time.Hour * 1337)
	}
}
