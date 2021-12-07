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
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/stats"
)

var httpLog = skogul.Logger("sender", "http")

/*
HTTP sender POSTs the Skogul JSON-encoded data to the provided URL.
*/
type HTTP struct {
	URL              string            `doc:"Fully qualified URL to send data to." example:"http://localhost:6081/ https://user:password@[::1]:6082/"`
	Headers          map[string]string `doc:"HTTP headers to be added to every request"`
	Timeout          skogul.Duration   `doc:"HTTP timeout."`
	Insecure         bool              `doc:"Disable TLS certificate validation."`
	ConnsPerHost     int               `doc:"Max concurrent connections per host. Should reflect ulimit -n. Defaults to unlimited."`
	IdleConnsPerHost int               `doc:"Max idle connections retained per host. Should reflect expected concurrency. Defaults to 2 + runtime.NumCPU."`
	RootCA           string            `doc:"Path to an alternate root CA used to verify server certificates. Leave blank to use system defaults."`
	Certfile         string            `doc:"Path to certificate file for TLS Client Certificate."`
	Keyfile          string            `doc:"Path to key file for TLS Client Certificate."`
	ok               bool              // set to OK if init worked. FIXME: Should Verify() check if this is OK? I'm thinking yes.
	stats            *httpStats
	ticker           *time.Ticker
	once             sync.Once
	client           *http.Client
}

type httpStats struct {
	Received          uint64         // Metrics received.
	Sent              uint64         // Metrics successfully sent.
	Errors            uint64         // Generic error cases in the module, such as failing to initialize or marshal/unmarshal data.
	RequestErrors     uint64         // Errors during requests, such as not being able to connect to a remote host.
	HttpResponseError map[int]uint64 // Error response codes from HTTP requests. Basically anything != 2XX.
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
		httpLog.Printf("unable to open alternate root CA: %v", err)
		return nil, skogul.Error{Source: "http sender", Reason: "unable to open custom root CA", Next: err}
	}
	defer func() {
		fd.Close()
	}()
	bytes := make([]byte, 1024000)
	n, err := fd.Read(bytes)
	if err != nil {
		httpLog.Printf("unable to read root ca: %v", err)
		return nil, skogul.Error{Source: "http sender", Reason: "unable to read custom root CA", Next: err}
	}
	ok := cp.AppendCertsFromPEM(bytes[:n])
	if !ok {
		httpLog.Printf("unable to append certificate to cert pool")
		return nil, skogul.Error{Source: "http sender", Reason: "unable to append certificate to root CA pool"}
	}
	return cp, nil
}

func (ht *HTTP) loadClientCert() (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(ht.Certfile, ht.Keyfile)
	if err != nil {
		httpLog.WithError(err).Error("Failed to load Client Certificate")
		return nil, err
	}
	return &cert, nil
}

func (ht *HTTP) init() {
	ht.ok = false
	if ht.Timeout.Duration == 0 {
		ht.Timeout.Duration = 20 * time.Second
	}
	if ht.Insecure {
		httpLog.Warning("Disabling certificate validation for HTTP sender - vulnerable to man-in-the-middle")
	}
	iconsph := ht.IdleConnsPerHost

	if iconsph == 0 {
		iconsph = 2 + runtime.NumCPU()
	}

	cp, err := getCertPool(ht.RootCA)
	if err != nil {
		httpLog.WithError(err).Error("Failed to initialize root CA pool")
		return
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: ht.Insecure,
		RootCAs:            cp,
	}

	if ht.Certfile != "" && ht.Keyfile != "" {
		c, err := ht.loadClientCert()
		if err != nil {
			httpLog.WithError(err).Error("Failed to load Client Certificate")
		}
		if c == nil {
			httpLog.Error("Certificate was nil after loading")
			return
		}
		tlsConfig.Certificates = []tls.Certificate{*c}
	}

	tran := http.Transport{
		TLSClientConfig:     tlsConfig,
		MaxConnsPerHost:     ht.ConnsPerHost,
		MaxIdleConnsPerHost: iconsph,
	}

	// Initialize the map if empty in config so we
	// can add headers programmatically.
	if ht.Headers == nil {
		ht.Headers = make(map[string]string)
	}
	contentTypeHeaderKey := ""
	contentTypeHeaderVal := "application/json"
	for header, val := range ht.Headers {
		if http.CanonicalHeaderKey(header) == http.CanonicalHeaderKey("content-type") {
			contentTypeHeaderKey = header
			contentTypeHeaderVal = val
			break
		}
	}
	if contentTypeHeaderKey != http.CanonicalHeaderKey(contentTypeHeaderKey) {
		// Enforce setting the header in the canonicalized way
		delete(ht.Headers, contentTypeHeaderKey)
	}
	ht.Headers[http.CanonicalHeaderKey("content-type")] = contentTypeHeaderVal

	ht.client = &http.Client{Transport: &tran, Timeout: ht.Timeout.Duration}

	ht.initStats()
	ht.ok = true
}

