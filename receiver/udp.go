/*
 * skogul, udp message receiver
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
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
	"net"
	"runtime"
	"sync"

	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/telenornms/skogul"
)

const (
	UDP_MAX_READ_SIZE = 65535
)

var udpLog = skogul.Logger("receiver", "udp")

// UDP contains the configuration for the receiver
type UDP struct {
	Address      string            `doc:"Address and port to listen to." example:"[::1]:3306"`
	Handler      skogul.HandlerRef `doc:"Handler used to parse, transform and send data."`
	Backlog      int               `doc:"Number of queued messages that are not delivered before the receiver starts blocking. Defaults to 100. Even when this is full, the kernel will still buffer data on top of this value. Higher values give smoother performance, but slightly more memory usage. Actual memory usage is a factor of average message size * backlog-fill."`
	Threads      int               `doc:"Number of worker go routines to use, which loosely translates to parallel execution. Defaults to number of CPU threads, with a minimum of 20. There is no correct number, but the value depends on how fast your senders are."`
	PacketSize   int               `doc:"UDP Packet size note: max. UDP read size is 65535"`
	FailureLevel string            `doc:"Level to log receiver failures as. Error, Warning, Info, Debug, or Trace. (default: Error)"`
	ch           chan []byte       // Used to pass messages from the accept/read-loop to the worker pool/threads.
	failureLevel logrus.Level
	once         sync.Once
}

// process is the worker-thread responsible for handling individual
// messages.
func (ud *UDP) process() {
	ud.once.Do(func() {
		if ud.FailureLevel == "" {
			ud.failureLevel = logrus.ErrorLevel
		} else {
			ud.failureLevel = skogul.GetLogLevelFromString(ud.FailureLevel)
		}
	})
	for {
		bytes := <-ud.ch
		if err := ud.Handler.H.Handle(bytes); err != nil {
			udpLog.WithError(err).Log(ud.failureLevel, "Unable to handle UDP message")
		}
	}
}

// Verify verifies the configuration for the UDP receiver
func (ud *UDP) Verify() error {
	if ud.PacketSize < 0 || ud.PacketSize > UDP_MAX_READ_SIZE {
		return skogul.Error{Source: "udp-receiver", Reason: "invalid udp packet size, maximum udp read size is between 0 and " + strconv.Itoa(UDP_MAX_READ_SIZE)}
	}
	return nil
}

// Start boots up ud.Threads number of worker threads, then starts
// listening for incoming UDP messages on the configured address. Start
// never returns.
func (ud *UDP) Start() error {
	if ud.PacketSize == 0 {
		ud.PacketSize = 9000
	}
	if ud.Backlog == 0 {
		ud.Backlog = 100
	}
	if ud.Threads == 0 {
		ud.Threads = runtime.NumCPU()
		if ud.Threads < 20 {
			ud.Threads = 20
		}
	}
	udpLog.Tracef("Got backlog size of %d and number of threads %d", ud.Backlog, ud.Threads)
	ud.ch = make(chan []byte, ud.Backlog)
	for i := 0; i < ud.Threads; i++ {
		go ud.process()
	}

	udpip, err := net.ResolveUDPAddr("udp", ud.Address)
	if err != nil {
		udpLog.WithError(err).WithField("address", ud.Address).Error("Can't resolve address")
		return err
	}
	ln, err := net.ListenUDP("udp", udpip)
	if err != nil {
		udpLog.WithError(err).WithField("address", ud.Address).Error("Can't listen on address")
		return err
	}
	for {
		bytes := make([]byte, ud.PacketSize)
		n, err := ln.Read(bytes)
		if err != nil || n == 0 {
			udpLog.WithError(err).WithField("bytes", n).Error("Unable to read UDP message")
			continue
		}
		ud.ch <- bytes[0:n]
	}
}
