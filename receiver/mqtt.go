/*
 * skogul, mqtt-receiver
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
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

package receiver

import (
	"fmt"
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
	Broker          string             `doc:"Address of broker to connect to." example:"[::1]:8888"`
	Topics          []string           `doc:"List of topics to subscribe to"`
	Handler         *skogul.HandlerRef `doc:"Handler used to parse, transform and send data."`
	Password        string             `doc:"Username for authenticating to the broker."`
	Username        string             `doc:"Password for authenticating."`
	ClientID        string             `doc:"Custom client id to use (default: random)"`
	RenewClientID   bool               `doc:"Renew the client ID on reconnects ([MQTT-3.1.4-2] @ https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc384800405)"`
	DisplayMQTTLogs bool

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

// reconnectInterval is how often Start rechecks that the receiver still
// has a broker connection. Nothing else retries a connect that never
// succeeded in the first place: the client runs without auto-reconnect,
// and its connection-lost handler only fires once a connection has been
// established. Same cadence as that handler uses.
const reconnectInterval = 5 * time.Second

// Start MQTT receiver.
func (handler *MQTT) Start() error {
	handler.mc.MQTTLogs = handler.DisplayMQTTLogs
	handler.mc.RenewClientID = handler.RenewClientID
	if err := handler.mc.Init(handler.Broker, handler.Username, handler.Password, handler.ClientID); err != nil {
		return fmt.Errorf("unable to set up MQTT client: %w", err)
	}
	for _, topic := range handler.Topics {
		handler.mc.Subscribe(topic, handler.receiver)
	}
	mqttLog.WithField("address", handler.Broker).Debug("Starting MQTT receiver")
	// Connect() establishes the subscriptions and returns immediately
	// once the client is connected, so this doubles as the supervisor
	// that satisfies the requirement that Start() never returns. It
	// covers the two cases the connection-lost handler never sees: a
	// broker that is down at startup, and a connect that timed out and
	// was aborted.
	timer := time.NewTicker(reconnectInterval)
	defer timer.Stop()
	for {
		if err := handler.mc.Connect(); err != nil {
			mqttLog.WithError(err).Errorf("Failed to connect to MQTT broker, retrying in %v", reconnectInterval)
		}
		<-timer.C
	}
}

// Verify makes sure required configuration options are set
func (handler *MQTT) Verify() error {
	if handler.Broker == "" {
		return skogul.MissingArgument("Broker")
	}
	if handler.Topics == nil {
		return skogul.MissingArgument("Topics")
	}
	if handler.RenewClientID && handler.ClientID != "" {
		mqttLog.Warning("RenewClientID AND ClientID is set - ClientID will change!")
	}
	return nil
}
