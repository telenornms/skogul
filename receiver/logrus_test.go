/*
 * skogul, logrus receiver
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
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

package receiver_test

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/telenornms/skogul/config"
	"github.com/telenornms/skogul/receiver"
	"github.com/telenornms/skogul/sender"
)

func TestLogrusLogReceivesData(t *testing.T) {

	config, err := config.Bytes([]byte(`
	{
		"receivers": {
			"logrus": {
				"type": "logrus",
				"handler": "test_h"
			}
		},
		"handlers": {
			"test_h": {
				"parser": "json",
				"transformers": [],
				"sender": "test"
			}
		},
		"senders": {
			"test": {
				"type": "test"
			}
		}
	}`))

	if err != nil {
		t.Errorf("Failed to load config (err: %s)", err)
		return
	}

	sTest := config.Senders["test"].Sender.(*sender.Test)
	rLog := config.Receivers["logrus"].Receiver.(*receiver.LogrusLog)
	go rLog.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
	logrus.Info("This works!")
	time.Sleep(time.Duration(10 * time.Millisecond))
	got := sTest.Received()
	if got != 1 {
		t.Errorf("receiver.Tester{}, x.Start() failed to receive data. Expected some data, got 0.")
	}
}
