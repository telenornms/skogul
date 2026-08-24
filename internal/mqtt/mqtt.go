/*
 * skogul, mqtt common functions
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.com>
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

/*
Package mqtt provides a bit of glue common between Skogul's MQTT sender and
receiver. Mostly providing mechanisms for setting up and maintaining a
connection to a broker. You really should not include this directly. Use the
MQTT sender and receiver instead.
*/
package mqtt

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/telenornms/skogul"
)

var mqttLog = skogul.Logger("sender", "mqtt")

// connectTimeout bounds how long a single CONNECT attempt may take.
// Connect is on the sender's send path, so waiting forever would pin the
// sending goroutine on a broker that accepts the TCP connection but never
// answers, instead of letting the sender's retry/backoff handle it. It is
// handed to paho as its own ConnectTimeout, which bounds both the dial and
// the handshake. A var, not a const, so tests can shrink it.
var connectTimeout = 30 * time.Second

// connectTimeoutMargin is how much longer than connectTimeout we wait on
// the connect token before giving up on our own. paho bounds its attempt
// by ConnectTimeout, so the margin means the usual outcome is a token
// carrying a real error, rather than our timeout racing paho's. A var
// for the same reason as connectTimeout - and it may be set negative in
// tests, to make our own wait expire first and exercise the abort path.
var connectTimeoutMargin = 5 * time.Second

// MQTT contains an MQTT client, its options and its configuration for handling messages
type MQTT struct {
	RenewClientID bool
	MQTTLogs      bool
	opts          *mqtt.ClientOptions
	topics        map[string]*MessageHandler

	// clientMu guards client, which is replaced whenever a connection
	// is abandoned or the client ID is renewed. Separate from the
	// connect bookkeeping below so that a publish never waits behind a
	// connect attempt.
	clientMu sync.RWMutex
	client   mqtt.Client

	// connectMu guards attempt, which holds the connect currently in
	// flight, if any, so that concurrent callers share one attempt
	// instead of queueing up to make their own.
	connectMu sync.Mutex
	attempt   *connectAttempt
}

