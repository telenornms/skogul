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
	"sync"

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

	once sync.Once
	mc   skmqtt.MQTT
}

// Send publishes the container in skogul JSON-encoded format on an MQTT
// topic.
func (handler *MQTT) Send(c *skogul.Container) error {
	handler.once.Do(func() {
		if handler.Topics == nil {
			handler.Topics = []string{"#"}
		}
		handler.mc.Init(handler.Broker, handler.Username, handler.Password, handler.ClientID)
		handler.mc.Connect()
	})
	b, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		mqttLog.WithError(err).Panic("Unable to marshal json for debug output")
		return err
	}
	for _, topic := range handler.Topics {
		handler.mc.Client.Publish(topic, 0, false, b)
	}
	return nil
}

// Verify makes sure required configuration options are set
func (handler *MQTT) Verify() error {
	if handler.Topics == nil {
		mqttLog.Warn("MQTT topic(s) not set, sending all messages to wildcard ('#')")
	}
	return nil
}
