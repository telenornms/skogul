/*
 * skogul, mqtt sender
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.com>
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
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/telenornms/skogul"
	skmqtt "github.com/telenornms/skogul/internal/mqtt"
)

var mqttLog = skogul.Logger("sender", "mqtt")

/*
MQTT Sender publishes messages on a MQTT message bus.

FIXME: The MQTT-sender and receiver should be updated to not use the
url-encoded scheme.
*/
type MQTT struct {
	Broker   string   `doc:"Address of broker to send to" example:"[::1]:8888"`
	Topics   []string `doc:"Topic(s) to publish events to"`
	Username string   `doc:"MQTT broker authorization username"`
	Password string   `doc:"MQTT broker authorization password"`
	ClientID string   `doc:"Custom client id to use (default: random)"`

	once    sync.Once
	mc      skmqtt.MQTT
	initErr error
	RetryConfig
}

// mqttPublishTimeout is how long we wait for a single publish to
// complete before treating it as failed. Note that we publish at QoS 0,
// where a publish completes once the packet has been written to the
// broker connection - it is not an acknowledgement from the broker, and
// a publish issued while the client is reconnecting completes without
// being written at all.
const mqttPublishTimeout = 10 * time.Second

// Send publishes the container in skogul JSON-encoded format on an MQTT
// topic.
func (handler *MQTT) Send(c *skogul.Container) error {
	handler.once.Do(func() {
		if handler.Topics == nil {
			handler.Topics = []string{"#"}
		}
		handler.initErr = handler.mc.Init(handler.Broker, handler.Username, handler.Password, handler.ClientID)
	})
	if handler.initErr != nil {
		return fmt.Errorf("MQTT sender initialization failed: %w", handler.initErr)
	}
	b, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		mqttLog.WithError(err).Panic("Unable to marshal json for debug output")
		return err
	}
	// Only the topics that failed are re-published, so a topic that was
	// already delivered does not get a second copy when a sibling topic
	// fails. Note the case this cannot cover: a publish that we timed out
	// on may still be written afterwards, since a QoS 0 token completing
	// is not an acknowledgement from the broker - so retrying a timed-out
	// topic can still duplicate.
	pending := handler.Topics
	return retryNetwork(&handler.RetryConfig, mqttLog, func() error {
		// Connect is a no-op if the client is already connected.
		if err := handler.mc.Connect(); err != nil {
			return err
		}
		// Issue all publishes before waiting on any of them, so they
		// proceed in parallel rather than one topic at a time. Each
		// wait still gets its own timeout, so a run of topics that
		// each stall can still add up.
		tokens := make([]pahomqtt.Token, len(pending))
		for i, topic := range pending {
			tokens[i] = handler.mc.Publish(topic, 0, false, b)
		}
		var failed []string
		var firstErr error
		for i, token := range tokens {
			var err error
			if !token.WaitTimeout(mqttPublishTimeout) {
				err = fmt.Errorf("MQTT publish to topic %s timed out", pending[i])
			} else if e := token.Error(); e != nil {
				err = fmt.Errorf("MQTT publish to topic %s: %w", pending[i], e)
			}
			if err != nil {
				failed = append(failed, pending[i])
				if firstErr == nil {
					firstErr = err
				}
			}
		}
		if firstErr == nil {
			return nil
		}
		total := len(pending)
		pending = failed
		if total == 1 {
			return firstErr
		}
		return fmt.Errorf("%d of %d topics failed, first error: %w", len(failed), total, firstErr)
	})
}

// Verify makes sure required configuration options are set
func (handler *MQTT) Verify() error {
	if handler.Broker == "" {
		return skogul.MissingArgument("Broker")
	}
	if handler.Topics == nil {
		mqttLog.Warn("MQTT topic(s) not set, sending all messages to wildcard ('#')")
	}
	return handler.verifyRetry()
}
