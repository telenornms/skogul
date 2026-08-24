/*
 * skogul, backoff sender tests
 *
 * Copyright (c) 2019-2026 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *
 * This library is free software; you can redistribute it and/or modify it
* under the terms of the GNU Lesser General Public License as published by the
* Free Software Foundation; either version 2.1 of the License, or (at your
* option) any later version.
 *
* This library is distributed in the hope that it will be useful, but WITHOUT
* ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
* FOR A PARTICULAR PURPOSE.  See the GNU Lesser General Public License for more
* details.
 *
* You should have received a copy of the GNU Lesser General Public License
* along with this library; if not, write to the Free Software Foundation, Inc.,
* 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301  USA
*/

package sender

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/telenornms/skogul"
)

// BackTester will fail until it has failed fails times.
type BackTester struct {
	fails int
}

func (bt *BackTester) Send(c *skogul.Container) error {
	if bt.fails > 0 {
		bt.fails--
		return fmt.Errorf("still failing")
	}
	return nil
}

// TestBackoff tests if backoff works at least a little bit
func TestBackoff(t *testing.T) {
	te := BackTester{fails: 1}
	bo := Backoff{
		Next:    skogul.SenderRef{S: &te},
		Base:    skogul.Duration{Duration: 10 * time.Millisecond},
		Retries: 2,
	}
	c := skogul.Container{}
	err := bo.Send(&c)
	if err != nil {
		t.Errorf("Got error from bo.Send(): %v", err)
	}
	te.fails = 10
	err = bo.Send(&c)
	if err == nil {
		t.Errorf("Didn't get error from bo.Send()")
	}
}

// Retries is a total attempt count, so zero means the send loop never
// runs at all: Send would return success without ever reaching Next,
// silently dropping the container.
func TestBackoffRejectsZeroRetries(t *testing.T) {
	bo := Backoff{
		Next: skogul.SenderRef{S: &BackTester{}},
		Base: skogul.Duration{Duration: 10 * time.Millisecond},
	}
	if err := bo.Verify(); err == nil {
		t.Error("expected Verify() to reject Retries: 0")
	}
}

// A Backoff built in code never passes through Verify, so Send has to
// hold the line too rather than silently dropping the container.
func TestBackoffZeroRetriesStillSends(t *testing.T) {
	sent := false
	bo := Backoff{
		Next: skogul.SenderRef{S: senderFunc(func(c *skogul.Container) error {
			sent = true
			return nil
		})},
		Base: skogul.Duration{Duration: 10 * time.Millisecond},
	}
	if err := bo.Send(&skogul.Container{}); err != nil {
		t.Errorf("Got error from bo.Send(): %v", err)
	}
	if !sent {
		t.Error("bo.Send() returned success without ever calling Next")
	}
}

// senderFunc adapts a function to skogul.Sender.
type senderFunc func(*skogul.Container) error

func (f senderFunc) Send(c *skogul.Container) error { return f(c) }

func TestRetryConfigDefaultEnabled(t *testing.T) {
	cfg := RetryConfig{}
	if !cfg.enabled() {
		t.Error("backoff should be enabled by default (nil BackoffEnabled)")
	}
}

func TestRetryConfigDisabled(t *testing.T) {
	disabled := false
	cfg := RetryConfig{BackoffEnabled: &disabled}
	if cfg.enabled() {
		t.Error("backoff should be disabled when BackoffEnabled is false")
	}
}

func TestRetryConfigDefaults(t *testing.T) {
	cfg := RetryConfig{}

	if cfg.maxRetries() != 3 {
		t.Errorf("maxRetries should default to 3, got %d", cfg.maxRetries())
	}
	if cfg.baseDelay() != 100*time.Millisecond {
		t.Errorf("baseDelay should default to 100ms, got %v", cfg.baseDelay())
	}
	if cfg.maxDelay() != 30*time.Second {
		t.Errorf("maxDelay should default to 30s, got %v", cfg.maxDelay())
	}
	if cfg.multiplier() != 2.0 {
		t.Errorf("multiplier should default to 2.0, got %f", cfg.multiplier())
	}
}

func TestWaitDurationBasic(t *testing.T) {
	cfg := defaultRetryConfig()

	d0 := cfg.waitDuration(0, nil)
	if d0 < 0 || d0 > 100*time.Millisecond {
		t.Errorf("attempt 0 delay should be in [0, 100ms], got %v", d0)
	}

	d1 := cfg.waitDuration(1, nil)
	if d1 < 0 || d1 > 200*time.Millisecond {
		t.Errorf("attempt 1 delay should be in [0, 200ms], got %v", d1)
	}
}

