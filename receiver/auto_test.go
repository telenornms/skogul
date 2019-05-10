/*
 * skogul, receiver tests
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
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/sender"
	"net/url"
	"testing"
)

var h skogul.Handler

func init() {
	h = skogul.Handler{Sender: sender.Debug{}}
}

func testGeneric(t *testing.T, in string) {
	x, err := New(in, h)
	if err != nil {
		t.Errorf("New(%s): Failed to create receiver automatically: %v", in, err)
	}
	if x == nil {
		t.Errorf("New(%s) returned no receiver", in)
	}
}

func testGenericNegative(t *testing.T, in string) {
	x, err := New(in, h)
	if err == nil {
		t.Errorf("New(%s): Supposed to fail, but worked.", in)
	}
	if x != nil {
		t.Errorf("New(%s): Should fail, but returned receiver: %v", in, x)
	}
}

func TestNew_test(t *testing.T) {
	testGeneric(t, "test:///?threads=15&values=1")
	testGeneric(t, "test:///?threads=15")
	testGeneric(t, "test://")
	testGenericNegative(t, "test:///?threads=x/y")
}

func TestNew_nonexistent(t *testing.T) {
	testGenericNegative(t, "nonexistent")
}

func TestNew_http(t *testing.T) {
	ok := []string{
		"http://localhost",
		"http://localhost/",
		"http://localhost:8080",
		"http://[::1]",
		"http://[::1]",
		"http://[::1]/foo/bar"}
	for _, url := range ok {
		testGeneric(t, url)
	}
	bad := []string{
		"http://[::",
	}
	for _, url := range bad {
		testGenericNegative(t, url)
	}
}

func intTest(t *testing.T) {
	defer func() {
		recover()
	}()
	d := func(url url.URL, h skogul.Handler) skogul.Receiver { return nil }
	addAutoReceiver("http", d, "foo")
	t.Errorf("addAutoReceiver() on existing sender didn't fail?")
}

func Test_addAutoReceiver(t *testing.T) {
	intTest(t)
}
