/*
 * skogul, mqtt-receiver
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

package receiver

import (
	"time"

	"github.com/telenornms/skogul"
	skmqtt "github.com/telenornms/skogul/internal/mqtt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttLog = skogul.Logger("receiver", "mqtt")

/*
MQTT connects to a MQTT broker and listens for messages on a topic.
*/
type MQTT struct {
	Broker   string             `doc:"Address of broker to connect to." example:"[::1]:8888"`
	Topics   []string           `doc:"List of topics to subscribe to"`
	Handler  *skogul.HandlerRef `doc:"Handler used to parse, transform and send data."`
	Password string             `doc:"Username for authenticating to the broker."`
	Username string             `doc:"Password for authenticating."`
	ClientID string             `doc:"Custom client id to use (default: random)"`

	mc skmqtt.MQTT
}

func appendTopic(container *skogul.Container, topic string) {
	for _, metric := range container.Metrics {
		if metric.Metadata == nil {
			metric.Metadata = make(map[string]interface{})
		}
		metric.Metadata["_mqtt_topic"] = topic
	}
}

// Handle a received message.
func (handler *MQTT) receiver(msg mqtt.Message) {
	container, err := handler.Handler.H.Parse(msg.Payload())

	if err != nil {
		mqttLog.WithError(err).Error("Failed to parse payload from MQTT message")
		return
	}

	appendTopic(container, msg.Topic())

	err = handler.Handler.H.TransformAndSend(container)
	if err != nil {
		mqttLog.WithError(err).Error("Error during transform or send container")
	}
}

// Start MQTT receiver.
func (handler *MQTT) Start() error {
	handler.mc.Init(handler.Broker, handler.Username, handler.Password, handler.ClientID)
	for _, topic := range handler.Topics {
		handler.mc.Subscribe(topic, handler.receiver)
	}
	mqttLog.WithField("address", handler.Broker).Debug("Starting MQTT receiver")
	handler.mc.Connect()
	// Note that handler.listen() DOES return, because it only sets up
	// subscriptions. This sillyness is to satisfy the requirement that
	// Start() never returns. It should PROBABLY be more sensible.
	timer := time.NewTicker(10 * time.Second)
	for range timer.C {
	}
	return skogul.Error{Reason: "Shouldn't reach this"}
}

// Verify makes sure required configuration options are set
func (handler *MQTT) Verify() error {
	if handler.Broker == "" {
		return skogul.Error{Reason: "Missing address for MQTT receiver", Source: "MQTT receiver"}
	}
	if handler.Topics == nil {
		return skogul.Error{Reason: "MQTT topic(s) not set", Source: "MQTT receiver"}
	}
	return nil
}
