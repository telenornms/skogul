/*
 * skogul, mqtt sender
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

package senders

import (
	"encoding/json"
	"github.com/KristianLyng/skogul/pkg"
	skmqtt "github.com/KristianLyng/skogul/pkg/mqtt"
	"log"
	"net/url"
	"sync"
)

/*
MQTT Sender publishes messages on a MQTT message bus.
*/
type MQTT struct {
	Address string

	init bool
	mux  sync.Mutex
	mc   skmqtt.MQTT
}

func (handler *MQTT) testInit() {
	if handler.init {
		return
	}
	handler.mux.Lock()
	defer handler.mux.Unlock()
	if handler.init {
		return
	}
	handler.mc.Address = handler.Address
	handler.mc.Init()
	handler.mc.Connect()
	handler.init = true
}
func (handler *MQTT) Send(c *skogul.Container) error {
	handler.testInit()
	b, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		log.Panicf("Unable to marshal json for debug output: %s", err)
		return err
	}
	handler.mc.Client.Publish(handler.mc.Topic, 0, false, b)
	return nil
}
func init() {
	addAutoSender("mqtt", NewMQTT, "MQTT sender publishes received metrics to an MQTT broker/topic")
}

/*
NewMQTT creates a new MQTT sender
*/
func NewMQTT(url url.URL) skogul.Sender {
	x := MQTT{Address: url.String()}
	return &x
}
