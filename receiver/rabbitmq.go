/*
 * skogul, rabbitmq-receiver
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

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/telenornms/skogul"
)

type Rabbitmq struct {
	Username string             `doc:"Username for rabbitmq instance"`
	Password string             `doc:"Password for rabbitmq instance"`
	Host     string             `doc:"Hostname for rabbitmq instance. Fallback is localhost"`
	Port     string             `doc:"Port for rabbitmq instance. Fallback is 5672"`
	Queue    string             `doc:"Queue to read from"`
	Handler  *skogul.HandlerRef `doc:"Handler used to parse, transform and send data. Default skogul."`
}

func (r *Rabbitmq) Start() error {
	if r.Username == "" || r.Password == "" {
		fmt.Print("Error missing username or password")
	}

	if r.Port == "" {
		r.Port = "5672"
	}

	if r.Host == "" {
		r.Host = "localhost"
	}

	if r.Handler == nil {
		r.Handler = &skogul.HandlerRef{}
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", r.Username, r.Password, r.Host, r.Port))
	if err != nil {
		return err
	}

	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		r.Queue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		r.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	for message := range msgs {
		container, err := r.Handler.H.Parse(message.Body)

		if err != nil {
			// fmt.Errorf("Error failed to parse body %v", err)
			return err
		}

		err = r.Handler.H.TransformAndSend(container)
		if err != nil {
			// fmt.Errorf("Error transforming %v", err)
			return err
		}
	}

	return nil
}
