/*
* skogul, mqtt common function tests
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

package mqtt

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// setTimeouts shrinks the connect timeouts for the duration of a test.
func setTimeouts(t *testing.T, timeout, margin time.Duration) {
	t.Helper()
	oldTimeout, oldMargin := connectTimeout, connectTimeoutMargin
	connectTimeout, connectTimeoutMargin = timeout, margin
	t.Cleanup(func() { connectTimeout, connectTimeoutMargin = oldTimeout, oldMargin })
}

// stalledBroker accepts TCP connections and then says nothing, so the
// MQTT handshake never completes. The returned function reports how many
// connections have been accepted, which is one per connect attempt.
func stalledBroker(t *testing.T) (string, func() int) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("unable to listen on loopback: %v", err)
	}
	t.Cleanup(func() { l.Close() })
	var accepted atomic.Int64
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			accepted.Add(1)
			// Hold the connection open without answering. Closed by
			// the client's deadline, or when the test ends.
			t.Cleanup(func() { conn.Close() })
		}
	}()
	return fmt.Sprintf("tcp://%s", l.Addr().String()), func() int { return int(accepted.Load()) }
}

// A connect that fails must leave the client in a state where the next
// connect starts over. paho refuses a Connect() on a client still stuck
// in its 'connecting' state ("status can only transition to connecting
// from disconnected", or "already connected or reconnecting" once it
// gets further), and neither reads as transient - so a sender retrying
// on top of this gives up on the second attempt instead of backing off.
// Every attempt must therefore report the underlying timeout.
func TestConnectFailureIsRepeatable(t *testing.T) {
	setTimeouts(t, 200*time.Millisecond, 5*time.Second)

	broker, _ := stalledBroker(t)
	handler := &MQTT{}
	if err := handler.Init(broker, "", "", "skogul-test"); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	for attempt := 1; attempt <= 3; attempt++ {
		err := handler.Connect()
		if err == nil {
			t.Fatalf("attempt %d: connect to a stalled broker unexpectedly succeeded", attempt)
		}
		msg := strings.ToLower(err.Error())
		if !strings.Contains(msg, "timeout") && !strings.Contains(msg, "timed out") {
			t.Fatalf("attempt %d: expected a timeout, got a status error from a client left mid-connect by the previous attempt: %v", attempt, err)
		}
	}
}

// stuckToken never completes, standing in for a paho connect that hangs
// past our own wait.
type stuckToken struct{}

func (stuckToken) Wait() bool                     { return false }
func (stuckToken) WaitTimeout(time.Duration) bool { return false }
func (stuckToken) Done() <-chan struct{}          { return make(chan struct{}) }
func (stuckToken) Error() error                   { return nil }

// stuckClient is a paho client whose Connect() never completes. It records
// whether Disconnect was called to abort it.
type stuckClient struct {
	paho.Client
	disconnected bool
}

func (c *stuckClient) IsConnected() bool      { return false }
func (c *stuckClient) IsConnectionOpen() bool { return false }
func (c *stuckClient) Connect() paho.Token    { return stuckToken{} }
func (c *stuckClient) Disconnect(uint)        { c.disconnected = true }

// When the wait does expire, the in-flight connect has to be aborted;
// leaving it running is what produces the "already connected or
// reconnecting" state tested above.
func TestConnectTimeoutAbortsAttempt(t *testing.T) {
	setTimeouts(t, 10*time.Millisecond, 10*time.Millisecond)

	handler := &MQTT{}
	if err := handler.Init("tcp://192.0.2.1:1883", "", "", "skogul-test"); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	client := &stuckClient{}
	handler.client = client

	err := handler.Connect()
	if err == nil {
		t.Fatal("expected a timeout error from a connect that never completes")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("expected a timeout error, got: %v", err)
	}
	if !client.disconnected {
		t.Error("the timed-out connect attempt was not aborted, so the next Connect() will fail with a paho status error instead of retrying")
	}
}

// The abort path is the one paho makes hard: Disconnect() does its work
// in a goroutine and returns when its quiesce timer fires, so a client
// aborted mid-connect is left in paho's "disconnecting" state and the
// next Connect() on it fails with a status error instead of connecting.
// Unlike TestConnectTimeoutAbortsAttempt above, this drives a real paho
// client, so it fails if we go back to relying on Disconnect alone.
//
// The negative margin makes our own wait (connectTimeout +
// connectTimeoutMargin) expire well before paho's ConnectTimeout, which
// is what puts us on the abort path at all.
func TestConnectAbortLeavesClientUsable(t *testing.T) {
	setTimeouts(t, 2*time.Second, -1900*time.Millisecond)

	broker, _ := stalledBroker(t)
	handler := &MQTT{}
	if err := handler.Init(broker, "", "", "skogul-test"); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	first := handler.Connect()
	if first == nil {
		t.Fatal("connect to a stalled broker unexpectedly succeeded")
	}
	if !strings.Contains(first.Error(), "timed out") {
		t.Fatalf("expected our own wait to expire first, got: %v", first)
	}

	// The aborted client is still finishing its connect in paho, so this
	// is the attempt that used to fail with a status error.
	second := handler.Connect()
	if second == nil {
		t.Fatal("connect to a stalled broker unexpectedly succeeded")
	}
	msg := strings.ToLower(second.Error())
	if !strings.Contains(msg, "timeout") && !strings.Contains(msg, "timed out") {
		t.Fatalf("second attempt hit a client left mid-connect by the first: %v", second)
	}
}

// Concurrent callers share one attempt. The MQTT sender connects on
// every send, so without this a pipeline fanning out over N goroutines
// serializes N attempts of up to connectTimeout each against a broker
// that is down.
func TestConcurrentConnectSharesOneAttempt(t *testing.T) {
	setTimeouts(t, 300*time.Millisecond, 5*time.Second)

	broker, dials := stalledBroker(t)
	handler := &MQTT{}
	if err := handler.Init(broker, "", "", "skogul-test"); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	// Pin the protocol version so one attempt means one dial: paho
	// otherwise dials a second time to fall back from 3.1.1 to 3.1 when
	// the broker does not answer. NewClient copies the options, so the
	// client has to be rebuilt after changing them.
	handler.opts.SetProtocolVersion(4)
	handler.newClient()

	const callers = 5
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := handler.Connect(); err == nil {
				t.Error("connect to a stalled broker unexpectedly succeeded")
			}
		}()
	}
	wg.Wait()

	if got := dials(); got != 1 {
		t.Errorf("%d concurrent connects made %d attempts, want 1 - they are not sharing one attempt", callers, got)
	}
}
