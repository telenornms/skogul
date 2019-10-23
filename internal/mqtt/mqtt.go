/*
 * skogul, mqtt common functions
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
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// MQTT shared data structure annd interal state
type MQTT struct {
	Address  string
	Client   mqtt.Client
	Topic    string
	Username string
	Password string
	opts     *mqtt.ClientOptions
	topics   map[string]*MessageHandler
	uri      *url.URL
	clientID string
}

// MessageHandler is used to establish a callback when a message is
// received.
type MessageHandler func(Message mqtt.Message)

// Subscribe to a topic. callback is called whenever a message is received.
// This also deals with re-subscribing when a reconnect takes place.
func (handler *MQTT) Subscribe(topic string, callback MessageHandler) {
	log.WithField("topic", topic).Debug("MQTT subscribed")
	if handler.topics == nil {
		handler.topics = make(map[string]*MessageHandler)
	}
	handler.topics[topic] = &callback
}

// Shim-layer that accepts a message and calls the appropriate callback.
func (handler *MQTT) receiver(client mqtt.Client, msg mqtt.Message) {
	t := msg.Topic()
	if handler.topics[t] == nil {
		log.WithFields(log.Fields{
			"topic": t,
			"msg":   msg,
		}).Debug("Message received on unknown topic")
		return
	}
	(*handler.topics[t])(msg)
}

// Connect to the broker and subscribe to the relevant topics, if any.
func (handler *MQTT) Connect() error {
	token := handler.Client.Connect()
	// Should probably be configurable, or at least not infinite.
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.WithError(err).Error("Failed to connect to MQTT broker")
		return err
	}
	for i := range handler.topics {
		handler.Client.Subscribe(i, 0, handler.receiver)
	}
	return nil
}

// Init the generic MQTT data structures, mostly parsing MQTT.Address.
func (handler *MQTT) Init() error {
	var err error
	handler.uri, err = url.Parse(handler.Address)
	if err != nil {
		log.WithError(err).Fatal("Initialization of MQTT failed")
	}
	handler.Topic = handler.uri.Path[1:len(handler.uri.Path)]
	if handler.Topic == "" {
		handler.Topic = "#"
	}
	handler.createClientOptions()
	handler.Client = mqtt.NewClient(handler.opts)
	return nil
}

// Handle reconnects when the connection drops.
func (handler *MQTT) connLostHandler(client mqtt.Client, e error) {
	log.WithError(e).Debug("Connection lost... Auto-reconnecting and re-subscribing.")
	for {
		e := handler.Connect()
		if e != nil {
			log.WithError(e).Debug("Failed to re-connect to MQTT broker. Retrying in 5 seconds")
			time.Sleep(time.Duration(5 * time.Second))
		} else {
			log.Debug("Reconnected to MQTT broker successfully.")
			break
		}
	}
}

// createClientOptions() sets up our default options.
func (handler *MQTT) createClientOptions() error {
	handler.opts = mqtt.NewClientOptions()
	handler.opts.AddBroker(fmt.Sprintf("tcp://%s", handler.uri.Host))
	if handler.Username != "" {
		handler.opts.SetUsername(handler.Username)
	} else {
		handler.opts.SetUsername(handler.uri.User.Username())
	}
	if handler.Password != "" {
		handler.opts.SetPassword(handler.Password)
	} else {
		password, _ := handler.uri.User.Password()
		handler.opts.SetPassword(password)
	}
	if handler.clientID == "" {
		handler.clientID = fmt.Sprintf("skogul-%d-%d", rand.Uint32(), rand.Uint32())
	}
	handler.opts.SetClientID(handler.clientID)
	handler.opts.SetAutoReconnect(false)
	handler.opts.SetPingTimeout(time.Duration(40 * time.Second))
	handler.opts.SetConnectionLostHandler(handler.connLostHandler)
	return nil
}
