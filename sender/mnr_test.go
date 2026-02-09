/*
 * skogul, M&R sender tests
 *
 * Copyright (c) 2026 Telenor Norge AS
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

package sender_test

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/sender"
)

// startTCPListener starts a TCP listener on a random port and returns
// the listener and the received data after one connection.
func startTCPListener(t *testing.T) (net.Listener, int, chan string) {
	t.Helper()
	port := 10000 + rand.Intn(10000)
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatalf("Failed to start TCP listener: %v", err)
	}
	ch := make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 4096)
		n, _ := conn.Read(buf)
		ch <- string(buf[:n])
	}()
	return ln, port, ch
}

func TestMnR_Verify(t *testing.T) {
	tests := []struct {
		name    string
		address string
		action  string
		wantErr bool
	}{
		{"missing address", "", "", true},
		{"empty action is valid", "localhost:1234", "", false},
		{"refresh is valid", "localhost:1234", "refresh", false},
		{"delete is valid", "localhost:1234", "delete", false},
		{"Refresh is valid", "localhost:1234", "Refresh", false},
		{"DELETE is valid", "localhost:1234", "DELETE", false},
		{"invalid action", "localhost:1234", "bogus", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mnr := &sender.MnR{
				Address: tt.address,
				Action:  tt.action,
			}
			err := mnr.Verify()
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMnR_Send(t *testing.T) {
	ln, port, ch := startTCPListener(t)
	defer ln.Close()

	now := time.Now()
	m := skogul.Metric{
		Time:     &now,
		Metadata: map[string]any{"group": "testgroup", "prefix": "dev."},
		Data:     map[string]any{"cpu": 42},
	}
	c := skogul.Container{Metrics: []*skogul.Metric{&m}}

	mnr := &sender.MnR{
		Address: fmt.Sprintf("localhost:%d", port),
	}
	err := mnr.Send(&c)
	if err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}

	received := <-ch
	if !strings.Contains(received, "testgroup") {
		t.Errorf("Expected 'testgroup' in output, got: %s", received)
	}
	if !strings.Contains(received, "dev.cpu") {
		t.Errorf("Expected 'dev.cpu' in output, got: %s", received)
	}
	if !strings.HasPrefix(received, fmt.Sprintf("%d\t", now.Unix())) {
		t.Errorf("Expected line to start with timestamp, got: %s", received)
	}
}

func TestMnR_SendDefaultGroup(t *testing.T) {
	ln, port, ch := startTCPListener(t)
	defer ln.Close()

	now := time.Now()
	m := skogul.Metric{
		Time:     &now,
		Metadata: map[string]any{},
		Data:     map[string]any{"val": 1},
	}
	c := skogul.Container{Metrics: []*skogul.Metric{&m}}

	mnr := &sender.MnR{
		Address:      fmt.Sprintf("localhost:%d", port),
		DefaultGroup: "mygroup",
	}
	err := mnr.Send(&c)
	if err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}

	received := <-ch
	if !strings.Contains(received, "mygroup") {
		t.Errorf("Expected 'mygroup' in output, got: %s", received)
	}
}

func TestMnR_SendFallbackGroup(t *testing.T) {
	ln, port, ch := startTCPListener(t)
	defer ln.Close()

	now := time.Now()
	m := skogul.Metric{
		Time:     &now,
		Metadata: map[string]any{},
		Data:     map[string]any{"val": 1},
	}
	c := skogul.Container{Metrics: []*skogul.Metric{&m}}

	mnr := &sender.MnR{
		Address: fmt.Sprintf("localhost:%d", port),
	}
	err := mnr.Send(&c)
	if err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}

	received := <-ch
	if !strings.Contains(received, "\tgroup\t") {
		t.Errorf("Expected fallback 'group' in output, got: %s", received)
	}
}

func TestMnR_SendActionRefresh(t *testing.T) {
	ln, port, ch := startTCPListener(t)
	defer ln.Close()

	now := time.Now()
	m := skogul.Metric{
		Time:     &now,
		Metadata: map[string]any{},
		Data:     map[string]any{"val": 1},
	}
	c := skogul.Container{Metrics: []*skogul.Metric{&m}}

	mnr := &sender.MnR{
		Address: fmt.Sprintf("localhost:%d", port),
		Action:  "refresh",
	}
	err := mnr.Send(&c)
	if err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}

	received := <-ch
	if !strings.HasPrefix(received, "+r\t") {
		t.Errorf("Expected line to start with '+r\\t', got: %s", received)
	}
}

func TestMnR_SendActionDelete(t *testing.T) {
	ln, port, ch := startTCPListener(t)
	defer ln.Close()

	now := time.Now()
	m := skogul.Metric{
		Time:     &now,
		Metadata: map[string]any{},
		Data:     map[string]any{"val": 1},
	}
	c := skogul.Container{Metrics: []*skogul.Metric{&m}}

	mnr := &sender.MnR{
		Address: fmt.Sprintf("localhost:%d", port),
		Action:  "delete",
	}
	err := mnr.Send(&c)
	if err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}

	received := <-ch
	if !strings.HasPrefix(received, "+d\t") {
		t.Errorf("Expected line to start with '+d\\t', got: %s", received)
	}
}

func TestMnR_SendBadAddress(t *testing.T) {
	mnr := &sender.MnR{
		Address: "localhost:1",
	}
	c := skogul.Container{Metrics: []*skogul.Metric{}}
	err := mnr.Send(&c)
	if err == nil {
		t.Error("Expected error when sending to bad address, got nil")
	}
}
