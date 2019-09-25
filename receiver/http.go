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
	Address  string                        `doc:"Address to listen to"`
	Handlers map[string]*skogul.HandlerRef `doc:"Map of urls to handlers" example:"{\"/\": \"someHandler\" }"`
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

// defaultAddress is the address used if none is provided to the HTTP
// instance. It doesn't really make much sense to change it, since you
// wont be able to start multiple HTTP receivers on the same address
// anyway, so it's a const, not var. If you want to try: Just set the same
// Address on each HTTP receiver....
var defaultAddress = "[::1]:8080"

type httpReturn struct {
	Message string
}

func (rcvr receiver) answer(w http.ResponseWriter, r *http.Request, code int, inerr error) {
	answer := "OK"

	if inerr != nil {
		answer = inerr.Error()
	}

	b, err := json.Marshal(httpReturn{Message: answer})
	if err != nil {
		log.Panic("Failed to marshal internal JSON")
		return
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s\n", b)
}

func (rcvr receiver) handle(w http.ResponseWriter, r *http.Request) (code int, oerr error) {
	if r.ContentLength == 0 {
		oerr = skogul.Error{Source: "http receiver", Reason: "Missing input data"}
		code = 400
		return
	}
	b := make([]byte, r.ContentLength)
	n, err := io.ReadFull(r.Body, b)
	if err != nil {
		log.Printf("Read error from client %v, read %d bytes: %s", r.RemoteAddr, n, err)
		code = 400
		oerr = skogul.Error{Source: "http receiver", Reason: "read failed", Next: err}
		return
	}
	m, err := rcvr.Handler.Parser.Parse(b)
	if err == nil {
		err = m.Validate()
	}
	if err != nil {
		oerr = skogul.Error{Source: "http receiver", Reason: "failed to parse JSON", Next: err}
		code = 400
		return
	}
	for _, t := range rcvr.Handler.Transformers {
		t.Transform(&m)
	}
	err = rcvr.Handler.Sender.Send(&m)
	if err != nil {
		code = 500
		oerr = skogul.Error{Source: "http receiver", Reason: "failed to send data", Next: err}
		return
	}
	oerr = nil
	code = 200
	return
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
	if htt.Username != "" {
		log.Printf("Enforcing basic authentication for user `%s'", htt.Username)
		htt.auth = true
	} else {
		htt.auth = false
	}
	for idx, h := range htt.Handlers {
		log.Printf("Adding handler for %v", idx)
		http.Handle(idx, receiver{Handler: h.H, settings: htt})
	}
	if htt.Address == "" {
		log.Printf("HTTP: No listen-address specified. Using %s", defaultAddress)
		htt.Address = defaultAddress
	}
	if htt.Certfile != "" {
		log.Printf("Starting http receiver with TLS at %s", htt.Address)
		log.Fatal(http.ListenAndServeTLS(htt.Address, htt.Certfile, htt.Keyfile, nil))
	} else {
		log.Printf("Starting INSECURE http receiver (no TLS) at %s", htt.Address)
		log.Fatal(http.ListenAndServe(htt.Address, nil))
	}
	return skogul.Error{Reason: "Shouldn't reach this"}
}
