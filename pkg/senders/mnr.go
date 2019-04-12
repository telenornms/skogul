/*
 * skogul, M&R port collector sender
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

package senders

import (
	"bytes"
	"fmt"
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"net"
)

/*
MnR sender writes to M&R port collector.

The output format is:

	<timestamp>\t<groupname>\t<variable>\t<value>(\t<property>=<value>)*

Example:

	1199145600 group myDevice.Variable1 100 device=myDevice name=MyVariable1

Two special metadata fields can be provided: "group" will set the M&R storage group,
and "prefix" will be used to prefix all individual data variables.

E.g:

	{
	    "template": {
		    "timestamp": "2019-03-15T11:08:02+01:00",
		    "metadata": {
			"server": "somewhere.example.com"
		    }
	    },
	    "metrics": [
		{
		    "metadata": {
			"prefix": "myDevice.",
			"key": "value",
			"paramkey": "paramvalue"
		    },
		    "data": {
			"astring": "text",
			"float": 1.11,
			"integer": 5
		    }
		}
	    ]
	}

Will result in:

	1552644482	group	myDevice.astring	text		key=value	paramkey=paramvalue	server=somewhere.example.com
	1552644482	group	myDevice.float	1.11		key=value	paramkey=paramvalue	server=somewhere.example.com
	1552644482	group	myDevice.integer	5		key=value	paramkey=paramvalue	server=somewhere.example.com

The default group is set to that of MnR DefaultGroup. If this is unset, the
default group is "group". Meaning:

- If metadata provides "group" key, this is used
- Otherwise, if DefaultGroup is set in MnR sender, this is used
- Otherwise, "group" is used.

*/
type MnR struct {
	Address      string
	DefaultGroup string
}

/*
Sends to MnR.

Implementation details: We need to write each value as its own variable to
MnR, so we start by constructing two buffers for what comes before and after
the key\tvalue, then iterate over m.Data.

Also, we open a new TCP connection for each call to Send() at the moment,
which is really suboptimal for large quantities of data, but ok for
occasional data dumps. If large metric containers are received, the cost will
be negligible. But this should, of course, be fixed in the future.
*/
func (mnr *MnR) Send(c *skogul.Container) error {
	d, err := net.Dial("tcp", mnr.Address)
	if err != nil {
		log.Print("Failed to connect to MnR: %s", err)
		return err
	}
	for _, m := range c.Metrics {
		var bufferpre bytes.Buffer
		var bufferpost bytes.Buffer
		fmt.Fprintf(&bufferpre, "%d\t", m.Time.Unix())
		if m.Metadata["group"] == nil {
			if mnr.DefaultGroup == "" {
				fmt.Fprintf(&bufferpre, "group\t")
			} else {
				fmt.Fprintf(&bufferpre, "%s\t", mnr.DefaultGroup)
			}
		} else {
			fmt.Fprintf(&bufferpre, "%s\t", m.Metadata["group"])
		}
		pre := ""
		if m.Metadata["prefix"] != nil {
			pre = m.Metadata["prefix"].(string)
		}
		for key, value := range m.Metadata {
			if key != "prefix" && key != "group" {
				fmt.Fprintf(&bufferpost, "\t%s=%v", key, value)
			}
		}
		for key, value := range m.Data {
			fmt.Fprintf(d, "%s%s%s\t%v\tname=%s%s\n", bufferpre.String(), pre, key, value, key, bufferpost.String())
		}
	}
	d.Close()
	return nil
}
