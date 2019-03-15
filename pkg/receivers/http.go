/*
 * gollector, generic receiver
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

package receivers

import (
	"encoding/json"
	"fmt"
	. "github.com/KristianLyng/gollector/pkg/common"
	"io"
	"log"
	"net/http"
)

type HTTPReceiver struct {
	Handler *Handler
}

func (handler HTTPReceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		b := make([]byte, r.ContentLength)
		n, err := io.ReadFull(r.Body, b)
		if err != nil {
			log.Panicf("Read error from client %v, read %d bytes: %s", r.RemoteAddr, n, err)
		}
		var m GollectorContainer
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
		for _, s := range handler.Handler.Senders {
			s.Send(&m)
		}
		fmt.Fprintf(w, "OK\n")
	}
}

func (handler HTTPReceiver) Start() error {
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	return Gerror{"Shouldn't reach this"}
}
