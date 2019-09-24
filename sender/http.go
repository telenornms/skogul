/*
 * skogul, http writer
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

package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/KristianLyng/skogul"
)

/*
HTTP sender POSTs the Skogul JSON-encoded data to the provided URL.
*/
type HTTP struct {
	URL     string          `doc:"Fully qualified URL to send data to" example:"http://localhost:6081/"`
	Timeout skogul.Duration `doc:"Timeout on POST"`
	once    sync.Once
	client  *http.Client
}

// Send POSTS data
func (ht *HTTP) Send(c *skogul.Container) error {
	b, err := json.Marshal(*c)
	if err != nil {
		e := skogul.Error{Source: "http sender", Reason: "unable to marshal JSON", Next: err}
		log.Print(e)
		return e
	}
	var buffer bytes.Buffer
	buffer.Write(b)
	req, err := http.NewRequest("POST", ht.URL, &buffer)
	if err != nil {
		e := skogul.Error{Source: "http sender", Reason: "unable to create new HTTP request", Next: err}
		log.Print(e)
		return e
	}
	req.Header.Set("Content-Type", "application/json")
	ht.once.Do(func() {
		if ht.Timeout.Duration == 0 {
			ht.Timeout.Duration = 20 * time.Second
		}
		ht.client = &http.Client{Timeout: ht.Timeout.Duration}
	})
	resp, err := ht.client.Do(req)
	if err != nil {
		e := skogul.Error{Source: "http sender", Reason: "unable to Do request", Next: err}
		log.Print(e)
		return e
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		e := skogul.Error{Source: "http sender", Reason: fmt.Sprintf("non-OK status code from target: %d", resp.StatusCode)}
		log.Print(e)
		return e
	}
	return nil
}
