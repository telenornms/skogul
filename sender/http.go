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
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/telenornms/skogul"
)

/*
HTTP sender POSTs the Skogul JSON-encoded data to the provided URL.
*/
type HTTP struct {
	URL              string          `doc:"Fully qualified URL to send data to." example:"http://localhost:6081/ https://user:password@[::1]:6082/"`
	Timeout          skogul.Duration `doc:"HTTP timeout."`
	Insecure         bool            `doc:"Disable TLS certificate validation."`
	ConnsPerHost     int             `doc:"Max concurrent connections per host. Should reflect ulimit -n. Defaults to unlimited."`
	IdleConnsPerHost int             `doc:"Max idle connections retained per host. Should reflect expected concurrency. Defaults to 2 + runtime.NumCPU."`
	RootCA           string          `doc:"Path to an alternate root CA used to verify server certificates. Leave blank to use system defaults."`
	ok               bool            // set to OK if init worked. FIXME: Should Verify() check if this is OK? I'm thinking yes.
	once             sync.Once
	client           *http.Client
}

// getCertPool reads the file specified in f and returns a CertPool with
// the parsed result, suitable for use as RootCAs. If f is empty, it
// returns nil (with no error), which generally means "use system-wide
// defaults".
//
// Called both during Verify(), and again on the first request. This is
// done to satisfy Verify()'s requirement of not modifying state and it
// being optional. We also need Verify to actually test this, so the user
// can be reasonably certain that a valid configuration is used.
func getCertPool(path string) (*x509.CertPool, error) {
	// this means "use system default"
	if path == "" {
		return nil, nil
	}
	cp := x509.NewCertPool()
	fd, err := os.Open(path)
	if err != nil {
		log.Printf("unable to open alternate root CA: %v", err)
		return nil, skogul.Error{Source: "http sender", Reason: "unable to open custom root CA", Next: err}
	}
	defer func() {
		fd.Close()
	}()
	bytes := make([]byte, 1024000)
	n, err := fd.Read(bytes)
	if err != nil {
		log.Printf("unable to read root ca: %v", err)
		return nil, skogul.Error{Source: "http sender", Reason: "unable to read custom root CA", Next: err}
	}
	ok := cp.AppendCertsFromPEM(bytes[:n])
	if !ok {
		log.Printf("unable to append certificate to cert pool")
		return nil, skogul.Error{Source: "http sender", Reason: "unable to append certificate to root CA pool"}
	}
	return cp, nil
}

func (ht *HTTP) init() {
	ht.ok = false
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

	cp, err := getCertPool(ht.RootCA)
	if err != nil {
		log.Print("Failed to initialize root CA pool")
		return
	}

	tran := http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: ht.Insecure,
			RootCAs:            cp,
		},
		MaxConnsPerHost:     ht.ConnsPerHost,
		MaxIdleConnsPerHost: iconsph,
	}
	ht.client = &http.Client{Transport: &tran, Timeout: ht.Timeout.Duration}
	ht.ok = true
}

// Send POSTS data
func (ht *HTTP) Send(c *skogul.Container) error {
	ht.once.Do(func() {
		ht.init()
	})
	if !ht.ok {
		e := skogul.Error{Source: "http sender", Reason: "initialization failed"}
		log.Print(e)
		return e
	}
	b, err := json.Marshal(*c)
	if err != nil {
		e := skogul.Error{Source: "http sender", Reason: "unable to marshal JSON", Next: err}
		log.Print(e)
		return e
	}
	var buffer bytes.Buffer
	buffer.Write(b)
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
	_, err := getCertPool(ht.RootCA)
	if err != nil {
		return skogul.Error{Source: "http sender", Reason: fmt.Sprintf("failed to read custom root CA (RootCA: %s)", ht.RootCA), Next: err}
	}
	return nil
}
