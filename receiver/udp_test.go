/*
 * skogul, udp receiver tests
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

package receiver_test

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/config"
	"github.com/telenornms/skogul/receiver"
	"github.com/telenornms/skogul/sender"
)

// FIXME: This file is very fond of global scope and init(), mainly for the
// sake of benchmarks... and it's a bit messy.
var (
	uConfig *config.Config
	pFile   []byte
	u1      *net.UDPConn
	u2      *net.UDPConn
	u3      *net.UDPConn
	u4      *net.UDPConn
	u5      *net.UDPConn
	u6      *net.UDPConn
	u7      *net.UDPConn
	pJSON   = []byte("{\"metrics\":[{\"timestamp\":\"2019-03-15T11:08:02+01:00\",\"metadata\":{\"key\":\"value\"},\"data\":{\"string\":\"text\",\"float\":1.11,\"integer\":5}}]}")
)

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
		},
		"sleep": {
			"type": "sleep",
			"next": "common",
			"maxdelay": "3ms",
			"base": "0.1ms"
		}
	},
	"receivers": {
		"udp1": {
			"type": "udp",
			"address": "localhost:1939",
			"handler": "protobuf"
		},
		"udp2": {
			"type": "udp",
			"address": "localhost:1959",
			"handler": "json-sleep",
			"Threads": 1
		},
		"udp3": {
			"type": "udp",
			"address": "localhost:1969",
			"handler": "json-sleep",
			"Threads": 10,
			"Buffer": 10240
		},
		"udp4": {
			"type": "udp",
			"address": "localhost:1979",
			"handler": "json-sleep",
			"Threads": 100,
			"Buffer": 1024
		},
		"udp5": {
			"type": "udp",
			"address": "localhost:1989",
			"handler": "json",
			"Threads": 1,
			"Buffer": 65536
		},
		"udp6": {
			"type": "udp",
			"address": "localhost:1999",
			"handler": "json",
			"Threads": 10,
			"Buffer": 65536
		},
		"udp7": {
			"type": "udp",
			"address": "localhost:2009",
			"handler": "json",
			"Threads": 100,
			"Buffer": 65536
		}
	},
	"handlers": {
		"protobuf": {
			"parser": "protobuf",
			"transformers": [],
			"sender": "common"
		},
		"json-sleep": {
			"parser": "json",
			"transformers": [],
			"sender": "sleep"
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
	pFile = readProtobufFile("../parser/testdata/protobuf-packet.bin")
	var udpAddr1 *net.UDPAddr
	var udpAddr2 *net.UDPAddr
	var udpAddr3 *net.UDPAddr
	var udpAddr4 *net.UDPAddr
	var udpAddr5 *net.UDPAddr
	var udpAddr6 *net.UDPAddr
	var udpAddr7 *net.UDPAddr
	udpAddr1, err = net.ResolveUDPAddr("udp", "localhost:1939")
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	udpAddr2, err = net.ResolveUDPAddr("udp", "localhost:1959")
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	udpAddr3, err = net.ResolveUDPAddr("udp", "localhost:1969")
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	udpAddr4, err = net.ResolveUDPAddr("udp", "localhost:1979")
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	udpAddr5, err = net.ResolveUDPAddr("udp", "localhost:1989")
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	udpAddr6, err = net.ResolveUDPAddr("udp", "localhost:1999")
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	udpAddr7, err = net.ResolveUDPAddr("udp", "localhost:2009")
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	u1, err = net.DialUDP("udp", nil, udpAddr1)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		os.Exit(1)
	}
	u1.SetWriteBuffer(9000)
	u2, err = net.DialUDP("udp", nil, udpAddr2)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		os.Exit(1)
	}
	u2.SetWriteBuffer(9000)
	u3, err = net.DialUDP("udp", nil, udpAddr3)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		os.Exit(1)
	}
	u3.SetWriteBuffer(9000)
	u4, err = net.DialUDP("udp", nil, udpAddr4)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		os.Exit(1)
	}
	u4.SetWriteBuffer(9000)
	u5, err = net.DialUDP("udp", nil, udpAddr5)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		os.Exit(1)
	}
	u5.SetWriteBuffer(9000)
	u6, err = net.DialUDP("udp", nil, udpAddr6)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		os.Exit(1)
	}
	u6.SetWriteBuffer(9000)
	u7, err = net.DialUDP("udp", nil, udpAddr7)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		os.Exit(1)
	}
	u7.SetWriteBuffer(9000)

	rUDP1 := uConfig.Receivers["udp1"].Receiver.(*receiver.UDP)
	rUDP2 := uConfig.Receivers["udp2"].Receiver.(*receiver.UDP)
	rUDP3 := uConfig.Receivers["udp3"].Receiver.(*receiver.UDP)
	rUDP4 := uConfig.Receivers["udp4"].Receiver.(*receiver.UDP)
	rUDP5 := uConfig.Receivers["udp5"].Receiver.(*receiver.UDP)
	rUDP6 := uConfig.Receivers["udp6"].Receiver.(*receiver.UDP)
	rUDP7 := uConfig.Receivers["udp7"].Receiver.(*receiver.UDP)
	go rUDP1.Start()
	go rUDP2.Start()
	go rUDP3.Start()
	go rUDP4.Start()
	go rUDP5.Start()
	go rUDP6.Start()
	go rUDP7.Start()
	time.Sleep(time.Duration(100 * time.Millisecond))
}

func sendUDP(u *net.UDPConn, b []byte) {
	n, err := u.Write(b)
	if n == 0 || err != nil {
		fmt.Printf("sent %d bytes of %d. err: %v. \n", n, len(b), err)
	}
}

type (
	dummySender     struct{}
	dummyJSONSender struct {
		sock       *net.UDPConn
		iterations int
	}
)

func (d *dummyJSONSender) Send(c *skogul.Container) error {
	for i := 0; i < d.iterations; i++ {
		sendUDP(d.sock, pJSON)
	}
	return nil
}

func (d *dummySender) Send(c *skogul.Container) error {
	sendUDP(u1, pFile)
	return nil
}

func TestUDP(t *testing.T) {
	sCommon := uConfig.Senders["common"].Sender.(*sender.Test)

	ds1 := &dummySender{}
	ds2 := &dummyJSONSender{u2, 1}
	sCommon.SetSync(true)
	sCommon.TestSync(t, ds1, &validContainer, 1, 1)
	sCommon.TestSync(t, ds2, &validContainer, 1, 1)
}

func BenchmarkUDP_protobuf(b *testing.B) {
	sCommon := uConfig.Senders["common"].Sender.(*sender.Test)
	ds := &dummySender{}
	sCommon.SetSync(true)
	for b.Loop() {
		sCommon.TestSync(b, ds, &validContainer, 10, 10)
	}
}

func BenchmarkUDP_json_Threads1(b *testing.B) {
	sCommon := uConfig.Senders["common"].Sender.(*sender.Test)
	ds := &dummyJSONSender{u5, 20}
	sCommon.SetSync(true)
	for b.Loop() {
		sCommon.TestSync(b, ds, &validContainer, 5, 100)
	}
}

func BenchmarkUDP_json_Threads10(b *testing.B) {
	sCommon := uConfig.Senders["common"].Sender.(*sender.Test)
	ds := &dummyJSONSender{u6, 20}
	sCommon.SetSync(true)
	for b.Loop() {
		sCommon.TestSync(b, ds, &validContainer, 5, 100)
	}
}

func BenchmarkUDP_json_Threads100(b *testing.B) {
	sCommon := uConfig.Senders["common"].Sender.(*sender.Test)
	ds := &dummyJSONSender{u7, 20}
	sCommon.SetSync(true)
	for b.Loop() {
		sCommon.TestSync(b, ds, &validContainer, 5, 100)
	}
}
