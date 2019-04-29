/*
 * skogul, mqtt to influx
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

package main
/*
Forward from MQTT to HTTP(/2).
*/
import (
	"flag"
	"github.com/KristianLyng/skogul/pkg"
	"github.com/KristianLyng/skogul/pkg/parsers"
	"github.com/KristianLyng/skogul/pkg/receivers"
	"github.com/KristianLyng/skogul/pkg/senders"
	"github.com/KristianLyng/skogul/pkg/transformers"
)

var flisten = flag.String("listen", "mqtt://localhost:1883/", "Address for MQTT broker")
var ftarget = flag.String("target", "http://127.0.0.1:8086/write?db=test", "InfluxDB to write to")

func main() {
	flag.Parse()
	target := &senders.HTTP{URL: *ftarget}
	fanout := &senders.Fanout{Next: target}

	h := skogul.Handler{
		Parser:       parsers.JSON{},
		Sender:       fanout,
		Transformers: []skogul.Transformer{transformers.Templater{}}}

	receiver := receivers.MQTT{Address: *flisten, Handler: &h}
	receiver.Start()
}
