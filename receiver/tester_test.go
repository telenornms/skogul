/*
 * skogul, test receiver tests (... I know)
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
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"testing"
	"time"
)

func TestTester_stack(t *testing.T) {
	one := &(sender.Test{})
	h := skogul.Handler{Sender: one, Parser: parser.JSON{}}
	rcv := receiver.Tester{Metrics: 10, Values: 5, Threads: 2, Handler: skogul.HandlerRef{H: &h}}
	go rcv.Start()

	zzz := 300 * time.Millisecond
	time.Sleep(time.Duration(zzz))

	// With atomic covermode and race condition testing, this is pretty
	// slow, hence the modest values
	got := one.Received()
	want := uint64(500)
	if got < want {
		t.Errorf("Less than %d received events after %v of tester running. Got %d", want, zzz, got)
	}
}

func TestTester_auto(t *testing.T) {
	one := &(sender.Test{})
	h := skogul.Handler{Sender: one, Parser: parser.JSON{}}
	x, err := receiver.New("test://", h)
	if err != nil {
		t.Errorf("receiver.New(\"test://\") gave error: %v", err)
	}
	if x == nil {
		t.Errorf("no receiver created with default test:// url")
	}

	x, err = receiver.New("test:///?values=fem", h)
	if err == nil {
		t.Errorf("receiver.New(\"test:///?values=fem\") gave no error.")
	}
	if x != nil {
		t.Errorf("Receiver created with test:///?values=fem url: %v", x)
	}
	x, err = receiver.New("test:///?metrics=fem", h)
	if err == nil {
		t.Errorf("receiver.New(\"test:///?metrics=fem\") gave no error.")
	}
	if x != nil {
		t.Errorf("Receiver created with test:///?metrics=fem url: %v", x)
	}
	x, err = receiver.New("test:///?delay=1s", h)
	if err != nil {
		t.Errorf("receiver.New(\"test:///?delay=1s\") gave error: %v", err)
	}
	if x == nil {
		t.Errorf("no receiver created with default test:// url")
	}
	go x.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
	if one.Received() < 1 {
		t.Errorf("receiver.New(\"test:///?delay=1s\"), x.Start() failed to receive data. Expected some data, got 0.")
	}
}
*/
