/*
 * skogul, SQL TLS utilities
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

// Package sql provides shared utilities for SQL sender and receiver TLS
// configuration.
package sql

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/telenornms/skogul"
)

var (
	log        = skogul.Logger("internal", "sql")
	tlsCounter atomic.Uint64
)

type TLSConfig struct {
	CAFile   string
	CertFile string
	KeyFile  string
	Insecure bool

	validated  bool
	rootCAs    *x509.CertPool
	clientCert *tls.Certificate
}

func hasMySQLTLSParam(connStr string) bool {
	queryStart := strings.Index(connStr, "?")
	if queryStart == -1 {
		return false
	}
	queryPart := connStr[queryStart:]
	return strings.Contains(queryPart, "tls=")
}

func hasPostgresSSLParam(connStr string) bool {
	for _, p := range []string{"sslmode=", "sslcert=", "sslkey=", "sslrootcert="} {
		if strings.HasPrefix(connStr, p) || strings.Contains(connStr, " "+p) ||
			strings.Contains(connStr, "?"+p) || strings.Contains(connStr, "&"+p) {
			return true
		}
	}
	return false
}

func isPostgresURL(connStr string) bool {
	return strings.HasPrefix(connStr, "postgres://") || strings.HasPrefix(connStr, "postgresql://")
}

func appendPostgresParam(connStr, key, value string) string {
	if isPostgresURL(connStr) {
		sep := "?"
		if strings.Contains(connStr, "?") {
			sep = "&"
		}
		return connStr + sep + key + "=" + url.QueryEscape(value)
	}
	return connStr + " " + key + "=" + pq.QuoteLiteral(value)
}

// SetupMySQLTLS configures TLS for MySQL connections.
func SetupMySQLTLS(connStr string, cfg *TLSConfig) (string, error) {
	if cfg.CAFile == "" && cfg.CertFile == "" && cfg.KeyFile == "" && !cfg.Insecure {
		return connStr, nil
	}

	if hasMySQLTLSParam(connStr) {
		log.Warn("TLS parameter already present in connection string, skipping TLS configuration from CAFile/CertFile/KeyFile/Insecure fields")
		return connStr, nil
	}

	if err := cfg.Verify(); err != nil {
		return connStr, err
	}
	if err := cfg.loadCerts(); err != nil {
		return connStr, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.Insecure,
		RootCAs:            cfg.rootCAs,
	}

	if cfg.Insecure {
		log.Warn("TLS certificate verification disabled - this is insecure")
	}

	if cfg.clientCert != nil {
		tlsConfig.Certificates = []tls.Certificate{*cfg.clientCert}
	}

	configName := fmt.Sprintf("skogul-tls-%d", tlsCounter.Add(1))
	if err := mysql.RegisterTLSConfig(configName, tlsConfig); err != nil {
		return connStr, fmt.Errorf("failed to register MySQL TLS config: %w", err)
	}

	sep := "?"
	if strings.Contains(connStr, "?") {
		sep = "&"
	}
	return connStr + sep + "tls=" + configName, nil
}

func SetupPostgresTLS(connStr string, cfg *TLSConfig) (string, error) {
	if cfg.CAFile == "" && cfg.CertFile == "" && cfg.KeyFile == "" && !cfg.Insecure {
		return connStr, nil
	}

	if hasPostgresSSLParam(connStr) {
		log.Warn("SSL parameters already present in connection string, skipping TLS configuration from CAFile/CertFile/KeyFile/Insecure fields")
		return connStr, nil
	}

	// Set sslmode based on configuration. Order matters: most secure first.
	switch {
	case cfg.CAFile != "":
		connStr = appendPostgresParam(connStr, "sslmode", "verify-full")
	case cfg.Insecure:
		log.Warn("TLS certificate verification disabled - this is insecure")
		connStr = appendPostgresParam(connStr, "sslmode", "require")
	}

	if cfg.CAFile != "" {
		connStr = appendPostgresParam(connStr, "sslrootcert", cfg.CAFile)
	}
	if cfg.CertFile != "" {
		connStr = appendPostgresParam(connStr, "sslcert", cfg.CertFile)
	}
	if cfg.KeyFile != "" {
		connStr = appendPostgresParam(connStr, "sslkey", cfg.KeyFile)
	}

	return connStr, nil
}

func (cfg *TLSConfig) Verify() error {
	if (cfg.CertFile != "" && cfg.KeyFile == "") || (cfg.CertFile == "" && cfg.KeyFile != "") {
		return fmt.Errorf("must provide both CertFile and KeyFile, or neither")
	}
	if cfg.CertFile != "" && cfg.CAFile == "" && !cfg.Insecure {
		return fmt.Errorf("client certificate requires either CAFile for server verification or Insecure: true")
	}
	return nil
}

func (cfg *TLSConfig) loadCerts() error {
	if cfg.validated {
		return nil
	}

	if cfg.CAFile != "" {
		cp, err := skogul.GetCertPool(cfg.CAFile)
		if err != nil {
			return fmt.Errorf("failed to load CA file: %w", err)
		}
		cfg.rootCAs = cp
	}

	if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load client certificate: %w", err)
		}
		cfg.clientCert = &cert
	}

	cfg.validated = true
	return nil
}
