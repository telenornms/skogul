/*
* skogul, SQL receiver tests
* Copyright (c) 2026 Telenor Norge AS
* Author(s):
* - Aslak Bakkeland <aslak.bakkeland@telenor.no>
* This library is free software; you can redistribute it and/or modify it under
* the terms of the GNU Lesser General Public License as published by the Free
* Software Foundation; either version 2.1 of the License, or (at your option)
* any later version.
* This library is distributed in the hope that it will be useful, but WITHOUT
* ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
* FOR A PARTICULAR PURPOSE.  See the GNU Lesser General Public License for more
* details.
* You should have received a copy of the GNU Lesser General Public License
* along with this library; if not, write to the Free Software Foundation, Inc.,
* 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301  USA */

package receiver_test

import (
	"fmt"
	"testing"

	"github.com/telenornms/skogul/config"
)

var sqlReceiverBase = `
{
	"receivers": {
		"sql": {
			"type": "sql",
			"handler": "test",
			%s
		}
	},
	"handlers": {
		"test": {
			"parser": "json",
			"sender": "test"
		}
	},
	"senders": {
		"test": {
			"type": "test"
		}
	}
}`

func sqlReceiverTestAuto(t *testing.T, params string) *config.Config {
	t.Helper()
	conf, err := config.Bytes(fmt.Appendf(nil, sqlReceiverBase, params))
	if conf == nil {
		t.Errorf("Bytes(\"%s\" failed", params)
	}
	if err != nil {
		t.Errorf("Bytes(\"%s\" failed: %v", params, err)
	}
	return conf
}

func sqlReceiverTestAutoNeg(t *testing.T, params string) {
	t.Helper()
	conf, err := config.Bytes(fmt.Appendf(nil, sqlReceiverBase, params))
	if conf != nil {
		t.Errorf("Bytes(\"%s\" succeeded, but expected failure. Val: %v", params, conf)
	}
	if err == nil {
		t.Errorf("Bytes(\"%s\" succeeded, but expected failure. Val: %v", params, conf)
	}
}

func TestSQLReceiver_auto(t *testing.T) {
	// Missing required fields
	sqlReceiverTestAutoNeg(t, `"driver":"mysql"`)
	sqlReceiverTestAutoNeg(t, `"driver":"mysql","connstr": "something"`)
	sqlReceiverTestAutoNeg(t, `"driver":"mysql","query": "something"`)

	// Valid configs
	sqlReceiverTestAuto(t, `"driver":"mysql","connstr":"something","query": "SELECT * FROM test"`)
	sqlReceiverTestAuto(t, `"driver":"postgres","connstr":"something","query": "SELECT * FROM test"`)

	// Invalid driver
	sqlReceiverTestAutoNeg(t, `"driver":"invalid","connstr":"something","query": "SELECT * FROM test"`)
}

func TestSQLReceiver_tls_config(t *testing.T) {
	// TLS validation is tested in detail in internal/sql/tls_test.go.
	// Here we verify the receiver integrates with that validation.

	// Invalid: cert without key
	sqlReceiverTestAutoNeg(t, `"driver":"mysql","connstr":"x","query":"SELECT 1","certfile":"/cert.pem"`)

	// Invalid: cert+key without CA or insecure
	sqlReceiverTestAutoNeg(t, `"driver":"mysql","connstr":"x","query":"SELECT 1","certfile":"/cert.pem","keyfile":"/key.pem"`)

	// Valid: insecure flag alone
	sqlReceiverTestAuto(t, `"driver":"mysql","connstr":"x","query":"SELECT 1","insecure":true`)

	// Valid: full mTLS config (file existence checked at runtime)
	sqlReceiverTestAuto(t, `"driver":"mysql","connstr":"x","query":"SELECT 1","cafile":"/ca.pem","certfile":"/cert.pem","keyfile":"/key.pem"`)
}
