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
	"log"
	"net/url"
	"time"

	skmqtt "github.com/KristianLyng/skogul/internal/mqtt"
	skogul "github.com/KristianLyng/skogul/pkg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
MQTT connects to a MQTT broker and listens for messages on a topic.
*/
type MQTT struct {
	Address  string
	Handler  *skogul.Handler
	Password string
	Username string

	mc skmqtt.MQTT
}

// Handle a received message.
func (handler *MQTT) receiver(msg mqtt.Message) {
	m, err := handler.Handler.Parser.Parse(msg.Payload())
	if err == nil {
		err = m.Validate()
	}
	if err != nil {
		log.Printf("Unable to parse payload: %s", err)
		return
	}
	for _, t := range handler.Handler.Transformers {
		t.Transform(&m)
	}
	handler.Handler.Sender.Send(&m)
}

// Start MQTT receiver.
func (handler *MQTT) Start() error {
	handler.mc.Address = handler.Address
	handler.mc.Username = handler.Username
	handler.mc.Password = handler.Password
	handler.mc.Init()
	handler.mc.Subscribe(handler.mc.Topic, handler.receiver)
	log.Printf("Starting MQTT receiver at %s", handler.Address)
	handler.mc.Connect()
	// Note that handler.listen() DOES return, because it only sets up
	// subscriptions. This sillyness is to satisfy the requirement that
	// Start() never returns. It should PROBABLY be more sensible.
	timer := time.NewTicker(10 * time.Second)
	for range timer.C {
	}
	return skogul.Error{Reason: "Shouldn't reach this"}
}

func init() {
	addAutoReceiver("mqtt", NewMQTT, "Listen for Skogul-formatted JSON on a MQTT endpoint")
}

/*
NewMQTT returns a new MQTT receiver built from provided URL, using
the path as the topic to subscribe to.
*/
func NewMQTT(ul url.URL, h skogul.Handler) skogul.Receiver {
	n := MQTT{Address: ul.String(), Handler: &h}
	return &n
}
