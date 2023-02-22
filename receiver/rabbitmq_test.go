/*
 * skogul, rabbitmq-receiver test
 *
 * Copyright (c) 2023 Telenor Norge AS
 * Author(s):
 *  - Kamil Oracz <kamil.oracz@telenor.no>
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
	"testing"

	"github.com/telenornms/skogul/config"
)

func TestRabbitmq(t *testing.T) {
	sconf := fmt.Sprintf(`
		 {
			 "receivers": {
					 "x": {
						 "type": "rabbitmq",
						 "handler": "kek",
						 "username":"guest",
						 "password":"guest",
						 "queue":"test-queue"
					 }
			 },
			 "handlers": {
					 "kek": {
							 "parser": "skogulmetric",
							 "transformers": [
									 "now"
							 ],
							 "sender": "test"
					 }
			 },
			 "senders": {
				 "test": {
					 "type": "test"
				 }
			 }
	 }`)

	conf, err := config.Bytes([]byte(sconf))

	if err != nil {
		t.Errorf("Failed to load config: %v", err)
		return
	}

	rcv := conf.Receivers["x"].Receiver.(*Rabbitmq)

	err = rcv.Start()

	if err != nil {
		t.Error(err)
	}
}
