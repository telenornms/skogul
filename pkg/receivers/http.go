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

/*
Receivers accept data and execute a handler. They are the "inbound"
API of Skogul.
*/
package receivers

import (
	"encoding/json"
	"fmt"
	"github.com/KristianLyng/skogul/pkg"
	"io"
	"log"
	"net/http"
)

/*
The HTTP receiver accepts HTTP connections on the Address
specified and directs valid Skogul metric containers to
the appropriate skogul.Handler.

Set it up similar to net/http:

        rcv := receiver.HTTP{Address: "localhost:8080"}
        rcv.Handle("/", foo)
        rcv.Handle("/blatti", bar)

*/

type HTTP struct {
	Address  string
	handlers map[string]*skogul.Handler
}

type receiver struct {
	Handler *skogul.Handler
}

func (handler receiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		b := make([]byte, r.ContentLength)
		n, err := io.ReadFull(r.Body, b)
		if err != nil {
			log.Printf("Read error from client %v, read %d bytes: %s", r.RemoteAddr, n, err)
		}
		var m skogul.Container
		err = json.Unmarshal(b, &m)
		if err == nil {
			err = m.Validate()
		}
		if err != nil {
			fmt.Fprintf(w, "Unable to parse JSON: %s", err)
		}
		for _, t := range handler.Handler.Transformers {
			t.Transform(&m)
		}
		handler.Handler.Sender.Send(&m)
		fmt.Fprintf(w, "OK\n")
	}
}

/*
Handle adds a handler to a URL-pattern (same as net/http). Mostly
a convenience function to get less-ugly assignements.
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

// Start the HTTP receiver
func (handler *HTTP) Start() error {
	for idx, h := range handler.handlers {
		log.Printf("Adding handler for %v", idx)
		http.Handle(idx, receiver{h})
	}
	if handler.Address == "" {
		log.Print("HTTP: No listen-address specified. Using localhost:8080")
		handler.Address = "localhost:8080"
	}
	log.Printf("Starting http receiver at http://%s", handler.Address)
	log.Fatal(http.ListenAndServe(handler.Address, nil))
	return skogul.Gerror{"Shouldn't reach this"}
}
