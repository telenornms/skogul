/*
 * skogul, file writer
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Håkon Solbjørg <Hakon.Solbjorg@telenor.com>
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

package sender

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/telenornms/skogul"
)

var fileLog = skogul.Logger("sender", "file")

const newLineChar byte = 10

type File struct {
	Path   string `doc:"Absolute path to file to write"`
	Append bool   `doc:"Whether to append to the file when starting. If false, will empty file before starting writes. Default: false"`
	ok     bool
	once   sync.Once
	f      *os.File
	c      chan []byte
}

func (f *File) init() {
	fileLog.Debug("Initializing File sender")

	var err error
	var file *os.File

	// Open file for append-only if it already exists and config says to append
	if finfo, err := os.Stat(f.Path); !os.IsNotExist(err) && f.Append {
		fileLog.Trace("File exists, let's open it for writing")
		file, err = os.OpenFile(f.Path, os.O_APPEND|os.O_WRONLY, finfo.Mode())
	} else {
		// Otherwise, create the file (which will truncate it if it already exists)
		fileLog.Trace("Creating file since it doesn't exist or we don't want to append to it")
		file, err = os.Create(f.Path)
	}
	if err != nil {
		fileLog.WithError(err).Errorf("Failed to open '%s'", f.Path)
		f.ok = false
		return
	}

	// startChan() handles closing the file as this function returns
	// and consequently would close the file
	f.f = file

	// Listening to a channel is blocking so we have
	// to start the channel listening in a goroutine
	// so that init() doesn't block.
	go f.startChan()

	f.ok = true
}

func (f *File) startChan() error {
	fileLog.Trace("Starting file writer channel")
	// Making sure we close the file if this function exits
	defer f.f.Close()
	f.c = make(chan []byte)
	for b := range f.c {
		written, err := f.f.Write(append(b, newLineChar))
		if err != nil {
			f.ok = false
			fileLog.WithError(err).Errorf("Failed to write to file. Wrote %d of %d bytes", written, len(b))
			return err
		}
		f.f.Sync()
	}
	fileLog.Warning("File writer chan closed, not handling any more writes!")
	return nil
}

// Send receives a skogul container and writes it to file.
func (f *File) Send(c *skogul.Container) error {
	f.once.Do(func() {
		f.init()
	})

	if !f.ok {
		e := skogul.Error{Reason: "File sender not in OK state", Source: "file sender"}
		fileLog.WithError(e).Error("Failied to initialize file sender, or an error occured in runtime")
		return e
	}

	b, err := json.Marshal(*c)

	if err != nil {
		fileLog.WithError(err).Error("Failed to marshal container data to json")
		return err
	}

	f.c <- b

	return nil
}

// Verify checks that the configuration options are set appropriately
func (f *File) Verify() error {
	fileLog.Debug("Verified file sender")
	return nil
}
