/*
 * skogul, linefile receiver
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
	"bufio"
	"github.com/telenornms/skogul"
	"log"
	"os"
	"time"
)

// LineFile will keep reading File over and over again, assuming one
// collection per line. Best suited for pointing at a FIFO, which will
// allow you to 'cat' stuff to Skogul.
type LineFile struct {
	File    string            `doc:"Path to the fifo or file from which to read from repeatedly."`
	Handler skogul.HandlerRef `doc:"Handler used to parse and transform and send data."`
	Delay   skogul.Duration   `doc:"Delay before re-opening the file, if any."`
}

// Common routine for both fifo and stdin
func (lf *LineFile) read() error {
	f, err := os.Open(lf.File)
	if err != nil {
		log.Printf("Unable to open file: %s", err)
		return err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		if err := lf.Handler.H.Handle(bytes); err != nil {
			log.Printf("Failed to send metric: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file: %s", err)
		return skogul.Error{Reason: "Error reading file"}
	}
	return nil
}

// Start never returns.
func (lf *LineFile) Start() error {
	for {
		lf.read()
		if lf.Delay.Duration != 0 {
			time.Sleep(lf.Delay.Duration)
		}
	}
}

// File reads from a FILE, a single JSON object per line, and
// exits at EOF.
type File struct {
	File    string            `doc:"Path to the file to read from once."`
	Handler skogul.HandlerRef `doc:"Handler used to parse, transform and send data."`
	lf      LineFile
}

// Start reads a file once, then returns.
func (s *File) Start() error {
	s.lf.File = s.File
	s.lf.Handler = s.Handler
	return s.lf.read()
}

// Stdin reads from /dev/stdin
type Stdin struct {
	Handler skogul.HandlerRef `doc:"Handler used to parse, transform and send data."`
	lf      LineFile
}

// Start reads from stdin until EOF, then returns
func (s *Stdin) Start() error {
	s.lf.File = "/dev/stdin"
	s.lf.Handler = s.Handler
	return s.lf.read()
}