// connectAttempt is a single connect in flight. done is closed once err
// holds the outcome, which every caller that joined the attempt returns.
type connectAttempt struct {
	done chan struct{}
	err  error
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

// Publish a payload on a topic, returning paho's token for it. Goes
// through the handler rather than the client directly, since the client
// is replaced whenever a connection is abandoned.
func (handler *MQTT) Publish(topic string, qos byte, retained bool, payload any) mqtt.Token {
	return handler.getClient().Publish(topic, qos, retained, payload)
}

// getClient returns the current client. It is replaced, not mutated, so
// the caller can keep using the one it got even if a reconnect swaps in
// a new one meanwhile.
func (handler *MQTT) getClient() mqtt.Client {
	handler.clientMu.RLock()
	defer handler.clientMu.RUnlock()
	return handler.client
}

// SetClient installs a client directly, bypassing Init. Exists for tests
// that drive a sender or receiver against a stand-in client rather than a
// broker.
func (handler *MQTT) SetClient(client mqtt.Client) {
	handler.clientMu.Lock()
	handler.client = client
	handler.clientMu.Unlock()
}

// newClient installs a fresh client built from the current options and
// returns it.
func (handler *MQTT) newClient() mqtt.Client {
	client := mqtt.NewClient(handler.opts)
	handler.SetClient(client)
	return client
}

// abandonClient gives up on a client and installs a fresh one, returning
// it. Disconnecting alone is not enough: paho does the work in a
// goroutine and returns as soon as the quiesce timer fires, while the
// transition itself waits for whatever the client is busy with - so a
// client we just disconnected can still be in paho's "disconnecting"
// state, and a Connect() on it fails outright with a status error
// ("status can only transition to connecting from disconnected") rather
// than trying. That does not look transient, so a sender retrying on top
// of it would give up on the next attempt instead of backing off. A
// fresh client always starts out disconnected; the old one finishes
// aborting in the background and is dropped.
func (handler *MQTT) abandonClient(client mqtt.Client, quiesce uint) mqtt.Client {
	client.Disconnect(quiesce)
	return handler.newClient()
}

// Connect to the broker and subscribe to the relevant topics, if any.
// Safe to call from multiple goroutines: a caller that finds the client
// already connected (because another goroutine just reconnected it)
// returns immediately, and callers that arrive while an attempt is in
// flight wait for it and share its result rather than starting another.
// The latter matters because the MQTT sender connects on every send: with
// the broker down, a pipeline fanning out over N goroutines would
// otherwise serialize N attempts of up to connectTimeout each.
func (handler *MQTT) Connect() error {
	if client := handler.getClient(); client != nil && client.IsConnected() {
		return nil
	}

	handler.connectMu.Lock()
	if attempt := handler.attempt; attempt != nil {
		handler.connectMu.Unlock()
		<-attempt.done
		return attempt.err
	}
	attempt := &connectAttempt{done: make(chan struct{})}
	handler.attempt = attempt
	handler.connectMu.Unlock()

	attempt.err = handler.connect()

	handler.connectMu.Lock()
	handler.attempt = nil
	handler.connectMu.Unlock()
	close(attempt.done)

	return attempt.err
}

// connect performs a single connect attempt. Only ever called by the
// goroutine that owns the current connectAttempt, so it has the client
// and the options to itself.
func (handler *MQTT) connect() error {
	client := handler.getClient()
	if client == nil {
		return fmt.Errorf("MQTT client not set up, Init has not been called")
	}
	if client.IsConnected() {
		return nil
	}
	if client.IsConnectionOpen() {
		mqttLog.Trace("Disconnecting client before (re)connecting")
		// Some quiesce time here, unlike the abort below: this
		// connection is alive, so it is worth giving paho a moment to
		// send the DISCONNECT rather than dropping the socket.
		client = handler.abandonClient(client, 100)
	}
	if handler.RenewClientID {
		clientID := fmt.Sprintf("skogul-%d-%d", rand.Uint32(), rand.Uint32())
		handler.opts.SetClientID(clientID)
		// renew client to set new clientID
		client = handler.newClient()
	}

	mqttLog.Debugf("Connecting to MQTT broker as '%s'", handler.opts.ClientID)
	token := client.Connect()
	wait := connectTimeout + connectTimeoutMargin
	if !token.WaitTimeout(wait) {
		// paho did not finish within its own ConnectTimeout either, so
		// the attempt is stuck. Abandon it rather than leaving it to
		// finish in the background: a client left mid-connect refuses a
		// fresh Connect(), and that refusal does not look transient, so
		// the sender's retry would give up on the next attempt instead
		// of backing off. See abandonClient.
		mqttLog.Errorf("Timed out after %v waiting to connect to MQTT broker, aborting the attempt", wait)
		handler.abandonClient(client, 0)
		return fmt.Errorf("connecting to MQTT broker timed out after %v", wait)
	}
	if err := token.Error(); err != nil {
		mqttLog.WithError(err).Error("Failed to connect to MQTT broker")
		return err
	}
	for i, messageHandler := range handler.topics {
		client.Subscribe(i, 0, func(_ mqtt.Client, msg mqtt.Message) { (*messageHandler)(msg) })
	}
	return nil
}

// Init sets up the MQTT client
func (handler *MQTT) Init(address, username, password, clientID string) error {
	if handler.MQTTLogs {
		mqtt.ERROR = mqttLog
		mqtt.CRITICAL = mqttLog
		mqtt.WARN = mqttLog
		mqtt.DEBUG = mqttLog
	}

	if err := handler.createClientOptions(address, username, password, clientID); err != nil {
		return err
	}
	handler.newClient()
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
	handler.opts.SetConnectTimeout(connectTimeout)
	handler.opts.SetPingTimeout(time.Duration(40 * time.Second))
	handler.opts.SetConnectionLostHandler(handler.connLostHandler)
	return nil
}
