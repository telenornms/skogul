/*
 * skogul, udp receiver tests
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
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/config"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"net"
	"os"
	"testing"
	"time"
)

// FIXME: This file is very fond of global scope and init(), mainly for the
// sake of benchmarks... and it's a bit messy.
var uConfig *config.Config
var pFile []byte
var u1 *net.UDPConn
var u2 *net.UDPConn
var pJSON = []byte("{\"metrics\":[{\"timestamp\":\"2019-03-15T11:08:02+01:00\",\"metadata\":{\"key\":\"value\"},\"data\":{\"string\":\"text\",\"float\":1.11,\"integer\":5}}]}")

func readProtobufFile(file string) []byte {
	b := make([]byte, 9000)
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("unable to open protobuf packet file: %v", err)
		os.Exit(1)
	}
	defer f.Close()
	n, err := f.Read(b)
	if err != nil {
		fmt.Printf("unable to read protobuf packet file: %v", err)
		os.Exit(1)
	}
	if n == 0 {
		fmt.Printf("read 0 bytes from protobuf packet file....")
		os.Exit(1)
	}
	return b[0:n]
}

func init() {
	var err error
	uConfig, err = config.Bytes([]byte(`
{
	"senders": {
		"common": {
			"type": "test"
		}
	},
	"receivers": {
		"udp1": {
			"type": "udp",
			"address": "[::1]:1939",
			"handler": "protobuf"
		},
		"udp2": {
			"type": "udp",
			"address": "[::1]:1959",
			"handler": "json"
		}
	},
	"handlers": {
		"protobuf": {
			"parser": "protobuf",
			"transformers": [],
			"sender": "common"
		},
		"json": {
			"parser": "json",
			"transformers": [],
			"sender": "common"
		}
	}
}`))

	if err != nil {
		fmt.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}
	pFile = readProtobufFile("../parser/protobuf-packet.bin")
	var udpAddr1 *net.UDPAddr
	var udpAddr2 *net.UDPAddr
	udpAddr1, err = net.ResolveUDPAddr("udp", "[::1]:1939")
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	udpAddr2, err = net.ResolveUDPAddr("udp", "[::1]:1959")
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	u1, _ = net.DialUDP("udp", nil, udpAddr1)
	u1.SetWriteBuffer(9000)
	u2, _ = net.DialUDP("udp", nil, udpAddr2)
	u2.SetWriteBuffer(9000)

	rUDP1 := uConfig.Receivers["udp1"].Receiver.(*receiver.UDP)
	rUDP2 := uConfig.Receivers["udp2"].Receiver.(*receiver.UDP)
	go rUDP1.Start()
	go rUDP2.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
}

func sendUDP(u *net.UDPConn, b []byte) {
	n, err := u.Write(b)
	if n == 0 || err != nil {
		fmt.Printf("sent %d bytes of %d. err: %v. \n", n, len(b), err)
	}
}

type dummySender struct{}
type dummyJSONSender struct{}

func (d *dummyJSONSender) Send(c *skogul.Container) error {
	sendUDP(u2, pJSON)
	return nil
}

func (d *dummySender) Send(c *skogul.Container) error {
	sendUDP(u1, pFile)
	return nil
}

func TestUDP(t *testing.T) {
	sCommon := uConfig.Senders["common"].Sender.(*sender.Test)

	ds1 := &dummySender{}
	ds2 := &dummyJSONSender{}
	sCommon.SetSync(true)
	sCommon.TestSync(t, ds1, &validContainer, 1, 1)
	sCommon.TestSync(t, ds2, &validContainer, 1, 1)
}

func BenchmarkUDP_protobuf(b *testing.B) {
	sCommon := uConfig.Senders["common"].Sender.(*sender.Test)
	ds := &dummySender{}
	sCommon.SetSync(true)
	for i := 0; i < b.N; i++ {
		sCommon.TestSync(b, ds, &validContainer, 10, 10)
	}
}

func BenchmarkUDP_json(b *testing.B) {
	sCommon := uConfig.Senders["common"].Sender.(*sender.Test)
	ds := &dummyJSONSender{}
	sCommon.SetSync(true)
	for i := 0; i < b.N; i++ {
		sCommon.TestSync(b, ds, &validContainer, 10, 10)
	}
}
