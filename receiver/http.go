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
	"log"
	"net/http"

	"github.com/KristianLyng/skogul"
)

/*
HTTP accepts HTTP connections on the Address specified, and requires at
least one handler to be set up, using Handle. This is done implicitly
if the HTTP receiver is created using New()
*/
type HTTP struct {
	Address  string                        `doc:"Address to listen to." example:"[::1]:80 [2001:db8::1]:443"`
	Handlers map[string]*skogul.HandlerRef `doc:"Paths to handlers. Need at least one." example:"{\"/\": \"someHandler\" }"`
	Username string                        `doc:"Username for basic authentication. No authentication is required if left blank."`
	Password string                        `doc:"Password for basic authentication."`
	Certfile string                        `doc:"Path to certificate file for TLS. If left blank, un-encrypted HTTP is used."`
	Keyfile  string                        `doc:"Path to key file for TLS."`
	auth     bool
}

// For each path we handle, we set up a receiver such as this
// to simplify things.
// FIXME: This should almost certianly have a more descriptive name to
// avoid collisions and confusion.
type receiver struct {
	Handler  *skogul.Handler
	settings *HTTP
}

type httpReturn struct {
	Message string
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
	if r.ContentLength == 0 {
		return 400, skogul.Error{Source: "http receiver", Reason: "Missing input data"}
	}

	b := make([]byte, r.ContentLength)

	if n, err := io.ReadFull(r.Body, b); err != nil {
		log.Printf("Read error from client %v, read %d bytes: %s", r.RemoteAddr, n, err)
		return 400, skogul.Error{Source: "http receiver", Reason: "read failed", Next: err}
	}

	if err := rcvr.Handler.Handle(b); err != nil {
		return 400, err
	}

	return 204, nil
}

// Core HTTP handler
func (rcvr receiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rcvr.settings.auth {
		user, pass, ok := r.BasicAuth()
		if !ok || user != rcvr.settings.Username || pass != rcvr.settings.Password {
			rcvr.answer(w, r, 401, skogul.Error{Source: "http receiver", Reason: "Authentication failed"})
			return
		}

	}
	code, err := rcvr.handle(w, r)
	rcvr.answer(w, r, code, err)
}

// Start never returns.
func (htt *HTTP) Start() error {
	server := http.Server{}
	serveMux := http.NewServeMux()
	server.Handler = serveMux
	if htt.Username != "" {
		log.Printf("Enforcing basic authentication for user `%s'", htt.Username)
		if htt.Password == "" {
			log.Fatal("HTTP receiver has a Username provided, but not a password? Probably a mistake.")
		}
		htt.auth = true
	} else {
		if htt.Password != "" {
			log.Fatal("Password provided for HTTP receiver, but not a username? Probably a mistake.")
		}
		htt.auth = false
	}
	for idx, h := range htt.Handlers {
		log.Printf("Adding handler %v -> %v", idx, h.Name)
		serveMux.Handle(idx, receiver{Handler: h.H, settings: htt})
	}
	if htt.Address == "" {
		log.Printf("HTTP: No listen-address specified. Using go default (probably :http or :https?)")
	}
	server.Addr = htt.Address
	if htt.Certfile != "" {
		log.Printf("Starting http receiver with TLS at %s", htt.Address)
		log.Fatal(server.ListenAndServeTLS(htt.Certfile, htt.Keyfile))
	} else {
		log.Printf("Starting INSECURE http receiver (no TLS) at %s", htt.Address)
		log.Fatal(server.ListenAndServe())
	}
	return skogul.Error{Reason: "Shouldn't reach this"}
}

// Verify ensures at least one handler exists for the HTTP receiver.
func (htt *HTTP) Verify() error {
	if htt.Handlers == nil || len(htt.Handlers) == 0 {
		return skogul.Error{Source: "http receiver", Reason: "No handlers specified. Need at least one."}
	}
	return nil
}
