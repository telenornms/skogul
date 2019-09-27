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
	Handler skogul.HandlerRef `doc:"Reference to a handler where the data is sent."`
}

// logContainer returns a container representing the log message in s.
func logContainer(s string) (*skogul.Container, error) {
	c := skogul.Container{}
	now := time.Now()
	c.Metrics = make([]*skogul.Metric, 1)
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
	c.Metrics[0] = &m
	return &c, nil
}

// Write splits the input on line-shift, assumes it
// follows the log format and parses it into a container to be sent to a
// sender.
//
// One issue we have is how to report errors. If we get errors, we do not
// want to trigger a feedback loop, so we can't use log.P... internally.
func (lg *Log) Write(p []byte) (n int, err error) {
	cpy := string(p)
	for _, line := range strings.Split(cpy, "\n") {
		if line == "" {
			continue
		}
		if lg.Echo {
			fmt.Println(line)
		}
		c, err := logContainer(line)
		if err == nil {
			lg.Handler.H.Sender.Send(c)
		} else {
			fmt.Printf("Log receiver failed to parse a log message into a container. The internal error was: %v\n", err)
			fmt.Printf("Original message was: %s\n", line)
		}
	}
	return len(p), nil
}

// Start acquires the standard log writer, sets appropriate flags and never
// returns.
func (lg *Log) Start() error {

	log.SetFlags(log.Lshortfile)
	log.SetOutput(lg)
	du, _ := time.ParseDuration("1h")
	for {
		time.Sleep(du)
	}
}
