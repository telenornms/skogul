/*
 * skogul, sender tests
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

package sender_test

import (
	"bytes"
	"github.com/KristianLyng/skogul/sender"
	"log"
	"testing"
)

var logBuffer bytes.Buffer

func init() {
	log.SetOutput(&logBuffer)
}
func testGeneric(t *testing.T, in string) {
	x, err := sender.New(in)
	if err != nil {
		t.Errorf("New(%s): Failed to create sender automatically: %v", in, err)
	}
	if x == nil {
		t.Errorf("New(%s) returned no sender", in)
	}
}

func testGenericNegative(t *testing.T, in string) {
	x, err := sender.New(in)
	if err == nil {
		t.Errorf("New(%s): Supposed to fail, but worked.", in)
	}
	if x != nil {
		t.Errorf("New(%s): Should fail, but returned sender: %v", in, x)
	}
}

func TestNew_debug(t *testing.T) {
	testGeneric(t, "debug://")
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

func TestNew_mnr(t *testing.T) {
	ok := []string{
		"mnr://localhost",
		"mnr://localhost/",
		"mnr://localhost:8080",
		"mnr://[::1]",
		"mnr://[::1]",
		"mnr://[::1]/foobar"}
	for _, url := range ok {
		testGeneric(t, url)
	}
	bad := []string{
		"mnr://[::",
	}
	for _, url := range bad {
		testGenericNegative(t, url)
	}
}
