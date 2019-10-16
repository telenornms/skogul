/*
 * skogul, tests
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

package skogul_test

import (
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/parser"
	"github.com/telenornms/skogul/sender"
	"testing"
)

func TestHandler(t *testing.T) {
	h1 := skogul.Handler{}
	h2 := skogul.Handler{Parser: parser.JSON{}}
	h3 := skogul.Handler{Parser: parser.JSON{}, Transformers: []skogul.Transformer{}}
	h4 := skogul.Handler{Parser: parser.JSON{}, Transformers: []skogul.Transformer{}, Sender: &(sender.Test{})}
	h5 := skogul.Handler{Parser: parser.JSON{}, Transformers: []skogul.Transformer{nil}, Sender: &(sender.Test{})}

	err := h1.Verify()
	if err == nil {
		t.Errorf("Handler verification didn't spot empty handler")
	}
	err = h2.Verify()
	if err == nil {
		t.Errorf("Handler verification didn't spot parser-only handler")
	}
	err = h3.Verify()
	if err == nil {
		t.Errorf("Handler verification didn't spot parser-and-transformers-only handler")
	}

	err = h4.Verify()
	if err != nil {
		t.Errorf("Supposedly valid handler actually failed verification: %v", err)
	}
	err = h5.Verify()
	if err == nil {
		t.Errorf("Handler verification didn't spot nil-transformer")
	}
}

func TestContainer(t *testing.T) {
	orig := skogul.Error{Source: "int", Reason: "fordi"}
	c := orig.Container()
	if c.Metrics[0] == nil {
		t.Errorf("Failed to get a metric from errror.Container")
	}
	if c.Metrics[0].Metadata["source"] != "int" {
		t.Errorf("error.Container() returned invalid source. Wanted %s got %s", "int", c.Metrics[0].Metadata["source"])
	}
	want := "fordi"
	got := c.Metrics[0].Data["reason"]
	if want != got {
		t.Errorf("error.Container() returned unexpected data/reason. Wanted %s got %s", want, got)
	}
}

func TestAssert(t *testing.T) {
	skogul.Assert(true)
	skogul.Assert(1+1 != 0)
	skogul.Assert(t != nil)
	skogul.Assert(true, "foo")
}

func TestAssert_fail(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered from assert: ", r)
		}
	}()
	skogul.Assert(false)
	t.Errorf("skogul.Error(false,\"test\") called, but execution continued.")
}

func TestAssert_fail_arg(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("Recovered from assert: ", r)
		}
	}()
	skogul.Assert(false, "something")
	t.Errorf("skogul.Error(false,\"test\") called, but execution continued.")
}
