/*
 * skogul, json parser
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
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
	"fmt"
	"os"
	"sync"

	"github.com/telenornms/skogul"
)

// DummyStore is a fake parser used to capture traffic Skogul doesn't
// understand, mainly for future development
type DummyStore struct {
	File   string `doc:"File name to write data to"`
	Append bool   `doc:"Append or overwrite. By default, the file is repeatedly overwritten to capture just 1 packet"`
	lock   sync.Mutex
}

// Parse accepts a byte slice of arbitrary data and returns nothing.
func (x *DummyStore) Parse(b []byte) (*skogul.Container, error) {
	container := skogul.Container{}
	container.Metrics = make([]*skogul.Metric, 0, 1)
	m := skogul.Metric{}
	m.Data = make(map[string]interface{})
	m.Metadata = make(map[string]interface{})
	container.Metrics = append(container.Metrics, &m)
	now := skogul.Now()
	m.Time = &now
	x.lock.Lock()
	defer x.lock.Unlock()
	var file *os.File
	var err error
	if finfo, lilerr := os.Stat(x.File); !os.IsNotExist(lilerr) && x.Append {
		file, err = os.OpenFile(x.File, os.O_APPEND|os.O_WRONLY, finfo.Mode())
	} else {
		file, err = os.Create(x.File)
	}
	if err != nil {
		return nil, fmt.Errorf("opening file failed: %w", err)
	}
	if file == nil {
		return nil, fmt.Errorf("open file returned nil-pointer")
	}
	defer file.Close()

	n, err := file.Write(b)
	if err != nil {
		return nil, fmt.Errorf("write error: %w", err)
	}
	if n != len(b) {
		return nil, fmt.Errorf("written bytes != received bytes")
	}
	return &container, nil
}
