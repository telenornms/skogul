/*
 * skogul, generic receiver
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
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

package receiver

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/telenornms/skogul"
)

var httpLog = skogul.Logger("receiver", "http")

// HTTPAuth contains ways to authenticate a HTTP request, e.g. Username/Password for Basic Auth.
type HTTPAuth struct {
	Username              string `doc:"Username for basic authentication. No authentication is required if left blank."`
	Password              string `doc:"Password for basic authentication."`
	SANDNSName            string `doc:"DNS name which has to be present in SAN extension of x509 certificate when using Client Certificate authentication"`
	SkipCertificateVerify bool   `doc:"Skip verifying certificate. (default: false)"`
	path                  string
}

/*
HTTP accepts HTTP connections on the Address specified, and requires at
least one handler to be set up, using Handle. This is done implicitly
if the HTTP receiver is created using New()
*/
type HTTP struct {
	Address              string                        `doc:"Address to listen to." example:"[::1]:80 [2001:db8::1]:443"`
	Handlers             map[string]*skogul.HandlerRef `doc:"Paths to handlers. Need at least one." example:"{\"/\": \"someHandler\" }"`
	Auth                 map[string]*HTTPAuth          `doc:"A map corresponding to Handlers; specifying authentication for the given path, if required."`
	Certfile             string                        `doc:"Path to certificate file for TLS. If left blank, un-encrypted HTTP is used."`
	Keyfile              string                        `doc:"Path to key file for TLS."`
	ClientCertificateCAs []string                      `doc:"Paths to files containing CAs which are accepted for Client Certificate authentication."`
}

// For each path we handle, we set up a receiver such as this
// to simplify things.
// FIXME: This should almost certianly have a more descriptive name to
// avoid collisions and confusion.
type receiver struct {
	Handler  *skogul.Handler
	settings *HTTP
	auth     *HTTPAuth
}

type httpReturn struct {
	Message string
}

func (auth *HTTPAuth) auth(r *http.Request) error {
	if auth.Username != "" && auth.Password != "" {
		username, pw, ok := r.BasicAuth()
		success := ok && auth.Username == username && auth.Password == pw
		if !success {
			return skogul.Error{Source: "http receiver", Reason: "Invalid credentials"}
		}

		return nil
	}

	if auth.SANDNSName != "" {
		httpLog.Trace("Verifying request using client certificates")
		if err := auth.verifyPeerCertificate(nil, r.TLS.VerifiedChains); err != nil {
			return err
		} else {
			return nil
		}
	}

	return skogul.Error{Source: "http receiver", Reason: "No matching authentication method"}
}

func (rcvr receiver) answer(w http.ResponseWriter, r *http.Request, code int, inerr error) {
	answer := "OK"

	w.WriteHeader(code)
	if code == 204 {
		return
	}
	if inerr != nil {
		answer = inerr.Error()
	}

	b, err := json.Marshal(httpReturn{Message: answer})
	skogul.Assert(err == nil, err)
	fmt.Fprintf(w, "%s\n", b)
}

func (rcvr receiver) handle(w http.ResponseWriter, r *http.Request) (int, error) {
	if rcvr.auth != nil {
		if err := rcvr.auth.auth(r); err != nil {
			return 401, err
		}
	}

	if r.ContentLength == 0 {
		return 400, skogul.Error{Source: "http receiver", Reason: "Missing input data"}
	}

	b := make([]byte, r.ContentLength)

	if n, err := io.ReadFull(r.Body, b); err != nil {
		httpLog.WithFields(log.Fields{
			"address":  r.RemoteAddr,
			"error":    err,
			"numbytes": n,
		}).Error("Read error from client")
		return 400, skogul.Error{Source: "http receiver", Reason: "read failed", Next: err}
	}

	if err := rcvr.Handler.Handle(b); err != nil {
		return 400, err
	}

	return 204, nil
}

// Core HTTP handler
func (rcvr receiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code, err := rcvr.handle(w, r)
	rcvr.answer(w, r, code, err)
}

// loadClientCertificateCAs loads a given list of strings as paths of
// acceptable CAs for use in Client Certificate authentication.
func loadClientCertificateCAs(paths []string) (*x509.CertPool, error) {
	httpLog.Debugf("Loading Client Certificates from %d file(s)", len(paths))
	pool := x509.NewCertPool()
	for _, path := range paths {
		if data, err := ioutil.ReadFile(path); err != nil {
			httpLog.WithError(err).WithField("path", path).Error("Failed to read certificate file")
			return nil, err
		} else {
			pool.AppendCertsFromPEM(data)
		}
	}
	return pool, nil
}

