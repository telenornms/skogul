/*
 * skogul, linefile tests
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

package receiver_test

/*

import (
	"encoding/json"
	"fmt"
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"math/rand"
	"os"
	"syscall"
	"testing"
	"time"
)

func deleteFile(t *testing.T, file string) {
	err := os.Remove(file)
	if err != nil {
		t.Errorf("Failed to remove old test file %s: %v", file, err)
	}
}

func TestLinefile(t *testing.T) {
	rand.Seed(int64(time.Now().Nanosecond()))
	one := &(sender.Test{})

	file := fmt.Sprintf("%s/skogul-linefiletest-%d-%d", os.TempDir(), os.Getpid(), rand.Int())

	_, err := os.Stat(file)

	if err != nil && !os.IsNotExist(err) {
		t.Errorf("Error statting tmp file %s: %v", file, err)
		return
	}

	if !os.IsNotExist(err) {
		t.Errorf("File possibly exists already: %s", file)
		return
	}

	err = syscall.Mkfifo(file, 0600)

	if err != nil {
		t.Errorf("Unable to make fifo %s: %v", file, err)
		return
	}
	defer deleteFile(t, file)

	h := skogul.Handler{Sender: one, Parser: parser.JSON{}}
	rcv, err := receiver.New(fmt.Sprintf("fifo:///%s", file), h)
	if err != nil {
		t.Errorf("receiver.New() failed: %v", err)
		return
	}
	if rcv == nil {
		t.Errorf("receiver.New() returned err == nil, but also no receiver")
		return
	}

	go rcv.Start()
	b, err := json.Marshal(validContainer)
	if err != nil {
		t.Errorf("Failed to marshal container: %v", err)
		return
	}
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		t.Errorf("Unable to open file/fifo for writing: %v", err)
		return
	}
	defer func() {
		f.Close()
	}()
	f.WriteString(fmt.Sprintf("%s\n", b))
	time.Sleep(time.Duration(10 * time.Millisecond))
	if one.Received() != 1 {
		t.Errorf("Didn't receive thing on other end!")
	}
	f.WriteString(fmt.Sprintf("%s\n", b))
	time.Sleep(time.Duration(10 * time.Millisecond))
	if one.Received() != 2 {
		t.Errorf("Didn't receive thing on other end!")
	}
	one.Set(0)
	f.WriteString(fmt.Sprintf("bad idea♥\n"))
	if one.Received() != 0 {
		t.Errorf("Receive thing on other end despite bogus data")
	}
}
*/
