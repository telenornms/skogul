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
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"net/url"
	"os"
)

// LineFile will keep reading File over and over again, assuming one
// collection per line. Best suited for pointing at a FIFO, which will
// allow you to 'cat' stuff to Skogul.
type LineFile struct {
	File    string
	Handler skogul.Handler
}

// Start never returns.
func (lf *LineFile) Start() error {
	for {
		f, err := os.Open(lf.File)
		if err != nil {
			log.Printf("Unable to open file: %s", err)
			return err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			bytes := scanner.Bytes()
			log.Printf("Read %s", bytes)
			m, err := lf.Handler.Parser.Parse(bytes)
			if err == nil {
				err = m.Validate()
			}
			if err != nil {
				log.Printf("Unable to parse JSON: %s", err)
				continue
			}
			for _, t := range lf.Handler.Transformers {
				t.Transform(&m)
			}
			lf.Handler.Sender.Send(&m)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading file: %s", err)
			return skogul.Error{Reason: "Error reading file"}
		}
	}
}

func init() {
	addAutoReceiver("fifo", NewLineFile, "Read from a FIFO on disk, reading one Skogul-formatted JSON per line. fifo:///var/skogul/foo")
}

// NewLineFile returns a LineFile receiver reading from the Path-element of
// the provided URL
func NewLineFile(ul url.URL, h skogul.Handler) skogul.Receiver {
	log.Printf("File: %s", ul.Path)
	return &LineFile{File: ul.Path, Handler: h}
}
