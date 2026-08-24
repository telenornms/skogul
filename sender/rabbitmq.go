/*
 * skogul, rabbitmq producer/sender
 *
 * Copyright (c) 2023 Telenor Norge AS
 * Author(s):
 *  - Kamil Oracz <kamil.oracz@telenor.no>
 *
 * This library is free software; you can redistribute it and/or modify it
* under the terms of the GNU Lesser General Public License as published by the
* Free Software Foundation; either version 2.1 of the License, or (at your
* option) any later version.
 *
* This library is distributed in the hope that it will be useful, but WITHOUT
* ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
* FOR A PARTICULAR PURPOSE.  See the GNU Lesser General Public License for more
* details.
 *
* You should have received a copy of the GNU Lesser General Public License
* along with this library; if not, write to the Free Software Foundation, Inc.,
* 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301  USA
*/

package sender

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/encoder"
)

type Rabbitmq struct {
	Username skogul.Secret     `doc:"Username for rabbitmq instance"`
	Password skogul.Secret     `doc:"Password for rabbitmq instance"`
	Host     string            `doc:"Hostname for rabbitmq instance. Fallback is localhost"`
	Port     string            `doc:"Port for rabbitmq instance. Fallback is 5672"`
	Queue    string            `doc:"Queue to write to"`
	Encoder  skogul.EncoderRef `doc:"Encoder to use. Fallback is json"`
	Timeout  int               `doc:"Timeout for rabbitmq instance connection. Fallback is 10 seconds."`
	conn     *amqp.Connection
	channel  *amqp.Channel
	mu       sync.RWMutex
	once     sync.Once
	RetryConfig
}

var rabbitmqLog = skogul.Logger("sender", "rabbitmq")

func (r *Rabbitmq) init() {
	if r.Port == "" {
		r.Port = "5672"
	}

	if r.Host == "" {
		r.Host = "localhost"
	}

	if r.Timeout == 0 {
		r.Timeout = 10
	}

	if r.Encoder.E == nil {
		r.Encoder.E = encoder.JSON{}
	}
}

// connect establishes the broker connection and channel, and declares the
// queue. On success r.conn and r.channel are set. Callers must hold the
// write lock.
func (r *Rabbitmq) connect() error {
	timeout := time.Duration(r.Timeout) * time.Second
	conn, err := amqp.DialConfig(
		fmt.Sprintf("amqp://%s:%s@%s:%s/", r.Username.Expose(), r.Password.Expose(), r.Host, r.Port),
		amqp.Config{
			// Bounds both the TCP connect and the AMQP
			// handshake. Without it the library's own 30s
			// default applies and Timeout does nothing, since
			// PublishWithContext ignores its context.
			Dial: amqp.DefaultDial(timeout),
			// Same as what amqp.Dial() would have used.
			Locale: "en_US",
		},
	)
	if err != nil {
		rabbitmqLog.WithError(err).Error("Failed initializing broker connection")
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		rabbitmqLog.WithError(err).Error("Failed initializing channel")
		conn.Close()
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
		rabbitmqLog.WithError(err).Error("Failed to declare a queue")
		conn.Close()
		return err
	}

	r.conn = conn
	r.channel = ch
	return nil
}

// disconnect drops the current connection and channel so the next attempt
// reconnects from scratch. Callers must hold the write lock.
func (r *Rabbitmq) disconnect() {
	if r.conn != nil {
		r.conn.Close()
	}
	r.conn = nil
	r.channel = nil
}

// getChannel returns the channel to publish on, connecting first if there
// is none. The returned channel is safe to use without holding the lock:
// amqp091 serializes publishing internally, and a channel belonging to a
// connection another goroutine has since dropped just fails the publish,
// which is retried.
func (r *Rabbitmq) getChannel() (*amqp.Channel, error) {
	r.mu.RLock()
	ch := r.channel
	r.mu.RUnlock()
	if ch != nil {
		return ch, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	// Another goroutine may have connected while we waited for the
	// write lock.
	if r.channel == nil {
		if err := r.connect(); err != nil {
			return nil, err
		}
	}
	return r.channel, nil
}

// dropChannel tears down the connection ch belongs to, unless another
// goroutine has already replaced it with a fresh one.
func (r *Rabbitmq) dropChannel(ch *amqp.Channel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.channel == ch {
		r.disconnect()
	}
}

func (r *Rabbitmq) Send(c *skogul.Container) error {
	r.once.Do(r.init)

	body, err := r.Encoder.E.Encode(c)
	if err != nil {
		return err
	}

	return retryNetwork(&r.RetryConfig, rabbitmqLog, func() error {
		// Send can be called from multiple goroutines. Only the
		// connect/disconnect bookkeeping is serialized: the publish
		// itself takes amqp091's own channel lock, so holding r.mu
		// across it would let one stuck publish block every other
		// sender goroutine. Backoff sleeps happen in retryNetwork,
		// outside the lock either way.
		ch, err := r.getChannel()
		if err != nil {
			return fmt.Errorf("no active rabbitmq connections: %w", err)
		}

		// PublishWithContext ignores this context in amqp091-go
		// v1.10 - the actual bound on a stuck publish is the
		// connection heartbeat. Passed anyway in case that changes.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.Timeout)*time.Second)
		defer cancel()

		err = ch.PublishWithContext(
			ctx,
			"",
			r.Queue,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			},
		)
		if err != nil {
			r.dropChannel(ch)
			return err
		}

		return nil
	})
}

func (r *Rabbitmq) Verify() error {
	if r.Username.Expose() == "" {
		return skogul.MissingArgument("Username")
	}

	if r.Password.Expose() == "" {
		return skogul.MissingArgument("Password")
	}

	if r.Queue == "" {
		return skogul.MissingArgument("Queue")
	}

	return r.verifyRetry()
}
