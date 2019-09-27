/*
 * skogul, log tests
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

import (
	"github.com/KristianLyng/skogul/config"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"log"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	config, err := config.Bytes([]byte(`
{
	"receivers": {
		"log": {
			"type": "log",
			"handler": "test_h"
		}
	},
	"handlers": {
		"test_h": {
			"parser": "json",
			"transformers": [],
			"sender": "dupe"
		}
	},
	"senders": {
		"dupe": {
			"type": "dupe",
			"next": ["test","print"]
		},
		"print": {
			"type": "debug"
		},
		"test": {
			"type": "test"
		}
	}
}`))

	if err != nil {
		t.Errorf("Failed to load config: %v", err)
		return
	}

	sTest := config.Senders["test"].Sender.(*sender.Test)
	rLog := config.Receivers["log"].Receiver.(*receiver.Log)
	go rLog.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
	log.Printf("This works!")
	time.Sleep(time.Duration(10 * time.Millisecond))
	got := sTest.Received()
	if got != 1 {
		t.Errorf("receiver.Tester{}, x.Start() failed to receive data. Expected some data, got 0.")
	}
}
