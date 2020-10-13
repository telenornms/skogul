/*
 * skogul, mqtt common functions
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

/*
Package mqtt provides a bit of glue common between Skogul's MQTT sender and
receiver. Mostly providing mechanisms for setting up and maintaining a
connection to a broker. You really should not include this directly. Use
the MQTT sender and receiver instead.
*/
package mqtt

import (
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/telenornms/skogul"
)

var mqttLog = skogul.Logger("sender", "mqtt")

// MQTT contains an MQTT client, its options and its configuration for handling messages
type MQTT struct {
	Client mqtt.Client
	opts   *mqtt.ClientOptions
	topics map[string]*MessageHandler
}

// MessageHandler is used to establish a callback when a message is
// received.
type MessageHandler func(Message mqtt.Message)

// Subscribe to a topic. callback is called whenever a message is received.
// This also deals with re-subscribing when a reconnect takes place.
func (handler *MQTT) Subscribe(topic string, callback MessageHandler) {
	mqttLog.WithField("topic", topic).Debug("MQTT subscribed")
	if handler.topics == nil {
		handler.topics = make(map[string]*MessageHandler)
	}
	handler.topics[topic] = &callback
}

// Connect to the broker and subscribe to the relevant topics, if any.
func (handler *MQTT) Connect() error {
	token := handler.Client.Connect()
	// Should probably be configurable, or at least not infinite.
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		mqttLog.WithError(err).Error("Failed to connect to MQTT broker")
		return err
	}
	for i, messageHandler := range handler.topics {
		handler.Client.Subscribe(i, 0, func(_ mqtt.Client, msg mqtt.Message) { (*messageHandler)(msg) })
	}
	return nil
}

// Init sets up the MQTT client
func (handler *MQTT) Init(address, username, password, clientID string) error {
	handler.createClientOptions(address, username, password, clientID)
	handler.Client = mqtt.NewClient(handler.opts)
	return nil
}

// connLostHandler handles reconnects if the connection drops.
func (handler *MQTT) connLostHandler(client mqtt.Client, e error) {
	mqttLog.WithError(e).Debug("Connection lost... Auto-reconnecting and re-subscribing.")
	for {
		e := handler.Connect()
		if e != nil {
			mqttLog.WithError(e).Debug("Failed to re-connect to MQTT broker. Retrying in 5 seconds")
			time.Sleep(time.Duration(5 * time.Second))
		} else {
			mqttLog.Debug("Reconnected to MQTT broker successfully.")
			break
		}
	}
}

// createClientOptions configures the MQTT client options
func (handler *MQTT) createClientOptions(address, username, password, clientID string) error {
	handler.opts = mqtt.NewClientOptions()
	handler.opts.AddBroker(address)
	if username != "" {
		handler.opts.SetUsername(username)
	}
	if password != "" {
		handler.opts.SetPassword(password)
	}
	if clientID == "" {
		clientID = fmt.Sprintf("skogul-%d-%d", rand.Uint32(), rand.Uint32())
	}
	handler.opts.SetClientID(clientID)
	handler.opts.SetAutoReconnect(false)
	handler.opts.SetPingTimeout(time.Duration(40 * time.Second))
	handler.opts.SetConnectionLostHandler(handler.connLostHandler)
	return nil
}
