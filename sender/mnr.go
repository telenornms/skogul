/*
 * skogul, M&R port collector sender
 *
 * Copyright (c) 2019-2026 Telenor Norge AS
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

package sender

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/telenornms/skogul"
)

var mnrLog = skogul.Logger("sender", "mnr")

/*
MnR sender writes to M&R port collector.

The output format is:

	<optional-action-flag>\t<timestamp>\t<groupname>\t<variable>\t<value>(\t<property>=<value>)*

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
	Address      string `doc:"Address to send data to" example:"192.168.1.99:1234"`
	DefaultGroup string `doc:"Default group to use if the metadatafield group is missing."`
	Action       string `doc:"Optional action flag to prepend to each line. Valid values are 'refresh' and 'delete'."`
	RetryConfig
}

// Verify checks the configuration of the MnR sender.
func (mnr *MnR) Verify() error {
	if mnr.Address == "" {
		return skogul.MissingArgument("Address")
	}
	if mnr.Action != "" {
		action := strings.ToLower(mnr.Action)
		if action != "refresh" && action != "delete" {
			return fmt.Errorf("invalid action %q: must be 'refresh' or 'delete'", mnr.Action)
		}
	}
	return mnr.verifyRetry()
}

/*
Send to MnR.

Implementation details: We need to write each value as its own variable to
MnR, so we start by constructing two buffers for what comes before and after
the key\tvalue, then iterate over m.Data.

Also, we open a new TCP connection for each call to Send() at the moment,
which is really suboptimal for large quantities of data, but ok for
occasional data dumps. If large metric containers are received, the cost will
be negligible. But this should, of course, be fixed in the future.
*/
func (mnr *MnR) Send(c *skogul.Container) error {
	actionPrefix := ""
	switch strings.ToLower(mnr.Action) {
	case "refresh":
		actionPrefix = "+r\t"
	case "delete":
		actionPrefix = "+d\t"
	}

	// Format the full payload up front so a retry can re-send it on a
	// fresh connection.
	var buffer bytes.Buffer
	for _, m := range c.Metrics {
		var bufferpre bytes.Buffer
		var bufferpost bytes.Buffer
		fmt.Fprintf(&bufferpre, "%s%d\t", actionPrefix, m.Time.Unix())
		switch {
		case m.Metadata["group"] != nil:
			fmt.Fprintf(&bufferpre, "%s\t", m.Metadata["group"])
		case mnr.DefaultGroup != "":
			fmt.Fprintf(&bufferpre, "%s\t", mnr.DefaultGroup)
		default:
			fmt.Fprintf(&bufferpre, "group\t")
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
			fmt.Fprintf(&buffer, "%s%s%s\t%v\tname=%s%s\n", bufferpre.String(), pre, key, value, key, bufferpost.String())
		}
	}

	return retryNetwork(&mnr.RetryConfig, mnrLog, func() error {
		d, err := net.DialTimeout("tcp", mnr.Address, defaultDialTimeout)
		if err != nil {
			return fmt.Errorf("unable to connect to MnR at %s: %w", mnr.Address, err)
		}
		defer d.Close()
		nbytes, err := d.Write(buffer.Bytes())
		if err != nil {
			return partialWrite(fmt.Errorf("unable to send data to MnR at %s: %w", mnr.Address, err), nbytes, buffer.Len())
		}
		if nbytes < buffer.Len() {
			return partialWrite(fmt.Errorf("write to MnR at %s succeeded, but not all data written", mnr.Address), nbytes, buffer.Len())
		}
		return nil
	})
}
