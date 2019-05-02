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
	"net/url"

	skogul "github.com/KristianLyng/skogul/pkg"
)

/*
HTTP accepts HTTP connections on the Address specified.

Set it up similar to net/http:

        rcv := receiver.HTTP{Address: "localhost:8080"}
        rcv.Handle("/", foo)
        rcv.Handle("/blatti", bar)

*/
type HTTP struct {
	Address  string
	handlers map[string]*skogul.Handler
}

// For each path we handle, we set up a receiver such as this
// to simplify things.
// FIXME: This should almost certianly have a more descriptive name to
// avoid collisions and confusion.
type receiver struct {
	Handler *skogul.Handler
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

func (handler receiver) answer(w http.ResponseWriter, r *http.Request, code int, inerr error) {
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

func (handler receiver) handle(w http.ResponseWriter, r *http.Request) (oerr error, code int) {
	if r.ContentLength == 0 {
		oerr = skogul.Error{Source: "http receiver", Reason: "Missing input data"}
		code = 400
		return
	}
	b := make([]byte, r.ContentLength)
	n, err := io.ReadFull(r.Body, b)
	if err != nil {
		log.Printf("Read error from client %v, read %d bytes: %s", r.RemoteAddr, n, err)
		return skogul.Error{Source: "http receiver", Reason: "read failed", Next: err}, 400
	}
	m, err := handler.Handler.Parser.Parse(b)
	if err == nil {
		err = m.Validate()
	}
	if err != nil {
		return skogul.Error{Source: "http receiver", Reason: "failed to parse JSON", Next: err}, 400
	}
	for _, t := range handler.Handler.Transformers {
		t.Transform(&m)
	}
	err = handler.Handler.Sender.Send(&m)
	if err != nil {
		return skogul.Error{Source: "http receiver", Reason: "failed to send data", Next: err}, 500
	}
	return nil, 200
}

// Core HTTP handler
func (handler receiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, code := handler.handle(w, r)
	handler.answer(w, r, code, err)
}

/*
Handle adds a handler to a URL-pattern (same as net/http). Mostly
a convenience function to get less-ugly assignements.

Example:

        rcv := receiver.HTTP{Address: "localhost:8080"}
        rcv.Handle("/", foo)
        rcv.Handle("/blatti", bar)
*/
func (handler *HTTP) Handle(idx string, h *skogul.Handler) {
	if handler.handlers == nil {
		handler.handlers = make(map[string]*skogul.Handler)
	}
	if handler.handlers[idx] != nil {
		log.Fatalf("Error: Refusing to overwrite existing handler for %s", idx)
	}
	handler.handlers[idx] = h
}

// Start never returns.
func (handler *HTTP) Start() error {
	for idx, h := range handler.handlers {
		log.Printf("Adding handler for %v", idx)
		http.Handle(idx, receiver{h})
	}
	if handler.Address == "" {
		log.Printf("HTTP: No listen-address specified. Using %s", defaultAddress)
		handler.Address = defaultAddress
	}
	log.Printf("Starting http receiver at http://%s", handler.Address)
	log.Fatal(http.ListenAndServe(handler.Address, nil))
	return skogul.Error{Reason: "Shouldn't reach this"}
}

func init() {
	addAutoReceiver("http", NewHTTP, "Listen for Skogul-formatted JSON on a HTTP endpoint")
}

/*
NewHTTP returns a HTTP receiver, with the Path of the url being the one to
listen to.
*/
func NewHTTP(ul url.URL, h skogul.Handler) skogul.Receiver {
	hl := HTTP{Address: ul.Host}
	hl.Handle(fmt.Sprintf("/%s", ul.Path), &h)
	return &hl
}
