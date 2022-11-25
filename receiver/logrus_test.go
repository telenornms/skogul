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
	"os"
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

	sleepDuration := 1000 * time.Millisecond
	// Allow overriding the sleep duration for this package in e.g. CI environments
	// which might be more unpredictable with timings.
	if testSleepDuration := os.Getenv("SKOGUL_TEST_SLEEP_DURATION"); testSleepDuration != "" {
		d, err := time.ParseDuration(testSleepDuration)
		if err != nil {
			t.Errorf("Failed to parse custom sleep duration for logrus test: %s", err)
			return
		}
		sleepDuration = d
	}

	sTest := config.Senders["test"].Sender.(*sender.Test)
	rLog := config.Receivers["logrus"].Receiver.(*receiver.LogrusLog)
	go rLog.Start()
	time.Sleep(sleepDuration)
	logrus.Info("This works!")
	time.Sleep(sleepDuration)
	got := sTest.Received()
	if got != 1 {
		t.Errorf("receiver.Tester{}, x.Start() failed to receive data. Expected some data, got 0.")
	}
}
