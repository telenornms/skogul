/*
* skogul, snmp sender tests
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
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/telenornms/skogul"
)

// trapListener is a UDP socket standing in for a trap receiver.
type trapListener struct {
	conn *net.UDPConn
	port uint16
}

func newTrapListener(t *testing.T) *trapListener {
	t.Helper()
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("unable to resolve loopback: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Skipf("unable to listen for traps: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return &trapListener{conn: conn, port: uint16(conn.LocalAddr().(*net.UDPAddr).Port)}
}

// receive counts the trap packets that arrive within a short window.
func (l *trapListener) receive(t *testing.T, want int) int {
	t.Helper()
	got := 0
	buf := make([]byte, 65535)
	for got < want {
		if err := l.conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			t.Fatalf("unable to set read deadline: %v", err)
		}
		if _, _, err := l.conn.ReadFrom(buf); err != nil {
			break
		}
		got++
	}
	return got
}

// snmpRetryConfig is a retry config with delays short enough not to slow
// the tests down. Kept local so this file stands on its own.
func snmpRetryConfig(maxRetries int) RetryConfig {
	enabled := true
	return RetryConfig{
		BackoffEnabled:    &enabled,
		MaxRetries:        &maxRetries,
		BaseDelay:         skogul.Duration{Duration: time.Millisecond},
		MaxDelay:          skogul.Duration{Duration: 10 * time.Millisecond},
		BackoffMultiplier: 2.0,
	}
}

func testSNMP(port uint16) *SNMP {
	return &SNMP{
		Port:      port,
		Community: "public",
		Version:   "2c",
		Target:    "127.0.0.1",
		Oidmap: map[string]any{
			"value": ".1.3.6.1.4.1.9.9.1.1",
			"name":  ".1.3.6.1.4.1.9.9.1.2",
		},
		SnmpTrapOID: ".1.3.6.1.4.1.9.9.1.0",
		RetryConfig: snmpRetryConfig(2),
	}
}

func testContainer(metrics int) *skogul.Container {
	c := skogul.Container{}
	for i := range metrics {
		c.Metrics = append(c.Metrics, &skogul.Metric{
			Data: map[string]any{"value": float64(i), "name": "test"},
		})
	}
	return &c
}

func TestSNMPVerify(t *testing.T) {
	x := testSNMP(1162)
	if err := x.Verify(); err != nil {
		t.Errorf("a valid config should verify, got: %v", err)
	}

	x.Target = ""
	if err := x.Verify(); err == nil {
		t.Error("a missing Target should fail verification")
	}

	x = testSNMP(1162)
	x.BackoffMultiplier = -1
	if err := x.Verify(); err == nil {
		t.Error("SNMP Verify should validate the retry config too")
	}
}

// The connection is established on demand rather than once at startup, so
// that a target which was unreachable earlier is picked up later. Sending
// after an explicit disconnect must reconnect rather than fail or panic.
func TestSNMPSendReconnects(t *testing.T) {
	l := newTrapListener(t)
	x := testSNMP(l.port)

	if err := x.Send(testContainer(2)); err != nil {
		t.Fatalf("unable to send traps: %v", err)
	}
	if got := l.receive(t, 2); got != 2 {
		t.Errorf("expected 2 traps, got %d", got)
	}

	// Drop the socket the way a failed send does, then send again.
	x.mu.Lock()
	x.disconnect()
	x.mu.Unlock()
	if x.g != nil {
		t.Fatal("disconnect should clear the handle")
	}

	if err := x.Send(testContainer(1)); err != nil {
		t.Fatalf("send after disconnect should reconnect, got: %v", err)
	}
	if got := l.receive(t, 1); got != 1 {
		t.Errorf("expected 1 trap after reconnect, got %d", got)
	}
}

// A connect failure must be reported and must not be permanent: the
// handle is left unset so the next send tries again. It used to be
// established once behind a sync.Once, with the error dropped on the
// floor, so a target that was down at startup stayed broken forever.
func TestSNMPConnectFailureIsNotPermanent(t *testing.T) {
	x := testSNMP(162)
	// An address literal that cannot be dialled, no DNS involved.
	x.Target = "300.300.300.300"

	err := x.Send(testContainer(1))
	if err == nil {
		t.Fatal("expected an error for an undialable target")
	}
	if !strings.Contains(err.Error(), "unable to connect") {
		t.Errorf("the connect error should be surfaced, got: %v", err)
	}
	if x.g != nil {
		t.Error("a failed connect must not leave a handle behind")
	}

	// Pointing at a real listener, the very next send must work.
	l := newTrapListener(t)
	x.Target = "127.0.0.1"
	x.Port = l.port
	if err := x.Send(testContainer(1)); err != nil {
		t.Fatalf("send should recover once the target is reachable, got: %v", err)
	}
	if got := l.receive(t, 1); got != 1 {
		t.Errorf("expected 1 trap, got %d", got)
	}
}

// gosnmp.GoSNMP is not safe for concurrent use, and a sender is shared by
// every handler pointing at it. Run under -race.
func TestSNMPConcurrentSend(t *testing.T) {
	l := newTrapListener(t)
	x := testSNMP(l.port)

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 5 {
				if err := x.Send(testContainer(1)); err != nil {
					t.Errorf("concurrent send failed: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	if got := l.receive(t, 40); got != 40 {
		t.Errorf("expected 40 traps, got %d", got)
	}
}

// A metric with nothing mappable produces a trap with no PDUs, which
// gosnmp refuses to send ("requires at least 1 PDU"). Skipping it keeps
// the rest of the container going, and a container of nothing but such
// metrics is not an error - the fields were dropped, same as buildTrap
// drops individual unmappable ones.
func TestSNMPSendSkipsEmptyTraps(t *testing.T) {
	l := newTrapListener(t)
	x := testSNMP(l.port)
	x.SnmpTrapOID = ""

	c := skogul.Container{
		Metrics: []*skogul.Metric{
			{Data: map[string]any{"nosuchfield": "a"}},
			{Data: map[string]any{"value": float64(1)}},
		},
	}
	if err := x.Send(&c); err != nil {
		t.Fatalf("a container with one usable metric should send, got: %v", err)
	}
	if got := l.receive(t, 1); got != 1 {
		t.Errorf("expected 1 trap, got %d", got)
	}

	empty := skogul.Container{
		Metrics: []*skogul.Metric{
			{Data: map[string]any{"nosuchfield": "a"}},
		},
	}
	if err := x.Send(&empty); err != nil {
		t.Errorf("a container with nothing to send should not error, got: %v", err)
	}
}

func TestSNMPBuildTrap(t *testing.T) {
	x := testSNMP(1162)

	tests := []struct {
		name  string
		data  map[string]any
		want  int // PDUs, including the trap OID
		named []string
	}{
		{
			name:  "mapped fields",
			data:  map[string]any{"value": float64(1), "name": "a"},
			want:  3,
			named: []string{".1.3.6.1.4.1.9.9.1.1", ".1.3.6.1.4.1.9.9.1.2"},
		},
		{
			// Used to become a PDU named "%!s(<nil>)", which only
			// fails later, inside gosnmp.
			name: "unmapped field is skipped",
			data: map[string]any{"value": float64(1), "nosuchfield": "a"},
			want: 2,
		},
		{
			// Used to append a nameless, valueless PDU.
			name: "unsupported type is skipped",
			data: map[string]any{"value": float64(1), "name": []string{"a"}},
			want: 2,
		},
		{
			name: "no usable data leaves just the trap OID",
			data: map[string]any{"nosuchfield": "a"},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trap := x.buildTrap(&skogul.Metric{Data: tt.data})
			if len(trap.Variables) != tt.want {
				t.Errorf("expected %d PDUs, got %d: %+v", tt.want, len(trap.Variables), trap.Variables)
			}
			for _, pdu := range trap.Variables {
				if pdu.Name == "" {
					t.Errorf("PDU with no name in trap: %+v", pdu)
				}
			}
			for _, name := range tt.named {
				found := false
				for _, pdu := range trap.Variables {
					if pdu.Name == name {
						found = true
					}
				}
				if !found {
					t.Errorf("expected a PDU named %s", name)
				}
			}
		})
	}
}
