/*
 * skogul, generic receiver
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

package receiver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/telenornms/skogul"
)

var httpLog = skogul.Logger("receiver", "http")

// HTTPAuth contains ways to authenticate a HTTP request, e.g. Username/Password for Basic Auth.
type HTTPAuth struct {
	Username string `doc:"Username for basic authentication. No authentication is required if left blank."`
	Password string `doc:"Password for basic authentication."`
}

/*
HTTP accepts HTTP connections on the Address specified, and requires at
least one handler to be set up, using Handle. This is done implicitly
if the HTTP receiver is created using New()
*/
type HTTP struct {
	Address  string                        `doc:"Address to listen to." example:"[::1]:80 [2001:db8::1]:443"`
	Handlers map[string]*skogul.HandlerRef `doc:"Paths to handlers. Need at least one." example:"{\"/\": \"someHandler\" }"`
	Auth     map[string]*HTTPAuth          `doc:"A map corresponding to Handlers; specifying authentication for the given path, if required."`
	Certfile string                        `doc:"Path to certificate file for TLS. If left blank, un-encrypted HTTP is used."`
	Keyfile  string                        `doc:"Path to key file for TLS."`
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

func (auth *HTTPAuth) auth(r *http.Request) (bool, error) {
	var authErrMsg skogul.Error

	if auth.Username != "" && auth.Password != "" {
		username, pw, ok := r.BasicAuth()
		success := ok && auth.Username == username && auth.Password == pw
		if !success {
			authErrMsg = skogul.Error{Source: "http receiver", Reason: "Invalid credentials"}
		}

		return success, authErrMsg
	}

	return false, nil
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
		ok, err := rcvr.auth.auth(r)
		if !ok && err == nil {
			return 401, skogul.Error{Source: "http receiver", Reason: "Auth error"}
		}

		if !ok {
			// !ok could be either "yeah no this didn't authenticate" or "yeah no this doesn't have access"...
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

	return nil
}
