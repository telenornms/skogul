/*
 * skogul, backoff helper
 *
* Copyright (c) 2026 Telenor Norge AS
 *
* This library is free software; you can redistribute it and/or modify it under
* the terms of the GNU Lesser General Public License as published by the Free
* Software Foundation; either version 2.1 of the License, or (at your option)
* any later version.
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
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/telenornms/skogul"
)

// RetryConfig holds the parameters for HTTP and network retry backoff.
// Embedded flat into sender structs so all fields appear at top level in JSON.
//
// Unset (zero) fields fall back to defaults at use time, so a zero
// RetryConfig is ready to use.
type RetryConfig struct {
	BackoffEnabled    *bool           `doc:"Retry sends that fail with what looks like a transient error. Enabled by default. Disable for fail-fast, at-most-once delivery."`
	MaxRetries        *int            `doc:"Number of retries after the first attempt, so the default of 3 means at most 4 attempts. 0 disables retries."`
	BaseDelay         skogul.Duration `doc:"Starting point for the delay between attempts, doubling (see BackoffMultiplier) for each attempt. The actual delay is drawn at random between zero and this computed value, so it averages half of it. Default: 100ms."`
	MaxDelay          skogul.Duration `doc:"Ceiling for the computed delay, and for Retry-After hints from a server. Default: 30s."`
	BackoffMultiplier float64         `doc:"What the delay is multiplied by for each attempt. Must be at least 1. Default: 2.0."`
	MaxElapsed        skogul.Duration `doc:"Time budget for the delays of a single send. Retrying stops once the next delay would take it past this. Note that an attempt already under way is never cut short, so the send itself can run past the budget by up to one attempt (bounded by the sender's own timeout, if it has one). Default: unset, in which case the bound is MaxRetries delays of at most MaxDelay each."`

	RetryPartialWrites *bool `doc:"Retry a send that failed after part of the payload had already reached the receiver. Only applies to the mnr and net senders, which write an unframed byte stream; ignored elsewhere. Disabled by default, since a retry re-sends the whole payload and the receiver would see the partial data followed by a complete copy. Enable for at-least-once delivery on receivers that tolerate duplicates."`
}

const (
	defaultMaxRetries        = 3
	defaultBaseDelay         = 100 * time.Millisecond
	defaultMaxDelay          = 30 * time.Second
	defaultBackoffMultiplier = 2.0

	// retryAfterJitter is the fraction of a Retry-After hint used as
	// the width of the random band applied on top of it.
	retryAfterJitter = 0.1

	// defaultDialTimeout bounds a single connection attempt in the
	// senders that dial per send. An unbounded net.Dial can sit in the
	// kernel's SYN retry for minutes, which each retry would then
	// repeat.
	defaultDialTimeout = 10 * time.Second
)

// enabled returns true if backoff is enabled.
// Returns true when BackoffEnabled is nil (not set by user) or explicitly true.
func (rc *RetryConfig) enabled() bool {
	return rc.BackoffEnabled == nil || *rc.BackoffEnabled
}

// maxRetries returns the configured number of retries, defaulting when
// unset. An explicit 0 means "no retries", so this is a pointer: with a
// plain int there is no way to tell an unset field from a deliberate 0.
func (rc *RetryConfig) maxRetries() int {
	if rc.MaxRetries == nil {
		return defaultMaxRetries
	}
	if *rc.MaxRetries < 0 {
		// Rejected by verifyRetry, but be safe if Verify was skipped.
		return 0
	}
	return *rc.MaxRetries
}

// retryPartialWrites reports whether a send that failed part-way through
// writing the payload may be retried. Off unless explicitly enabled.
func (rc *RetryConfig) retryPartialWrites() bool {
	return rc.RetryPartialWrites != nil && *rc.RetryPartialWrites
}

// maxElapsed returns the total time budget for a send, or 0 if unset.
func (rc *RetryConfig) maxElapsed() time.Duration {
	if rc.MaxElapsed.Duration < 0 {
		return 0
	}
	return rc.MaxElapsed.Duration
}

// verifyRetry validates the retry settings. Senders with a RetryConfig
// call this from their own Verify(). It is not called Verify() itself:
// that name would be promoted to the embedding sender and silently
// stand in for a missing sender-level Verify().
func (rc *RetryConfig) verifyRetry() error {
	if rc.MaxRetries != nil && *rc.MaxRetries < 0 {
		return fmt.Errorf("MaxRetries must be zero or more, got %d", *rc.MaxRetries)
	}
	if rc.BaseDelay.Duration < 0 {
		return fmt.Errorf("BaseDelay must not be negative, got %v", rc.BaseDelay.Duration)
	}
	if rc.MaxDelay.Duration < 0 {
		return fmt.Errorf("MaxDelay must not be negative, got %v", rc.MaxDelay.Duration)
	}
	if rc.MaxElapsed.Duration < 0 {
		return fmt.Errorf("MaxElapsed must not be negative, got %v", rc.MaxElapsed.Duration)
	}
	// A multiplier below 1 shrinks the delay on each attempt, and a
	// negative one alternates sign, which makes time.Sleep return
	// immediately and turns the backoff into a hot retry loop.
	if rc.BackoffMultiplier != 0 && rc.BackoffMultiplier < 1 {
		return fmt.Errorf("BackoffMultiplier must be at least 1, got %v", rc.BackoffMultiplier)
	}
	// Only flagged when both are set: an unset BaseDelay falling back to
	// a default larger than an explicit MaxDelay is harmless, since the
	// delay is capped at MaxDelay either way.
	if rc.BaseDelay.Duration > 0 && rc.MaxDelay.Duration > 0 && rc.BaseDelay.Duration > rc.MaxDelay.Duration {
		return fmt.Errorf("BaseDelay (%v) must not exceed MaxDelay (%v)", rc.BaseDelay.Duration, rc.MaxDelay.Duration)
	}
	return nil
}

func (rc *RetryConfig) baseDelay() time.Duration {
	if rc.BaseDelay.Duration == 0 {
		return defaultBaseDelay
	}
	return rc.BaseDelay.Duration
}

func (rc *RetryConfig) maxDelay() time.Duration {
	if rc.MaxDelay.Duration == 0 {
		return defaultMaxDelay
	}
	return rc.MaxDelay.Duration
}

func (rc *RetryConfig) multiplier() float64 {
	if rc.BackoffMultiplier == 0 {
		return defaultBackoffMultiplier
	}
	return rc.BackoffMultiplier
}

// waitDuration calculates the jittered exponential backoff delay.
// Uses Full Jitter (AWS-recommended): random(0, 1) * computed_delay.
// If retryAfter is set (from Retry-After header) and is longer than the
// computed delay, it is honoured as a floor - but never beyond MaxDelay,
// since a remote server must not be able to stall the pipeline
// indefinitely, and with a random band on top, see jitterRetryAfter.
func (rc *RetryConfig) waitDuration(attempt int, retryAfter *time.Duration) time.Duration {
	maxDelay := rc.maxDelay()
	delay := float64(rc.baseDelay()) * math.Pow(rc.multiplier(), float64(attempt))
	if delay > float64(maxDelay) {
		delay = float64(maxDelay)
	}
	delay *= rand.Float64() // Full Jitter: random in [0, delay]
	if retryAfter != nil {
		ra := *retryAfter
		if ra > maxDelay {
			ra = maxDelay
		}
		if float64(ra) > delay {
			delay = rc.jitterRetryAfter(ra)
		}
	}
	return time.Duration(delay)
}

// jitterRetryAfter spreads a Retry-After hint over a random band, so that
// every client told the same number by an overloaded server does not come
// back in the same instant. The band normally sits above the hint, to
// still honour it as a floor; if that would breach MaxDelay it is placed
// below the cap instead.
func (rc *RetryConfig) jitterRetryAfter(retryAfter time.Duration) float64 {
	maxDelay := float64(rc.maxDelay())
	span := float64(retryAfter) * retryAfterJitter
	low := float64(retryAfter)
	if low+span > maxDelay {
		low = maxDelay - span
		if low < 0 {
			low, span = 0, maxDelay
		}
	}
	return low + rand.Float64()*span
}

// retryHTTP runs send, retrying with backoff on transient network errors
// and retryable HTTP status codes (5xx, 429, 408). send always runs at
// least once. Honours Retry-After hints found in the error chain as a
// floor, capped at MaxDelay and jittered, see waitDuration.
//
// Note that a retry can duplicate data: a request that reached the server
// but whose reply was lost is indistinguishable from one that never
// arrived, so this trades at-most-once for at-least-once delivery.
func retryHTTP(cfg *RetryConfig, log *logrus.Entry, send func() error) error {
	return retryLoop(cfg, log, send, func(attempt int, err error) bool {
		return shouldHTTPRetry(cfg, attempt, httpStatusCode(err), err)
	})
}

// retryNetwork runs send, retrying with backoff on transient network
// errors. send always runs at least once. Errors marked with
// partialWrite are not retried unless RetryPartialWrites is set, see
// partialWriteError.
func retryNetwork(cfg *RetryConfig, log *logrus.Entry, send func() error) error {
	return retryLoop(cfg, log, send, func(attempt int, err error) bool {
		return shouldNetworkRetry(cfg, attempt, err)
	})
}

// retryLoop is the common retry loop behind retryHTTP and retryNetwork.
// It returns nil on success, the error as-is if no retry was performed,
// or the last error annotated with the attempt count once retries are
// exhausted, the time budget runs out, or a non-retryable error shows up.
func retryLoop(cfg *RetryConfig, log *logrus.Entry, send func() error, shouldRetry func(attempt int, err error) bool) error {
	start := time.Now()
	for attempt := 0; ; attempt++ {
		err := send()
		if err == nil {
			return nil
		}
		if !shouldRetry(attempt, err) {
			if attempt > 0 {
				return fmt.Errorf("giving up after %d attempts: %w", attempt+1, err)
			}
			return err
		}
		delay := cfg.waitDuration(attempt, extractRetryAfterFromError(err))
		// Stop if the next delay would take us past the total budget.
		// Nothing cancels a send once it is under way, so the budget
		// is the only thing keeping a slow remote from parking a
		// pipeline goroutine for maxRetries*MaxDelay.
		if budget := cfg.maxElapsed(); budget > 0 {
			if elapsed := time.Since(start); elapsed+delay > budget {
				return fmt.Errorf("giving up after %d attempts, %v of %v retry budget used: %w", attempt+1, elapsed.Round(time.Millisecond), budget, err)
			}
		}
		log.WithError(err).Warnf("send failed (attempt %d of %d), retrying in %v", attempt+1, cfg.maxRetries()+1, delay)
		time.Sleep(delay)
	}
}

// --- Partial writes ---

// partialWriteError marks a send that failed after part of the payload
// had already reached the receiver. A retry re-sends the whole payload,
// so the receiver would see the partial data followed by a complete copy.
// Senders that write an unframed, unacknowledged payload (MnR, Net) use
// this to keep the retry from duplicating data; whether it is retried
// anyway is up to RetryPartialWrites.
type partialWriteError struct {
	Err     error
	written int
	total   int
}

func (p *partialWriteError) Error() string {
	return fmt.Sprintf("%s (wrote %d of %d bytes before failing)", p.Err.Error(), p.written, p.total)
}

func (p *partialWriteError) Unwrap() error {
	return p.Err
}

// partialWrite marks err as a partial write of a total-byte payload.
// Writing nothing at all is not a partial write - there is nothing to
// duplicate - so err is returned unchanged, and stays retryable.
func partialWrite(err error, written, total int) error {
	if written <= 0 {
		return err
	}
	return &partialWriteError{Err: err, written: written, total: total}
}

// isPartialWrite reports whether the error chain carries a partial write.
func isPartialWrite(err error) bool {
	var p *partialWriteError
	return errors.As(err, &p)
}

// retryablePartialWrite reports whether a partial write, if that is what
// err is, may be retried under this config.
func retryablePartialWrite(cfg *RetryConfig, err error) bool {
	return cfg.retryPartialWrites() || !isPartialWrite(err)
}

// shouldHTTPRetry determines if an HTTP error/status is retryable.
// Retries on: 5xx, 429, 408, and transient network errors.
// Does not retry on: 4xx (except 429/408).
func shouldHTTPRetry(cfg *RetryConfig, attempt int, statusCode int, err error) bool {
	if !cfg.enabled() || attempt >= cfg.maxRetries() {
		return false
	}
	if !retryablePartialWrite(cfg, err) {
		return false
	}

	switch {
	case statusCode >= 500 || statusCode == 429 || statusCode == 408:
		return true
	case statusCode >= 100:
		// A definite, non-retryable HTTP status (e.g. 4xx). Don't
		// fall through to text matching: the error text can contain
		// the response body (the InfluxDB sender does this), and a
		// body mentioning e.g. "timeout" must not trigger a retry
		// of a permanent failure.
		return false
	default:
		// No HTTP status at all: transport-level failure, retry if
		// it looks transient.
		return isTransientError(err)
	}
}

// shouldNetworkRetry determines if a network error is retryable.
// Retries on timeouts, DNS failures, connection refused/reset,
// unexpected EOF, broker errors, and channel closed errors.
func shouldNetworkRetry(cfg *RetryConfig, attempt int, err error) bool {
	if !cfg.enabled() || attempt >= cfg.maxRetries() {
		return false
	}
	if !retryablePartialWrite(cfg, err) {
		return false
	}
	return isTransientError(err)
}

// transientErrorTexts are substrings that mark an error as transient when
// the typed checks in isTransientError don't apply, for errors that
// third-party libraries only expose as text: "timeout"/"timed out" covers
// our own publish timeouts and library errors that don't implement
// net.Error, "not connected" covers paho publishing while the broker
// connection is down, "channel/connection is not open" is amqp's.
//
// The last two are paho's connection-status errors, returned when a
// client is asked to connect while it is still leaving a previous state.
// paho's Disconnect() does its work in a goroutine and returns on its own
// timer, so a client we have given up on can still be mid-transition when
// the next attempt arrives - see internal/mqtt, which replaces the client
// to avoid it. These are the fallback for the window that leaves.
var transientErrorTexts = []string{
	"timeout",
	"timed out",
	"no such host",
	"connection refused",
	"connection reset",
	"reset by peer",
	"unexpected eof",
	"broker not available",
	"not connected",
	"channel closed",
	"connection closed",
	"channel/connection is not open",
	"can only transition to connecting",
	"already connected or reconnecting",
}

// isTransientError reports whether the error looks like a transient
// failure worth retrying. Typed checks (net.Error, syscall errnos,
// sentinel errors) are preferred; string matching on a lower-cased copy
// of the error text (libraries are inconsistent about casing, e.g. paho's
// "not Connected") is the fallback. Since %w-wrapping embeds the wrapped
// error's text, matching the top-level message covers the whole chain.
func isTransientError(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	// NXDOMAIN is usually a persistent config error, but it can also be
	// transient DNS trouble, so it is retried too.
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) && dnsErr.IsNotFound {
		return true
	}
	// EPIPE is what a write to a socket the peer already closed returns
	// once the reset has been seen - the usual outcome when a receiver
	// restarts under us. EHOSTUNREACH and ENETUNREACH cover a route that
	// went away. None of these report Timeout(), and none of them shows
	// up in transientErrorTexts, so they need naming here.
	if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, syscall.ECONNREFUSED) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.EPIPE) ||
		errors.Is(err, syscall.EHOSTUNREACH) ||
		errors.Is(err, syscall.ENETUNREACH) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	msg := strings.ToLower(err.Error())
	for _, text := range transientErrorTexts {
		if strings.Contains(msg, text) {
			return true
		}
	}
	return false
}

// --- Retry-After header parsing (RFC 7231) ---

// parseHTTPRetryAfter parses the Retry-After response header per RFC 7231.
// Returns the parsed duration or nil if absent.
// Handles both integer-seconds and HTTP-date formats.
func parseHTTPRetryAfter(resp *http.Response) *time.Duration {
	if resp == nil || resp.Header == nil {
		return nil
	}

	raw := resp.Header.Get("Retry-After")
	if raw == "" {
		return nil
	}

	// Try integer seconds first
	if seconds, err := strconv.ParseInt(raw, 10, 64); err == nil {
		d := time.Duration(seconds) * time.Second
		return &d
	}
	// Fallback to HTTP-date format
	if t, err := http.ParseTime(raw); err == nil {
		d := time.Until(t)
		if d > 0 {
			return &d
		}
	}
	return nil
}

// RetryableHTTPError wraps an HTTP error with the status code and any
// Retry-After hint from the response, so the retry loop can make
// per-status decisions on errors returned by sendBytes.
type RetryableHTTPError struct {
	Err        error
	StatusCode int
	retryAfter *time.Duration
}

func (r *RetryableHTTPError) Error() string {
	return r.Err.Error()
}

func (r *RetryableHTTPError) Unwrap() error {
	return r.Err
}

// newRetryableHTTPError wraps an HTTP response error, capturing the
// status code and Retry-After hint. The response itself is not retained,
// so the error does not pin the request/response buffers.
func newRetryableHTTPError(err error, resp *http.Response) *RetryableHTTPError {
	return &RetryableHTTPError{
		Err:        err,
		StatusCode: resp.StatusCode,
		retryAfter: parseHTTPRetryAfter(resp),
	}
}

// httpStatusCode extracts an HTTP status code from the error chain.
// Returns -1 if the chain carries none.
func httpStatusCode(err error) int {
	var re *RetryableHTTPError
	if errors.As(err, &re) {
		return re.StatusCode
	}
	return -1
}

// extractRetryAfterFromError returns the Retry-After hint captured in a
// RetryableHTTPError in the error chain, or nil.
func extractRetryAfterFromError(err error) *time.Duration {
	var re *RetryableHTTPError
	if errors.As(err, &re) {
		return re.retryAfter
	}
	return nil
}