func TestWaitDurationCaps(t *testing.T) {
	cfg := defaultRetryConfig()

	d := cfg.waitDuration(9, nil)
	if d > 30*time.Second {
		t.Errorf("attempt 9 delay should be capped at 30s, got %v", d)
	}
	if d <= 0 {
		t.Errorf("attempt 9 delay should be > 0 due to jitter, got %v", d)
	}
}

func TestWaitDurationRetryAfter(t *testing.T) {
	cfg := defaultRetryConfig()

	retryAfter := 5 * time.Second
	d := cfg.waitDuration(0, &retryAfter)
	if d < 5*time.Second {
		t.Errorf("with Retry-After=5s, delay should be at least 5s, got %v", d)
	}
}

func TestShouldHTTPRetry(t *testing.T) {
	cfg := defaultRetryConfig()

	tests := []struct {
		name       string
		statusCode int
		err        error
		want       bool
	}{
		{"200 OK", 200, nil, false},
		{"400 Bad Request", 400, nil, false},
		{"401 Unauthorized", 401, nil, false},
		{"403 Forbidden", 403, nil, false},
		{"404 Not Found", 404, nil, false},
		{"429 Too Many Requests", 429, nil, true},
		{"408 Request Timeout", 408, nil, true},
		{"500 Internal Server Error", 500, nil, true},
		{"502 Bad Gateway", 502, nil, true},
		{"503 Service Unavailable", 503, nil, true},
		{"504 Gateway Timeout", 504, nil, true},
		{"network error", -1, &net.OpError{Op: "dial", Net: "tcp", Addr: nil, Err: fmt.Errorf("connection refused")}, true},
		// A non-retryable status must win over transient-looking
		// error text: the InfluxDB sender embeds the response body
		// in the error message.
		{"400 with timeout in body", 400, fmt.Errorf("bad response from InfluxDB: 400 Bad Request - partial write: timeout"), false},
		{"404 with connection closed in body", 404, fmt.Errorf("bad response from InfluxDB: 404 Not Found - connection closed"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldHTTPRetry(&cfg, 0, tt.statusCode, tt.err)
			if got != tt.want {
				t.Errorf("shouldHTTPRetry(attempt=0, %d, %v) = %v, want %v", tt.statusCode, tt.err, got, tt.want)
			}
		})
	}
}

func TestShouldHTTPRetryExhausted(t *testing.T) {
	cfg := defaultRetryConfig()
	cfg.MaxRetries = retries(2)
	got := shouldHTTPRetry(&cfg, 2, 500, nil)
	if got {
		t.Error("shouldHTTPRetry should return false when retries exhausted")
	}
}

func TestShouldHTTPRetryDisabled(t *testing.T) {
	disabled := false
	cfg := RetryConfig{BackoffEnabled: &disabled}
	got := shouldHTTPRetry(&cfg, 0, 503, nil)
	if got {
		t.Error("shouldHTTPRetry should return false when backoff disabled")
	}
}

func TestShouldNetworkRetry(t *testing.T) {
	cfg := defaultRetryConfig()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"connection refused", &net.OpError{Op: "dial", Net: "tcp", Addr: nil, Err: fmt.Errorf("connection refused")}, true},
		{"connection reset", fmt.Errorf("connection reset by peer"), true},
		{"EOF", io.EOF, true},
		{"context deadline", context.DeadlineExceeded, true},
		{"no such host", fmt.Errorf("dial tcp: lookup bad.example.com: no such host"), true},
		{"broker not available", fmt.Errorf("broker not available"), true},
		{"channel closed", fmt.Errorf("channel closed"), true},
		{"encoding error", fmt.Errorf("json: unsupported type"), false},
		{"unknown error", fmt.Errorf("some random error"), false},
		{"no error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldNetworkRetry(&cfg, 0, tt.err)
			if got != tt.want {
				t.Errorf("shouldNetworkRetry(err=%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestHTTPStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"nil error", nil, -1},
		{"plain error", fmt.Errorf("oops"), -1},
		{"RetryableHTTPError 404", &RetryableHTTPError{Err: fmt.Errorf("not found"), StatusCode: 404}, 404},
		{"wrapped RetryableHTTPError 503", fmt.Errorf("send failed: %w", &RetryableHTTPError{Err: fmt.Errorf("service unavailable"), StatusCode: 503}), 503},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := httpStatusCode(tt.err)
			if code != tt.wantCode {
				t.Errorf("httpStatusCode() = %d, want %d", code, tt.wantCode)
			}
		})
	}
}

