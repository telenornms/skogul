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
	"time"

	"github.com/KristianLyng/skogul"
	skmqtt "github.com/KristianLyng/skogul/internal/mqtt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
MQTT connects to a MQTT broker and listens for messages on a topic.
*/
type MQTT struct {
	Address  string             `doc:"Address to connect to."`
	Handler  *skogul.HandlerRef `doc:"Handler used to parse, transform and send data."`
	Password string             `doc:"Username for authenticating to the broker."`
	Username string             `doc:"Password for authenticating."`

	mc skmqtt.MQTT
}

// Handle a received message.
func (handler *MQTT) receiver(msg mqtt.Message) {
	m, err := handler.Handler.H.Parser.Parse(msg.Payload())
	if err == nil {
		err = m.Validate()
	}
	if err != nil {
		log.Printf("Unable to parse payload: %s", err)
		return
	}
	for _, t := range handler.Handler.H.Transformers {
		if e := t.Transform(&m); e != nil {
			log.Printf("failed to transform mqtt message: %v", e)
			return
		}
	}
	handler.Handler.H.Sender.Send(&m)
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
