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
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/sender"
	"testing"
)

func TestHandler(t *testing.T) {
	h1 := skogul.Handler{}
	h2 := skogul.Handler{ Parser: parser.JSON{} }
	h3 := skogul.Handler{ Parser: parser.JSON{}, Transformers: []skogul.Transformer{}}
	h4 := skogul.Handler{ Parser: parser.JSON{}, Transformers: []skogul.Transformer{}, Sender: &(sender.Test{}) }
	h5 := skogul.Handler{ Parser: parser.JSON{}, Transformers: []skogul.Transformer{nil}, Sender: &(sender.Test{}) }

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
