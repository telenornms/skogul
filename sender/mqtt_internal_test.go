/*
 * skogul, mqtt sender tests
 *
 * Copyright (c) 2026 Telenor Norge AS
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
	"fmt"
	"sync"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/telenornms/skogul"
)

// publishToken is a completed paho token carrying a fixed outcome.
type publishToken struct {
	err      error
	timedOut bool
}

func (t publishToken) Wait() bool { return !t.timedOut }
func (t publishToken) WaitTimeout(time.Duration) bool {
	return !t.timedOut
}
func (t publishToken) Done() <-chan struct{} {
	c := make(chan struct{})
	close(c)
	return c
}
func (t publishToken) Error() error { return t.err }

// recordingClient stands in for a connected paho client and counts the
// publishes it sees per topic. failures maps a topic to the number of
// leading publishes that should fail before it starts succeeding.
type recordingClient struct {
	paho.Client

	mu        sync.Mutex
	failures  map[string]int
	timeouts  map[string]int
	published map[string]int
}

func newRecordingClient() *recordingClient {
	return &recordingClient{
		failures:  make(map[string]int),
		timeouts:  make(map[string]int),
		published: make(map[string]int),
	}
}

func (c *recordingClient) IsConnected() bool      { return true }
func (c *recordingClient) IsConnectionOpen() bool { return true }

func (c *recordingClient) Publish(topic string, _ byte, _ bool, _ any) paho.Token {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.published[topic]++
	if c.timeouts[topic] > 0 {
		c.timeouts[topic]--
		return publishToken{timedOut: true}
	}
	if c.failures[topic] > 0 {
		c.failures[topic]--
		return publishToken{err: fmt.Errorf("connection reset by peer")}
	}
	return publishToken{}
}

func (c *recordingClient) count(topic string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.published[topic]
}

// newTestMQTT builds an MQTT sender wired to client, skipping the real
// broker setup that Send would otherwise do on first use.
func newTestMQTT(topics []string, client paho.Client) *MQTT {
	handler := &MQTT{Topics: topics}
	// Consume the once so Send does not replace the client via mc.Init.
	handler.once.Do(func() {})
	handler.mc.SetClient(client)
	handler.BaseDelay = skogul.Duration{Duration: time.Millisecond}
	return handler
}

// A retry must re-publish only the topics that failed. Re-sending to a
// topic that already took the payload hands its subscribers a duplicate.
func TestMQTTRetriesOnlyFailedTopics(t *testing.T) {
	client := newRecordingClient()
	client.failures["b"] = 1
	handler := newTestMQTT([]string{"a", "b"}, client)

	if err := handler.Send(&skogul.Container{}); err != nil {
		t.Fatalf("Send() failed: %v", err)
	}
	if got := client.count("a"); got != 1 {
		t.Errorf("topic a was published %d times, want 1 - the retry duplicated a topic that had already succeeded", got)
	}
	if got := client.count("b"); got != 2 {
		t.Errorf("topic b was published %d times, want 2", got)
	}
}

// Same, for a publish that times out rather than erroring.
func TestMQTTRetriesOnlyTimedOutTopics(t *testing.T) {
	client := newRecordingClient()
	client.timeouts["b"] = 1
	handler := newTestMQTT([]string{"a", "b", "c"}, client)

	if err := handler.Send(&skogul.Container{}); err != nil {
		t.Fatalf("Send() failed: %v", err)
	}
	for _, topic := range []string{"a", "c"} {
		if got := client.count(topic); got != 1 {
			t.Errorf("topic %s was published %d times, want 1", topic, got)
		}
	}
	if got := client.count("b"); got != 2 {
		t.Errorf("topic b was published %d times, want 2", got)
	}
}

// A topic that never comes back must still fail the send, and must not
// drag the topics that succeeded into further attempts.
func TestMQTTExhaustedRetriesReportFailure(t *testing.T) {
	client := newRecordingClient()
	client.failures["b"] = 99
	handler := newTestMQTT([]string{"a", "b"}, client)

	err := handler.Send(&skogul.Container{})
	if err == nil {
		t.Fatal("expected Send() to fail when a topic never publishes")
	}
	if got := client.count("a"); got != 1 {
		t.Errorf("topic a was published %d times, want 1", got)
	}
	// Default MaxRetries is 3, so 4 attempts.
	if got := client.count("b"); got != 4 {
		t.Errorf("topic b was published %d times, want 4", got)
	}
}