// initStats initializes the required components and structs for stats collection.
func (ht *HTTP) initStats() {
	ht.stats = &httpStats{
		Received:          0,
		Sent:              0,
		Errors:            0,
		RequestErrors:     0,
		HttpResponseError: make(map[int]uint64),
	}
	ht.ticker = time.NewTicker(stats.DefaultInterval)
	go ht.sendStats()
}

// sendBytes uses a configured HTTP client to
// send a request.
// This makes it possible for other senders to
// reuse the HTTP sender options without having
// to re-implement them.
func (ht *HTTP) sendBytes(b []byte) error {
	if !ht.ok {
		return skogul.Error{Reason: "HTTP sender not in OK state", Source: "http-sender"}
	}

	var buffer bytes.Buffer
	buffer.Write(b)
	ht.once.Do(func() {
		ht.init()
	})
	req, err := http.NewRequest("POST", ht.URL, &buffer)
	if err != nil {
		atomic.AddUint64(&ht.stats.Errors, 1)
		e := skogul.Error{Source: "http sender", Reason: "unable to create request", Next: err}
		httpLog.WithError(e).Error("Failed to create request")
		return e
	}
	for header, value := range ht.Headers {
		req.Header.Add(http.CanonicalHeaderKey(header), value)
	}
	resp, err := ht.client.Do(req)
	if err != nil {
		atomic.AddUint64(&ht.stats.RequestErrors, 1)
		e := skogul.Error{Source: "http sender", Reason: "unable to POST request", Next: err}
		httpLog.WithError(e).Error("Failed to do POST request")
		return e
	}
	if resp.ContentLength > 0 {
		tmp := make([]byte, resp.ContentLength)
		if n, err := io.ReadFull(resp.Body, tmp); err != nil {
			atomic.AddUint64(&ht.stats.Errors, 1)
			httpLog.WithError(err).WithFields(log.Fields{
				"expected": resp.ContentLength,
				"actual":   n,
			}).Error("Failed to read http response body")
			resp.Body.Close()
			return err
		}
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		httpResponseCodeStats := ht.stats.HttpResponseError[resp.StatusCode]
		atomic.AddUint64(&httpResponseCodeStats, 1)
		e := skogul.Error{Source: "http sender", Reason: fmt.Sprintf("non-OK status code from target: %d", resp.StatusCode)}
		httpLog.WithError(e).Error("HTTP response status code was not in OK range")
		return e
	}
	atomic.AddUint64(&ht.stats.Sent, 1)
	return nil
}

// Send POSTS data
func (ht *HTTP) Send(c *skogul.Container) error {
	// This is called both here and in sendBytes to make sure
	// we init the sender before sending. Since the first thing we check
	// is if the init was OK, we have to do the check in both functions.
	// (if we only do it in sendBytes(), the following !ht.ok check would fail
	// on first use, and if we only do it in Send(), internal re-use of
	// sendBytes() would fail their check.)
	ht.once.Do(func() {
		ht.init()
	})
	atomic.AddUint64(&ht.stats.Received, 1)
	if !ht.ok {
		atomic.AddUint64(&ht.stats.Errors, 1)
		e := skogul.Error{Source: "http sender", Reason: "initialization failed"}
		httpLog.Print(e)
		return e
	}
	b, err := json.Marshal(*c)
	if err != nil {
		atomic.AddUint64(&ht.stats.Errors, 1)
		e := skogul.Error{Source: "http sender", Reason: "unable to marshal JSON", Next: err}
		httpLog.WithError(e).Error("Configuring HTTP sender failed")
		return e
	}
	err = ht.sendBytes(b)
	if err != nil {
		return err
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
	if (ht.Certfile != "" && ht.Keyfile == "") || (ht.Certfile == "" && ht.Keyfile != "") {
		return skogul.Error{Source: "http sender", Reason: "Specify both Certfile and Keyfile if either is specified."}
	}
	return nil
}

// GetStats prepares a skogul metric with stats
// for the HTTP sender.
func (ht *HTTP) GetStats() *skogul.Metric {
	now := skogul.Now()
	metric := skogul.Metric{
		Time:     &now,
		Metadata: make(map[string]interface{}),
		Data:     make(map[string]interface{}),
	}
	metric.Metadata["component"] = "sender"
	metric.Metadata["type"] = "HTTP"
	metric.Metadata["identity"] = skogul.Identity[ht]

	metric.Data["received"] = ht.stats.Received
	metric.Data["sent"] = ht.stats.Sent
	metric.Data["errors"] = ht.stats.Errors
	metric.Data["request_errors"] = ht.stats.RequestErrors
	for key, val := range ht.stats.HttpResponseError {
		metric.Data[fmt.Sprintf("http_response_%d", key)] = val
	}
	return &metric
}

// sendStats sets up a forever-running loop which sends stats
// to the global skogul stats channel at the configured interval.
func (ht *HTTP) sendStats() {
	for range ht.ticker.C {
		httpLog.Trace("sending stats")
		stats.Chan <- ht.GetStats()
	}
}
