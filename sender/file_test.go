/*
 * skogul, file writer tests
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
 * Author(s):
 *  - Roshini Narasimha Raghavan <roshiragavi@gmail.com>
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

package sender_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"syscall"
	"testing"
	"time"

	skogul "github.com/telenornms/skogul"
	sender "github.com/telenornms/skogul/sender"
)

func createConf() {

}

// createContainer is a simple helper func which
// creates a skogul.Container with some data
func createContainer() skogul.Container {
	meta := make(map[string]interface{})
	meta["foo"] = "bar"
	data := make(map[string]interface{})
	data["baz"] = "qux"

	metric := skogul.Metric{
		Metadata: meta,
		Data:     data,
	}
	metrics := make([]*skogul.Metric, 0)
	metrics = append(metrics, &metric)
	return skogul.Container{
		Metrics: metrics,
	}
}

func TestWriteToNonExistingFile(t *testing.T) {
	filename := "skogul-file-sender-nonexisting-file.txt"
	tmpdir := os.TempDir()
	path := path.Join(tmpdir, filename)

	// Ensure file does not exist already
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// Already exists..
		// Let's assume it's safe to remove ?
		os.Remove(path)
	}

	// Now let's initialize a config which writes to that file which does not exist
	sender := &sender.File{
		File:   path,
		Append: false,
	}

	c := createContainer()
	sender.Send(&c)

	// Since the write is done by a goroutine
	// we have to make sure it is properly
	// flushed before we try to read it back
	time.Sleep(time.Second)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
		return
	}

	var j map[string]interface{}
	err = json.Unmarshal(b, &j)
	if err != nil {
		t.Error(err)
		return
	}
	// Assume that if we managed to read the file
	// and unmarshal the contents to JSON
	// the write succeeded.
}

func TestAppendToExistingFile(t *testing.T) {
	filename := "skogul-file-sender-existing-file-append.txt"
	path := path.Join(os.TempDir(), filename)

	f, err := os.Create(path)
	if err != nil {
		t.Error(err)
		return
	}

	f.Write([]byte("some data\n"))
	f.Sync()

	time.Sleep(time.Second)

	sender := &sender.File{
		File:   path,
		Append: true,
	}

	c := createContainer()
	sender.Send(&c)

	// Since the write is done by a goroutine
	// we have to make sure it is properly
	// flushed before we try to read it back
	time.Sleep(time.Second)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
		return
	}
	str := string(b)
	if !strings.Contains(str, "some data") {
		t.Errorf("Test file does not contain test string 'some data', was it overwritten? Contents: %s", str)
	}
}

func TestSignals(t *testing.T) {
	filename := "skogul-file-1.txt"
	path := path.Join(os.TempDir(), filename)

	f, err := os.Create(path)
	if err != nil {
		t.Error(err)
		return
	}

	f.Write([]byte("some data\n"))
	f.Sync()

	time.Sleep(time.Second)

	sender := &sender.File{
		File:   path,
		Append: false,
	}

	c := createContainer()

	sender.Send(&c)

	os.Rename(path, path+".old")

	sender.Send(&c)

	// the file should not exists since we renamed. If it continues to exists, then fail the test.
	if _, err := os.Stat(path); err == nil {
		t.Error(err)
		t.FailNow()
	}

	// Since the write is done by a goroutine
	// we have to make sure it is properly
	// flushed before we try to read it back
	time.Sleep(time.Second)
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Error(err)
		return
	}

	// provides an interrup signal to the current process handling the file.
	if err := proc.Signal(syscall.SIGHUP); err != nil {
		t.Error(err)
		return
	}

	// give enough time for the sighup signal to be handled and send the container again. This should reopen the file and write the values sent by the container.
	time.Sleep(time.Second * 4)
	sender.Send(&c)

	// If the file remains closed, then the following commands will fail the test.
	_, err = os.Stat(path)
	if err != nil {
		t.Error(err)
		return
	}

}
