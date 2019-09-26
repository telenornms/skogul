/*
 * skogul, complex receiver tests
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
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/config"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"testing"
	"time"
)

var validContainer = skogul.Container{}

func init() {

	now := time.Now()
	m := skogul.Metric{}
	m.Time = &now
	m.Metadata = make(map[string]interface{})
	m.Data = make(map[string]interface{})
	m.Metadata["foo"] = "bar"
	m.Data["tall"] = 5
	validContainer.Metrics = []*skogul.Metric{&m}
}

// Tests http receiver, sender and JSON parser implicitly
func TestHttp_stack(t *testing.T) {
	config, err := config.Bytes([]byte(`
{
	"senders": {
		"plain_origin": {
			"type": "http",
			"url": "http://[::1]:1339"
		},
		"auth_plain_origin": {
			"type": "http",
			"url": "http://god:hunter2@[::1]:3000"
		},
		"auth_plain_fail1": {
			"type": "http",
			"url": "http://gad:hunter2@[::1]:3000"
		},
		"auth_plain_fail2": {
			"type": "http",
			"url": "http://gad:@[::1]:3000"
		},
		"auth_plain_fail3": {
			"type": "http",
			"url": "http://:hunter2@[::1]:3000"
		},
		"auth_plain_fail4": {
			"type": "http",
			"url": "http://[::1]:3000"
		},
		"ssl_ok": {
			"type": "http",
			"url": "https://[::1]:3443",
			"insecure": true
		},
		"ssl_bad1": {
			"type": "http",
			"url": "https://[::1]:3443",
			"insecure": false
		},
		"ssl_bad2": {
			"type": "http",
			"url": "http://[::1]:3443"
		},
		"ssl_auth_ok": {
			"type": "http",
			"url": "https://god:hunter2@[::1]:5443",
			"insecure": true
		},
		"ssl_auth_bad1": {
			"type": "http",
			"url": "https://gad:hunter2@[::1]:5443",
			"insecure": true
		},
		"ssl_auth_bad2": {
			"type": "http",
			"url": "https://god:hunter3@[::1]:5443",
			"insecure": true
		},
		"ssl_auth_bad3": {
			"type": "http",
			"url": "https://god:@[::1]:5443",
			"insecure": true
		},
		"ssl_auth_bad4": {
			"type": "http",
			"url": "https://:hunter2@[::1]:5443",
			"insecure": true
		},
		"ssl_auth_bad5": {
			"type": "http",
			"url": "http://god:hunter2@[::1]:5443"
		},
		"ssl_auth_bad6": {
			"type": "http",
			"url": "https://god:hunter2@[::1]:5443",
			"insecure": false
		},
		"common": {
			"type": "test"
		}
	},
	"receivers": {
		"plain": {
			"type": "http",
			"address": "[::1]:1339",
			"handlers": { "/": "common"}
		},
		"auth": {
			"type": "http",
			"address": "[::1]:3000",
			"handlers": { "/": "common"},
			"username": "god",
			"password": "hunter2"
		},
		"ssl_noauth": {
			"type": "http",
			"address": "[::1]:3443",
			"handlers": { "/": "common"},
			"certfile": "../examples/cacert-snakeoil.pem",
			"keyfile": "../examples/privkey-snakeoil.pem"
		},
		"ssl_auth": {
			"type": "http",
			"address": "[::1]:5443",
			"handlers": { "/": "common"},
			"certfile": "../examples/cacert-snakeoil.pem",
			"keyfile": "../examples/privkey-snakeoil.pem",
			"username": "god",
			"password": "hunter2"
		}
	},
	"handlers": {
		"common": {
			"parser": "json",
			"transformers": [],
			"sender": "common"
		}
	}
}`))

	if err != nil {
		t.Errorf("Failed to load config: %v", err)
		return
	}

	sPlainOrigin := config.Senders["plain_origin"].Sender.(*sender.HTTP)
	sAuthPlainOrigin := config.Senders["auth_plain_origin"].Sender.(*sender.HTTP)
	sAuthPlainFail1 := config.Senders["auth_plain_fail1"].Sender.(*sender.HTTP)
	sAuthPlainFail2 := config.Senders["auth_plain_fail2"].Sender.(*sender.HTTP)
	sAuthPlainFail3 := config.Senders["auth_plain_fail3"].Sender.(*sender.HTTP)
	sAuthPlainFail4 := config.Senders["auth_plain_fail4"].Sender.(*sender.HTTP)
	sSSLOk1 := config.Senders["ssl_ok"].Sender.(*sender.HTTP)
	sSSLBad1 := config.Senders["ssl_bad1"].Sender.(*sender.HTTP)
	sSSLBad2 := config.Senders["ssl_bad2"].Sender.(*sender.HTTP)
	sSSLAuthOk1 := config.Senders["ssl_auth_ok"].Sender.(*sender.HTTP)
	sSSLAuthBad1 := config.Senders["ssl_auth_bad1"].Sender.(*sender.HTTP)
	sSSLAuthBad2 := config.Senders["ssl_auth_bad2"].Sender.(*sender.HTTP)
	sSSLAuthBad3 := config.Senders["ssl_auth_bad3"].Sender.(*sender.HTTP)
	sSSLAuthBad4 := config.Senders["ssl_auth_bad4"].Sender.(*sender.HTTP)
	sSSLAuthBad5 := config.Senders["ssl_auth_bad5"].Sender.(*sender.HTTP)
	sSSLAuthBad6 := config.Senders["ssl_auth_bad6"].Sender.(*sender.HTTP)
	sCommon := config.Senders["common"].Sender.(*sender.Test)

	rPlain := config.Receivers["plain"].Receiver.(*receiver.HTTP)
	rAuth := config.Receivers["auth"].Receiver.(*receiver.HTTP)
	rSSLNoAuth := config.Receivers["ssl_noauth"].Receiver.(*receiver.HTTP)
	rSSLAuth := config.Receivers["ssl_auth"].Receiver.(*receiver.HTTP)

	go rPlain.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
	sCommon.TestQuick(t, sPlainOrigin, &validContainer, 1)
	blank := skogul.Container{}
	sCommon.TestNegative(t, sPlainOrigin, &blank)

	go rAuth.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
	sCommon.TestQuick(t, sAuthPlainOrigin, &validContainer, 1)
	sCommon.TestNegative(t, sAuthPlainOrigin, &blank)
	sCommon.TestNegative(t, sAuthPlainFail1, &validContainer)
	sCommon.TestNegative(t, sAuthPlainFail2, &validContainer)
	sCommon.TestNegative(t, sAuthPlainFail3, &validContainer)
	sCommon.TestNegative(t, sAuthPlainFail4, &validContainer)
	sCommon.TestQuick(t, sAuthPlainOrigin, &validContainer, 1)

	go rSSLNoAuth.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
	sCommon.TestQuick(t, sSSLOk1, &validContainer, 1)
	sCommon.TestNegative(t, sSSLBad1, &validContainer)
	sCommon.TestNegative(t, sSSLBad2, &validContainer)
	sCommon.TestQuick(t, sSSLOk1, &validContainer, 1)

	go rSSLAuth.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
	sCommon.TestQuick(t, sSSLAuthOk1, &validContainer, 1)
	sCommon.TestNegative(t, sSSLAuthBad1, &validContainer)
	sCommon.TestQuick(t, sSSLAuthOk1, &validContainer, 1)
	sCommon.TestNegative(t, sSSLAuthBad2, &validContainer)
	sCommon.TestQuick(t, sSSLAuthOk1, &validContainer, 1)
	sCommon.TestNegative(t, sSSLAuthBad3, &validContainer)
	sCommon.TestQuick(t, sSSLAuthOk1, &validContainer, 1)
	sCommon.TestNegative(t, sSSLAuthBad4, &validContainer)
	sCommon.TestQuick(t, sSSLAuthOk1, &validContainer, 1)
	sCommon.TestNegative(t, sSSLAuthBad5, &validContainer)
	sCommon.TestQuick(t, sSSLAuthOk1, &validContainer, 1)
	sCommon.TestNegative(t, sSSLAuthBad6, &validContainer)
	sCommon.TestQuick(t, sSSLAuthOk1, &validContainer, 1)

}
