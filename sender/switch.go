/*
 * skogul, switch-sender
 *
 * Copyright (c) 2022 Telenor Norge AS
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
	"github.com/telenornms/skogul"
)

/*
Switch sender sends metrics selectively based on metadata.

Example config:

	{
		"type": "switch",
		"map": [
			{
				"conditions": [{
					"router": "foo",
					"interface": "ae12"
				},
				{
					"router": "bar",
					"interface": "ae0"
				}],
				"next": "customerExportA"
			},
			{
				"conditions": [{
					"router": "foo",
					"interface": "ae5"
				}],
				"next": "customerExportB"
			}
		],
		"default": "log-no-customer"
*/
type Switch struct {
	Default *skogul.SenderRef   `doc:"Default sender to use if no other match is made. If not specified, metrics are discarded."`
	Map     []Match             `doc:"List of match conditions."`
	Next    []*skogul.SenderRef `doc:"Ordered list of senders that will potentially receive metrics."`
}

// Match describes a list of conditions that need to match for a sender to
// receive metrics.
type Match struct {
	Conditions []map[string]interface{} `doc:"Array of metadata headers and required values."`
	Next       *skogul.SenderRef        `doc:"Sender to use in case of a match."`
}

var swLog = skogul.Logger("sender", "switch")

func (cond Match) check(metric *skogul.Metric) bool {
	for _, c := range cond.Conditions {
		match := true
		for key, value := range c {
			if metric.Metadata[key] != value {
				match = false
			}
		}
		if match {
			return true
		}
	}
	return false
}

// Send sends data down stream. Note that it is allowed to create new
// containers, but we CAN NOT modify the original, and we CAN NOT modify
// the metrics them self. This means that this is not an optimal
func (sw *Switch) Send(c *skogul.Container) error {
	newDefault := skogul.Container{}
	newCond := make(map[skogul.Sender]*skogul.Container)
	for _, metric := range c.Metrics {
		nMatch := 0
		for _, mp := range sw.Map {
			match := mp.check(metric)
			if match {
				if newCond[mp.Next.S] == nil {
					cont := skogul.Container{}
					newCond[mp.Next.S] = &cont
				}
				newCond[mp.Next.S].Metrics = append(newCond[mp.Next.S].Metrics, metric)
				nMatch++
				continue
			}
		}
		if nMatch == 0 {
			newDefault.Metrics = append(newDefault.Metrics, metric)
		}
	}
	for sender, cont := range newCond {
		err := sender.Send(cont)
		if err != nil {
			swLog.WithError(err).Warnf("failed to send conditional metrics")
		}
	}
	if len(newDefault.Metrics) > 0 && sw.Default != nil {
		return sw.Default.S.Send(&newDefault)
	}
	return nil
}