// Start never returns.
func (htt *HTTP) Start() error {
	server := http.Server{}
	serveMux := http.NewServeMux()
	server.Handler = serveMux
	for idx, h := range htt.Handlers {
		httpLog.WithFields(log.Fields{
			"configuredHandler": idx,
			"selectedHandler":   h.Name,
			"hasAuth":           htt.Auth[idx] != nil,
		}).Debug("Adding handler")

		serveMux.Handle(idx, receiver{Handler: h.H, settings: htt, auth: htt.Auth[idx]})
	}

	if len(htt.ClientCertificateCAs) > 0 {
		pool, err := loadClientCertificateCAs(htt.ClientCertificateCAs)
		if err != nil {
			httpLog.WithError(err).Error("Failed to load Client Certificates")
			return err
		}
		server.TLSConfig = &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.VerifyClientCertIfGiven,
		}
		httpLog.Info("Configured HTTP receiver with Client Certificate authentication")
		for _, auth := range htt.Auth {
			if auth.SANDNSName != "" {
				httpLog.Info("Configured HTTP receiver with Client Certificate verification")
				break
			}
		}
	}

	server.Addr = htt.Address
	if htt.Certfile != "" {
		httpLog.WithField("address", htt.Address).Info("Starting http receiver with TLS")
		httpLog.Fatal(server.ListenAndServeTLS(htt.Certfile, htt.Keyfile))
	} else {
		httpLog.WithField("address", htt.Address).Info("Starting INSECURE http receiver (no TLS)")
		httpLog.Fatal(server.ListenAndServe())
	}
	return skogul.Error{Reason: "Shouldn't reach this"}
}

// verifyPeerCertificate verifies a client certificate presented to us
// during TLS handshake by comparing its extensions (such as SAN) to
// some expected value(s)
func (auth *HTTPAuth) verifyPeerCertificate(_ [][]byte, verifiedChains [][]*x509.Certificate) error {
	if auth.SkipCertificateVerify || auth.SANDNSName == "" {
		httpLog.WithFields(log.Fields{"skip": auth.SkipCertificateVerify, "dns_name": auth.SANDNSName}).Trace("Skipping verifying certificate")
		return nil
	}

	certLogger := httpLog
	certLogger.Tracef("Verifying %d certificate chains", len(verifiedChains))
	for _, chain := range verifiedChains {
		if len(chain) < 1 {
			// Chain has no certificates (?)
			continue
		}
		cert := chain[0]
		certDebugLogger := certLogger.WithFields(log.Fields{
			"issuer":           cert.Issuer,
			"subject":          cert.Subject,
			"num_dns_names":    len(cert.DNSNames),
			"num_emails":       len(cert.EmailAddresses),
			"num_ip_addresses": len(cert.IPAddresses),
			"num_uris":         len(cert.URIs),
			"num_x509_ext":     len(cert.Extensions),
		})
		certDebugLogger.Trace("Verifying certificate")
		for _, dnsName := range cert.DNSNames {
			if auth.SANDNSName != "" && strings.ToLower(dnsName) == strings.ToLower(auth.SANDNSName) {
				// If we find a matching DNS name in the SANs, we return non-error
				// which specifies that we're done verifying with access granted.
				return nil
			}
		}
	}
	// If no checks until now have returned a non-error,
	// we return an error to tell the verifying function that this certificate
	// is not verified, and access is denied.
	// This will present the user with a 'bad certificate' alert.
	return skogul.Error{Reason: "Failed to verify x509 SAN DNS Name", Source: "http-receiver"}
}

// Verify verifies the configuration for the HTTP receiver
func (htt *HTTP) Verify() error {
	if htt.Handlers == nil || len(htt.Handlers) == 0 {
		httpLog.Error("No handlers specified. Need at least one.")
		return skogul.Error{Source: "http receiver", Reason: "No handlers specified. Need at least one."}
	}

	if htt.Address == "" {
		httpLog.Warn("Missing listen address for http receiver, using Go default")
	}

	if htt.Certfile == "" && htt.Auth != nil {
		httpLog.Warn("HTTP receiver configured with authentication but not with TLS! Auth will happen in the open!")
	}
	if (htt.Certfile != "" && htt.Keyfile == "") || (htt.Certfile == "" && htt.Keyfile != "") {
		return skogul.Error{Source: "http-receiver", Reason: "Specify both Certfile and Keyfile if either is specified."}
	}
	cas, err := loadClientCertificateCAs(htt.ClientCertificateCAs)
	if err != nil {
		return skogul.Error{Source: "http-receiver", Reason: "Failed to load Client Certificates CAs", Next: err}
	}
	for _, auth := range htt.Auth {
		if auth.SANDNSName != "" && cas == nil {
			return skogul.Error{Source: "http-receiver", Reason: "No Client Certificate CAs defined, but DNS Name for SAN specified. Specify ClientCertificateCAs configuration element."}
		}
	}

	return nil
}
