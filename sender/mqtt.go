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

package sender

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/telenornms/skogul"
	skmqtt "github.com/telenornms/skogul/internal/mqtt"
)

/*
MQTT Sender publishes messages on a MQTT message bus.

FIXME: The MQTT-sender and receiver should be updated to not use the
url-encoded scheme.
*/
type MQTT struct {
	Address string `doc:"URL-encoded address." example:"mqtt://user:password@server/topic"`

	once sync.Once
	mc   skmqtt.MQTT
}

// Send publishes the container in skogul JSON-encoded format on an MQTT
// topic.
func (handler *MQTT) Send(c *skogul.Container) error {
	handler.once.Do(func() {
		handler.mc.Address = handler.Address
		handler.mc.Init()
		handler.mc.Connect()
	})
	b, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		log.Panicf("Unable to marshal json for debug output: %s", err)
		return err
	}
	handler.mc.Client.Publish(handler.mc.Topic, 0, false, b)
	return nil
}
