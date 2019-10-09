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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/KristianLyng/skogul"
)

/*
HTTP sender POSTs the Skogul JSON-encoded data to the provided URL.
*/
type HTTP struct {
	URL              string          `doc:"Fully qualified URL to send data to." example:"http://localhost:6081/ https://user:password@[::1]:6082/"`
	Timeout          skogul.Duration `doc:"HTTP timeout."`
	Insecure         bool            `doc:"Disable TLS certificate validation."`
	ConnsPerHost     int             `doc:"Max concurrent connections per host. Should reflect ulimit -n. Defaults to unlimited."`
	IdleConnsPerHost int             `doc:"Mas idle connections retained per host. Should reflect expected concurrency. Defaults to 2 + runtime.NumCPU."`
	once             sync.Once
	client           *http.Client
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
	ht.once.Do(func() {
		if ht.Timeout.Duration == 0 {
			ht.Timeout.Duration = 20 * time.Second
		}
		if ht.Insecure {
			log.Print("Warning: Disabeling certificate validation for HTTP sender - vulnerable to man-in-the-middle")
		}
		iconsph := ht.IdleConnsPerHost

		if iconsph == 0 {
			iconsph = 2 + runtime.NumCPU()
		}
		tran := http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ht.Insecure,
			},
			MaxConnsPerHost:     ht.ConnsPerHost,
			MaxIdleConnsPerHost: iconsph,
		}
		ht.client = &http.Client{Transport: &tran, Timeout: ht.Timeout.Duration}
	})
	resp, err := ht.client.Post(ht.URL, "application/json", &buffer)
	if err != nil {
		e := skogul.Error{Source: "http sender", Reason: "unable to POST request", Next: err}
		log.Print(e)
		return e
	}
	if resp.ContentLength > 0 {
		tmp := make([]byte, resp.ContentLength)
		if n, err := io.ReadFull(resp.Body, tmp); err != nil {
			log.Printf("Read %d of %d bytes, returned: %v", n, resp.ContentLength, err)
			resp.Body.Close()
			return err
		}
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		e := skogul.Error{Source: "http sender", Reason: fmt.Sprintf("non-OK status code from target: %d", resp.StatusCode)}
		log.Print(e)
		return e
	}
	return nil
}

// Verify checks that configuration is sensible
func (ht *HTTP) Verify() error {
	if ht.URL == "" {
		return skogul.Error{Source: "http sender", Reason: "no URL specified"}
	}
	return nil
}
