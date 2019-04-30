/*
 * skogul, tcpline receiver
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

package receivers

import (
	"bufio"
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"net"
	"net/url"
)

/*
TCPLine listens on a IP:TCP port specified in the Address string and accepts
one container per line to be sent to the parser.

Example usage, assuming JSON parser:

	$ cat payloads/simple.json  | jq -c . | nc '::1' '1234'

Since this is not possible to secure, it should be avoided where possible and
placed as close to the data source. A good use of this model is to use a TCPLine
receiver on the same box that needs to write to it, combined with
skogul.senders.HTTP to forward over a more sensible channel.
*/
type TCPLine struct {
	Address string
	Handler *skogul.Handler
}

/*

Start the TCP line receiver and run forever.

We close the write-side of the connection leaving it to the other side to
finish up. We should probably add a read-timeout in the future.
*/
func (tl *TCPLine) Start() error {
	if tl.Address == "" {
		tl.Address = "[::1]:1234"
	}
	tcpip, err := net.ResolveTCPAddr("tcp", tl.Address)
	if err != nil {
		log.Printf("Can't resolve %s: %v", tl.Address, err)
		return err
	}
	ln, err := net.ListenTCP("tcp", tcpip)
	if err != nil {
		log.Printf("Can't listen on %s: %v", tl.Address, err)
		return err
	}
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Printf("Unable to accept connection: %v", err)
			continue
		}
		conn.CloseWrite()
		go tl.handleConnection(conn)
	}
	return skogul.Error{Reason: "Shouldn't reach this"}
}

func (tl *TCPLine) handleConnection(conn *net.TCPConn) error {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		log.Printf("Read %s", bytes)
		m, err := tl.Handler.Parser.Parse(bytes)
		if err == nil {
			err = m.Validate()
		}
		if err != nil {
			log.Printf("Unable to parse JSON: %s", err)
			continue
		}
		for _, t := range tl.Handler.Transformers {
			t.Transform(&m)
		}
		tl.Handler.Sender.Send(&m)
	}
	if err := scanner.Err(); err != nil {
		log.Print("Error reading: %s", err)
		return skogul.Error{Reason: "Error reading file"}
	}
	return nil
}

func init() {
	addAutoReceiver("tcp", NewTCPLine, "Listen for Skogul-formatted JSON on a line-separate tcp socket")
}

/*
NewTCPLine returns a new TCPLine receiver built from the url. Correct format is
tcp://ip:port
*/
func NewTCPLine(ul url.URL, h skogul.Handler) skogul.Receiver {
	return &TCPLine{Address: ul.String(), Handler: &h}
}