func TestParseHTTPRetryAfter(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		wantMin time.Duration
		wantMax time.Duration
		wantNil bool
	}{
		{"absent header", "", 0, 0, true},
		{"integer seconds (60)", "60", 60 * time.Second, 60 * time.Second, false},
		{"integer seconds (120)", "120", 120 * time.Second, 120 * time.Second, false},
		{"past date", "Thu, 01 Jan 1970 00:00:00 GMT", 0, 0, true},
		{"future date", "Thu, 01 Jan 2099 00:00:00 GMT", 0, time.Hour * 24 * 365 * 73, false},
		{"invalid", "not-a-date", 0, 0, true},
		{"empty", "", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{Header: http.Header{}}
			if tt.header != "" {
				resp.Header.Set("Retry-After", tt.header)
			}
			d := parseHTTPRetryAfter(resp)

			if tt.wantNil {
				if d != nil {
					t.Errorf("parseHTTPRetryAfter(%q) returned %v, expected nil", tt.header, *d)
				}
				return
			}
			if d == nil {
				t.Errorf("parseHTTPRetryAfter(%q) returned nil, expected value", tt.header)
				return
			}
			if *d < tt.wantMin || *d > tt.wantMax {
				t.Errorf("parseHTTPRetryAfter(%q) = %v, expected in [%v, %v]", tt.header, *d, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func retries(n int) *int {
	return &n
}

func defaultRetryConfig() RetryConfig {
	enabled := true
	return RetryConfig{
		BackoffEnabled:    &enabled,
		MaxRetries:        retries(3),
		BaseDelay:         skogulDur(100 * time.Millisecond),
		MaxDelay:          skogulDur(30 * time.Second),
		BackoffMultiplier: 2.0,
	}
}

func skogulDur(d time.Duration) skogul.Duration {
	return skogul.Duration{Duration: d}
}

func TestIsTransientError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{"connection refused", fmt.Errorf("dial tcp: connection refused"), true},
		{"connection refused OpError", &net.OpError{Err: fmt.Errorf("connection refused")}, true},
		{"connection reset", fmt.Errorf("connection reset by peer"), true},
		{"EOF", io.EOF, true},
		{"unexpected EOF", fmt.Errorf("unexpected EOF"), true},
		{"broker not available", fmt.Errorf("kafka: broker not available"), true},
		{"paho not connected", fmt.Errorf("not Connected"), true},
		{"nats timeout", fmt.Errorf("nats: timeout"), true},
		{"publish timed out", fmt.Errorf("MQTT publish to topic foo timed out"), true},
		{"channel closed", fmt.Errorf("amqp: channel closed"), true},
		{"no such host", fmt.Errorf("lookup example.com: no such host"), true},
		{"NXDOMAIN DNSError", &net.DNSError{Err: "no answer", Name: "example.com", IsNotFound: true}, true},
		{"context exceeded", context.DeadlineExceeded, true},
		{"syscall ECONNRESET", syscall.ECONNRESET, true},
		{"syscall EPIPE", syscall.EPIPE, true},
		{"EPIPE OpError", &net.OpError{Op: "write", Err: syscall.EPIPE}, true},
		{"syscall EHOSTUNREACH", syscall.EHOSTUNREACH, true},
		{"syscall ENETUNREACH", syscall.ENETUNREACH, true},
		{"paho status error", fmt.Errorf("status can only transition to connecting from disconnected"), true},
		{"wrapped transient", fmt.Errorf("send failed: %w", syscall.ECONNREFUSED), true},
		{"no match", fmt.Errorf("some other error"), false},
		{"nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTransientError(tt.err)
			if got != tt.expect {
				t.Errorf("isTransientError(%v) = %v, want %v", tt.err, got, tt.expect)
			}
		})
	}
}

func TestSyscallErrors(t *testing.T) {
	cfg := defaultRetryConfig()

	if !shouldNetworkRetry(&cfg, 0, syscall.ECONNREFUSED) {
		t.Error("ECONNREFUSED should be retryable")
	}

	if shouldNetworkRetry(&cfg, 0, syscall.ENOENT) {
		t.Error("ENOENT should not be retryable")
	}
}

func TestFullJitterVariety(t *testing.T) {
	cfg := defaultRetryConfig()

	results := make([]float64, 100)
	for i := range 100 {
		d := cfg.waitDuration(1, nil)
		results[i] = float64(d)
	}

	unique := make(map[float64]bool)
	for _, r := range results {
		unique[r] = true
	}
	if len(unique) < 90 {
		t.Errorf("Full Jitter should produce mostly unique values, got %d unique out of 100", len(unique))
	}

	var min, max time.Duration = math.MaxInt64, 0
	for _, r := range results {
		d := time.Duration(r)
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}

	if min < 0 || max > 200*time.Millisecond {
		t.Errorf("All values should be in [0, 200ms], min=%v, max=%v", min, max)
	}
}

// fastRetryConfig is a resolved config with tiny delays so retry tests
// run fast.
func fastRetryConfig(maxRetries int) RetryConfig {
	enabled := true
	return RetryConfig{
		BackoffEnabled:    &enabled,
		MaxRetries:        retries(maxRetries),
		BaseDelay:         skogulDur(time.Millisecond),
		MaxDelay:          skogulDur(10 * time.Millisecond),
		BackoffMultiplier: 2.0,
	}
}

var testRetryLog = skogul.Logger("sender", "backoff-test")

func TestRetryNetworkEventuallySucceeds(t *testing.T) {
	cfg := fastRetryConfig(3)
	calls := 0
	err := retryNetwork(&cfg, testRetryLog, func() error {
		calls++
		if calls < 3 {
			return fmt.Errorf("connection refused")
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected success after retries, got: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 attempts, got %d", calls)
	}
}

func TestRetryNetworkNonRetryable(t *testing.T) {
	cfg := fastRetryConfig(3)
	calls := 0
	err := retryNetwork(&cfg, testRetryLog, func() error {
		calls++
		return fmt.Errorf("json: unsupported type")
	})
	if err == nil {
		t.Error("expected error for non-retryable failure")
	}
	if calls != 1 {
		t.Errorf("non-retryable error should not be retried, got %d attempts", calls)
	}
}

func TestRetryNetworkExhausted(t *testing.T) {
	cfg := fastRetryConfig(2)
	calls := 0
	err := retryNetwork(&cfg, testRetryLog, func() error {
		calls++
		return fmt.Errorf("connection refused")
	})
	if err == nil {
		t.Error("expected error after exhausting retries")
	}
	if calls != 3 {
		t.Errorf("MaxRetries=2 should give 3 attempts, got %d", calls)
	}
	if !strings.Contains(err.Error(), "giving up after 3 attempts") {
		t.Errorf("error should mention attempt count, got: %v", err)
	}
}

func TestRetryHTTPStatusCodes(t *testing.T) {
	cfg := fastRetryConfig(3)
	calls := 0
	err := retryHTTP(&cfg, testRetryLog, func() error {
		calls++
		if calls < 2 {
			return &RetryableHTTPError{Err: fmt.Errorf("server error"), StatusCode: 500}
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected success after retrying 500, got: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 attempts, got %d", calls)
	}

	calls = 0
	err = retryHTTP(&cfg, testRetryLog, func() error {
		calls++
		return &RetryableHTTPError{Err: fmt.Errorf("not found"), StatusCode: 404}
	})
	if err == nil {
		t.Error("expected error for 404")
	}
	if calls != 1 {
		t.Errorf("404 should not be retried, got %d attempts", calls)
	}
}

func TestWaitDurationRetryAfterCappedAtMaxDelay(t *testing.T) {
	cfg := defaultRetryConfig()

	retryAfter := 10 * time.Minute
	d := cfg.waitDuration(0, &retryAfter)
	if d > cfg.MaxDelay.Duration {
		t.Errorf("Retry-After must be capped at MaxDelay (%v), got %v", cfg.MaxDelay.Duration, d)
	}
}

// Retry-After must be spread over a random band, so that every client
// told the same number by an overloaded server does not come back in the
// same instant.
func TestWaitDurationRetryAfterJittered(t *testing.T) {
	cfg := defaultRetryConfig()
	retryAfter := 5 * time.Second

	seen := make(map[time.Duration]bool)
	for range 50 {
		d := cfg.waitDuration(0, &retryAfter)
		if d < retryAfter {
			t.Fatalf("Retry-After must be honoured as a floor, got %v < %v", d, retryAfter)
		}
		if d > retryAfter+time.Duration(float64(retryAfter)*retryAfterJitter) {
			t.Fatalf("jitter band too wide: got %v for Retry-After %v", d, retryAfter)
		}
		seen[d] = true
	}
	if len(seen) < 10 {
		t.Errorf("Retry-After delays should vary, got %d distinct values in 50 draws", len(seen))
	}
}

// When the hint is already at MaxDelay there is no room to jitter above
// it, so the band has to sit below the cap instead.
func TestWaitDurationRetryAfterAtMaxDelayStillJitters(t *testing.T) {
	cfg := defaultRetryConfig()
	retryAfter := 10 * time.Minute // way past MaxDelay

	seen := make(map[time.Duration]bool)
	for range 50 {
		d := cfg.waitDuration(0, &retryAfter)
		if d > cfg.MaxDelay.Duration {
			t.Fatalf("delay must be capped at MaxDelay (%v), got %v", cfg.MaxDelay.Duration, d)
		}
		seen[d] = true
	}
	if len(seen) < 10 {
		t.Errorf("delays at the cap should still vary, got %d distinct values", len(seen))
	}
}

func TestMaxRetriesUnsetAndZero(t *testing.T) {
	var cfg RetryConfig
	if cfg.maxRetries() != defaultMaxRetries {
		t.Errorf("unset MaxRetries should default to %d, got %d", defaultMaxRetries, cfg.maxRetries())
	}

	// An explicit 0 must mean "no retries", not "unset".
	cfg.MaxRetries = retries(0)
	if cfg.maxRetries() != 0 {
		t.Errorf("MaxRetries=0 should mean no retries, got %d", cfg.maxRetries())
	}
	calls := 0
	err := retryNetwork(&cfg, testRetryLog, func() error {
		calls++
		return fmt.Errorf("connection refused")
	})
	if err == nil {
		t.Error("expected an error")
	}
	if calls != 1 {
		t.Errorf("MaxRetries=0 should give exactly 1 attempt, got %d", calls)
	}
}

// MaxElapsed bounds the total time a single send can spend retrying,
// regardless of MaxRetries.
func TestMaxElapsedStopsRetrying(t *testing.T) {
	cfg := fastRetryConfig(100)
	cfg.BaseDelay = skogulDur(20 * time.Millisecond)
	cfg.MaxDelay = skogulDur(20 * time.Millisecond)
	cfg.MaxElapsed = skogulDur(100 * time.Millisecond)

	calls := 0
	start := time.Now()
	err := retryNetwork(&cfg, testRetryLog, func() error {
		calls++
		return fmt.Errorf("connection refused")
	})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected an error once the budget ran out")
	}
	if !strings.Contains(err.Error(), "retry budget") {
		t.Errorf("error should mention the budget, got: %v", err)
	}
	if calls >= 100 {
		t.Errorf("MaxElapsed should have cut the retries short, got %d attempts", calls)
	}
	if elapsed > time.Second {
		t.Errorf("MaxElapsed=100ms should not take %v", elapsed)
	}
}

func TestVerifyRetry(t *testing.T) {
	tests := []struct {
		name    string
		cfg     RetryConfig
		wantErr bool
	}{
		{"zero value", RetryConfig{}, false},
		{"sane", defaultRetryConfig(), false},
		{"no retries", RetryConfig{MaxRetries: retries(0)}, false},
		{"negative retries", RetryConfig{MaxRetries: retries(-1)}, true},
		{"negative base delay", RetryConfig{BaseDelay: skogulDur(-time.Second)}, true},
		{"negative max delay", RetryConfig{MaxDelay: skogulDur(-time.Second)}, true},
		{"negative max elapsed", RetryConfig{MaxElapsed: skogulDur(-time.Second)}, true},
		{"negative multiplier", RetryConfig{BackoffMultiplier: -2}, true},
		{"shrinking multiplier", RetryConfig{BackoffMultiplier: 0.5}, true},
		{"multiplier of 1", RetryConfig{BackoffMultiplier: 1}, false},
		{"base delay past max delay", RetryConfig{BaseDelay: skogulDur(time.Minute), MaxDelay: skogulDur(time.Second)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.verifyRetry()
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyRetry() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// A send that failed after part of the payload reached the receiver must
// not be retried by default: the retry re-sends the whole payload, so the
// receiver would see the partial data followed by a complete copy.
func TestPartialWriteNotRetriedByDefault(t *testing.T) {
	cfg := fastRetryConfig(3)
	calls := 0
	err := retryNetwork(&cfg, testRetryLog, func() error {
		calls++
		return partialWrite(fmt.Errorf("unable to send data: %w", syscall.ECONNRESET), 10, 100)
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	if calls != 1 {
		t.Errorf("a partial write should not be retried, got %d attempts", calls)
	}
	if !strings.Contains(err.Error(), "wrote 10 of 100 bytes") {
		t.Errorf("error should say how much was written, got: %v", err)
	}
	// The underlying error must stay reachable for callers.
	if !errors.Is(err, syscall.ECONNRESET) {
		t.Error("partialWrite must not break the error chain")
	}
}

// Writing nothing at all is not a partial write - there is nothing to
// duplicate, so it stays retryable.
func TestZeroByteWriteIsRetried(t *testing.T) {
	cfg := fastRetryConfig(2)
	calls := 0
	err := retryNetwork(&cfg, testRetryLog, func() error {
		calls++
		return partialWrite(fmt.Errorf("unable to send data: %w", syscall.ECONNRESET), 0, 100)
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	if calls != 3 {
		t.Errorf("a zero-byte write should be retried, got %d attempts", calls)
	}
}

// RetryPartialWrites opts back in to at-least-once delivery.
func TestPartialWriteRetriedWhenEnabled(t *testing.T) {
	cfg := fastRetryConfig(2)
	yes := true
	cfg.RetryPartialWrites = &yes
	calls := 0
	err := retryNetwork(&cfg, testRetryLog, func() error {
		calls++
		return partialWrite(fmt.Errorf("unable to send data: %w", syscall.ECONNRESET), 10, 100)
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	if calls != 3 {
		t.Errorf("RetryPartialWrites should allow retries, got %d attempts", calls)
	}
}

// The guard applies to the HTTP path too, even though no HTTP sender
// produces partial writes today.
func TestPartialWriteNotRetriedOnHTTPPath(t *testing.T) {
	cfg := fastRetryConfig(3)
	if shouldHTTPRetry(&cfg, 0, -1, partialWrite(fmt.Errorf("timeout"), 10, 100)) {
		t.Error("a partial write should not be retried on the HTTP path either")
	}
}

// The Net sender must not re-send a payload the receiver already got part
// of, but must still retry when the connection could not be established.
func TestNetSenderPartialWrite(t *testing.T) {
	// A listener that accepts, reads a little, then resets, so the
	// sender's write fails with data already delivered.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("unable to listen: %v", err)
	}
	defer l.Close()

	var accepted atomic.Int64
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			conn, connerr := l.Accept()
			if connerr != nil {
				return
			}
			accepted.Add(1)
			// Read a little, then reset the connection mid-payload.
			buf := make([]byte, 1)
			conn.Read(buf)
			if tcp, ok := conn.(*net.TCPConn); ok {
				tcp.SetLinger(0)
			}
			conn.Close()
		}
	}()

	n := Net{
		Address:     l.Addr().String(),
		Network:     "tcp",
		RetryConfig: fastRetryConfig(3),
	}
	// A payload large enough that the write cannot complete into the
	// socket buffer before the reset lands.
	metrics := make([]*skogul.Metric, 0, 20000)
	for i := range 20000 {
		metrics = append(metrics, &skogul.Metric{
			Metadata: map[string]any{"key": "some reasonably long metadata value"},
			Data:     map[string]any{"value": i},
		})
	}
	c := skogul.Container{Metrics: metrics}

	// Whether the reset lands mid-write or before the first byte is up
	// to the kernel, so assert conditionally on what actually happened -
	// deterministic either way, and never flaky.
	err = n.Send(&c)
	switch {
	case err == nil:
		t.Log("write completed before the reset; nothing to assert here")
	case strings.Contains(err.Error(), "bytes before failing"):
		if got := accepted.Load(); got != 1 {
			t.Errorf("a partial write must not be retried, got %d connection attempts", got)
		}
	default:
		t.Logf("write failed without delivering data, retried as expected: %v", err)
	}

	l.Close()
	<-done

	// Dial failures, by contrast, must be retried: nothing was sent.
	n2 := Net{
		Address:     l.Addr().String(),
		Network:     "tcp",
		RetryConfig: fastRetryConfig(2),
	}
	if err := n2.Send(&skogul.Container{}); err == nil {
		t.Error("expected an error when dialling a closed port")
	} else if !strings.Contains(err.Error(), "giving up after 3 attempts") {
		t.Errorf("dial failures should be retried, got: %v", err)
	}
}
