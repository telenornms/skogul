/*
 * skogul, common receover test code
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

package receiver_test

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/config"
	"github.com/telenornms/skogul/receiver"
	"github.com/telenornms/skogul/sender"
)

var validContainer = skogul.Container{}

func init() {

	now := time.Now()
	m := skogul.Metric{}
	m.Time = &now
	m.Metadata = make(map[string]interface{})
	m.Data = make(map[string]interface{})
	m.Metadata["foo"] = "bar"
	m.Data["tall"] = 5
	validContainer.Metrics = []*skogul.Metric{&m}
}

// UpdateHttpConfigWithUsablePorts updates the configuration with
// ports that are available for use. This can be useful in a CI
// environment where we don't care about the port number.
// Note: We might end up with a port which is written in the test
// as supposed to fail. ... left as an exercise for a future reader?
// Note2: We only support HTTP receivers.
func UpdateHttpConfigWithUsablePorts(config *config.Config) {
	for key := range config.Receivers {
		conf, ok := config.Receivers[key].Receiver.(*receiver.HTTP)
		if !ok {
			fmt.Println("failed to cast config to http")
			continue
		}
		addr := conf.Address
		origAddress := addr

		add := strings.Split(addr, ":")[0]
		port, _ := strconv.Atoi(strings.Split(addr, ":")[1])
		for i := 0; i < 10; i++ { // only try 10 times
			addr = fmt.Sprintf("%s:%d", add, port)
			if sock, err := net.Listen("tcp", addr); err != nil {
				// port might be in use, find another
				port++
			} else {
				sock.Close()
				break
			}

		}
		conf.Address = addr
		// fix senders
		for key, send := range config.Senders {
			if strings.ToLower(send.Type) != "http" {
				continue
			}

			sConf, ok := config.Senders[key].Sender.(*sender.HTTP)
			if !ok {
				fmt.Println("failed to cast config to http sender")
				continue
			}
			sConf.URL = strings.ReplaceAll(sConf.URL, origAddress, conf.Address)
		}
	}
}

// UpdateTcpConfigWithUsablePorts does the same as UpdateHttpConfigWithUsablePorts,
// only that it does it for TCP instead...
func UpdateTcpConfigWithUsablePorts(config *config.Config) {
	for key := range config.Receivers {
		conf, ok := config.Receivers[key].Receiver.(*receiver.TCPLine)
		if !ok {
			fmt.Println("failed to cast config to http")
			continue
		}
		addr := conf.Address
		origAddress := addr

		add := strings.Split(addr, ":")[0]
		port, _ := strconv.Atoi(strings.Split(addr, ":")[1])
		for i := 0; i < 10; i++ { // only try 10 times
			addr = fmt.Sprintf("%s:%d", add, port)
			if sock, err := net.Listen("tcp", addr); err != nil {
				// port might be in use, find another
				port++
			} else {
				sock.Close()
				break
			}

		}
		conf.Address = addr
		// fix senders
		for key, send := range config.Senders {
			if strings.ToLower(send.Type) != "net" {
				continue
			}

			sConf, ok := config.Senders[key].Sender.(*sender.Net)
			if !ok {
				fmt.Println("failed to cast config to http sender")
				continue
			}
			fmt.Println("updating", origAddress, conf.Address)
			sConf.Address = strings.ReplaceAll(sConf.Address, origAddress, conf.Address)
		}
	}
}
