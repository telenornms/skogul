/*
 * skogul, SQL TLS utilities tests
 *
 * Copyright (c) 2026 Telenor Norge AS
 * Author(s):
 *  - Aslak Bakkeland <aslak.bakkeland@telenor.no>
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

package sql

import (
	"strings"
	"testing"
)

func TestHasMySQLTLSParam(t *testing.T) {
	tests := []struct {
		name     string
		connStr  string
		expected bool
	}{
		{"no params", "user:pass@tcp(host:3306)/db", false},
		{"other params only", "user:pass@tcp(host:3306)/db?parseTime=true", false},
		{"tls param", "user:pass@tcp(host:3306)/db?tls=custom", true},
		{"tls with other params", "user:pass@tcp(host:3306)/db?parseTime=true&tls=skip-verify", true},
		{"tls first param", "user:pass@tcp(host:3306)/db?tls=true&parseTime=true", true},
		{"tls in path not query", "user:pass@tcp(host:3306)/tls_database", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasMySQLTLSParam(tt.connStr)
			if result != tt.expected {
				t.Errorf("hasMySQLTLSParam(%q) = %v, want %v", tt.connStr, result, tt.expected)
			}
		})
	}
}

func TestHasPostgresSSLParam(t *testing.T) {
	tests := []struct {
		name     string
		connStr  string
		expected bool
	}{
		// DSN format tests
		{"dsn no ssl", "host=localhost user=postgres dbname=test", false},
		{"dsn sslmode", "host=localhost user=postgres sslmode=require", true},
		{"dsn sslmode at start", "sslmode=require host=localhost user=postgres", true},
		{"dsn sslcert", "host=localhost sslcert=/path/to/cert user=postgres", true},
		{"dsn sslkey", "host=localhost sslkey=/path/to/key user=postgres", true},
		{"dsn sslrootcert", "host=localhost sslrootcert=/path/to/ca user=postgres", true},

		// URL format tests
		{"url no ssl", "postgres://user:pass@localhost/db", false},
		{"url sslmode", "postgres://user:pass@localhost/db?sslmode=require", true},
		{"url sslmode with other", "postgres://user:pass@localhost/db?connect_timeout=10&sslmode=require", true},
		{"url sslcert", "postgresql://user:pass@localhost/db?sslcert=/path", true},

		// Edge cases
		{"ssl in hostname not param", "host=sslserver.example.com user=postgres", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPostgresSSLParam(tt.connStr)
			if result != tt.expected {
				t.Errorf("hasPostgresSSLParam(%q) = %v, want %v", tt.connStr, result, tt.expected)
			}
		})
	}
}

func TestIsPostgresURL(t *testing.T) {
	tests := []struct {
		connStr  string
		expected bool
	}{
		{"postgres://localhost/db", true},
		{"postgresql://localhost/db", true},
		{"host=localhost dbname=db", false},
		{"mysql://localhost/db", false},
	}

	for _, tt := range tests {
		t.Run(tt.connStr, func(t *testing.T) {
			result := isPostgresURL(tt.connStr)
			if result != tt.expected {
				t.Errorf("isPostgresURL(%q) = %v, want %v", tt.connStr, result, tt.expected)
			}
		})
	}
}

func TestAppendPostgresParam(t *testing.T) {
	tests := []struct {
		name     string
		connStr  string
		key      string
		value    string
		contains []string
	}{
		{
			"dsn format",
			"host=localhost user=postgres",
			"sslmode", "require",
			[]string{"host=localhost user=postgres ", "sslmode="},
		},
		{
			"url format no params",
			"postgres://localhost/db",
			"sslmode", "require",
			[]string{"postgres://localhost/db?sslmode=require"},
		},
		{
			"url format with params",
			"postgres://localhost/db?timeout=10",
			"sslmode", "require",
			[]string{"postgres://localhost/db?timeout=10&sslmode=require"},
		},
		{
			"dsn with special chars in value",
			"host=localhost",
			"sslcert", "/path/with spaces/cert.pem",
			[]string{"sslcert="},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appendPostgresParam(tt.connStr, tt.key, tt.value)
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("appendPostgresParam result %q missing expected substring %q", result, substr)
				}
			}
		})
	}
}

func TestSetupMySQLTLS(t *testing.T) {
	tests := []struct {
		name, connStr, contains string
		cfg                     TLSConfig
		wantErr, unchanged      bool
	}{
		{"no config", "user:pass@tcp(host)/db", "", TLSConfig{}, false, true},
		{"existing tls param", "user:pass@tcp(host)/db?tls=skip-verify", "", TLSConfig{Insecure: true}, false, true},
		{"insecure adds param", "user:pass@tcp(host)/db", "tls=skogul-tls-", TLSConfig{Insecure: true}, false, false},
		{"with existing params", "user:pass@tcp(host)/db?parseTime=true", "&tls=skogul-tls-", TLSConfig{Insecure: true}, false, false},
		{"nonexistent CA", "user:pass@tcp(host)/db", "", TLSConfig{CAFile: "/nonexistent/ca.pem"}, true, false},
		{"nonexistent cert", "user:pass@tcp(host)/db", "", TLSConfig{CertFile: "/x.pem", KeyFile: "/x.pem", Insecure: true}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SetupMySQLTLS(tt.connStr, &tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.unchanged && result != tt.connStr {
				t.Errorf("expected unchanged, got %q", result)
			}
			if tt.contains != "" && !strings.Contains(result, tt.contains) {
				t.Errorf("expected %q in result %q", tt.contains, result)
			}
		})
	}
}

func TestSetupPostgresTLS(t *testing.T) {
	tests := []struct {
		name, connStr string
		cfg           TLSConfig
		contains      []string
		unchanged     bool
	}{
		{"no config", "host=localhost", TLSConfig{}, nil, true},
		{"existing ssl", "host=localhost sslmode=disable", TLSConfig{Insecure: true}, nil, true},
		{"CA file", "host=localhost", TLSConfig{CAFile: "/ca.pem"}, []string{"sslmode=", "sslrootcert="}, false},
		{"insecure", "host=localhost", TLSConfig{Insecure: true}, []string{"sslmode="}, false},
		{"full mTLS", "host=localhost", TLSConfig{CAFile: "/ca.pem", CertFile: "/c.pem", KeyFile: "/k.pem"}, []string{"sslrootcert=", "sslcert=", "sslkey="}, false},
		{"URL format", "postgres://localhost/db", TLSConfig{Insecure: true}, []string{"?sslmode="}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SetupPostgresTLS(tt.connStr, &tt.cfg)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.unchanged && result != tt.connStr {
				t.Errorf("expected unchanged, got %q", result)
			}
			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("expected %q in result %q", s, result)
				}
			}
		})
	}
}

func TestTLSConfig_Verify(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *TLSConfig
		wantError bool
	}{
		{"empty config", &TLSConfig{}, false},
		{"insecure only", &TLSConfig{Insecure: true}, false},
		{"CA only", &TLSConfig{CAFile: "/ca.pem"}, false},
		{"cert without key", &TLSConfig{CertFile: "/cert.pem"}, true},
		{"key without cert", &TLSConfig{KeyFile: "/key.pem"}, true},
		{"cert+key without CA or insecure", &TLSConfig{CertFile: "/cert.pem", KeyFile: "/key.pem"}, true},
		{"cert+key with CA", &TLSConfig{CertFile: "/cert.pem", KeyFile: "/key.pem", CAFile: "/ca.pem"}, false},
		{"cert+key with insecure", &TLSConfig{CertFile: "/cert.pem", KeyFile: "/key.pem", Insecure: true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Verify()
			if (err != nil) != tt.wantError {
				t.Errorf("TLSConfig.Verify() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
