/*
 * Copyright (c) 2024 Telenor Norge AS
 * Author(s):
 *  - Hans Rafaelsen <hans.rafaelsen@telenor.no>
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

package transformer_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/transformer"
)

var ratefields_json = `[
        ["if_in_octets", "if_in_bps" ],
        ["if_out_octets", "if_out_bps"],
        ["if_in_bcast_pkts", "if_in_bcast_pps"],
        ["if_in_mcast_pkts", "if_in_mcast_pps"],
        ["if_in_pause_pkts", "if_in_pause_pps"],
        ["if_in_pkts", "if_in_pps"],
        ["if_in_ucast_pkts", "if_in_ucast_pps"],
        ["if_out_bcast_pkts", "if_out_bcast_pps"],
        ["if_out_mcast_pkts", "if_out_mcast_pps"],
        ["if_out_pause_pkts", "if_out_pause_pps"],
        ["if_out_pkts", "if_out_pps"],
        ["if_out_ucast_pkts", "if_out_ucast_pps"],
        ["in_stats__if_unknown_proto_pkts", "in_stats__if_unknown_proto_pps"],
        ["out_stats__if_unknown_proto_pkts", "out_stats__if_unknown_proto_pps"],
        ["if_out_queue_bytes", "if_out_queue_bps"],
        ["if_out_queue_packets", "if_out_queue_pps"],
        ["if_out_queue_red_drop_bytes", "if_out_queue_red_drop_bps"],
        ["if_out_queue_red_drop_packets", "if_out_queue_red_drop_pps"],
        ["if_out_queue_rl_drop_bytes", "if_out_queue_rl_drop_bps"],
        ["if_out_queue_rl_drop_packets", "if_out_queue_rl_drop_pps"],
        ["if_out_queue_tail_drop_packets", "if_out_queue_tail_drop_pps" ]
      ]`

var bitfields_json = `[ 
        "if_in_octets",
        "if_out_octets",
        "if_out_queue_bytes",
        "if_out_queue_packets",
        "if_out_queue_red_drop_bytes",
        "if_out_queue_rl_drop_bytes"
      ]`

var qosfields_json = `[
          "if_out_queue_bytes",
          "if_out_queue_packets",
          "if_out_queue_red_drop_bytes",
          "if_out_queue_red_drop_packets",
          "if_out_queue_rl_drop_bytes",
          "if_out_queue_rl_drop_packets",
          "if_out_queue_tail_drop_packets"
      ]`

func TestIfInOctets(t *testing.T) {

	metric := skogul.Metric{}
	metric.Metadata = make(map[string]any)
	metric.Metadata["systemId"] = "xxx1"
	ti, err := time.Parse(time.RFC3339, "2025-01-01T12:00:00.000+02:00")
	if err != nil {
		t.Errorf("Failed to parse time %v", err)
		return
	} else {
		metric.Time = &ti
	}
	metric.Metadata["if_name"] = "SNMPv2-SMI::enterprises.2011.5.25.31.1.1.3.1.8.67305550"
	metric.Data = make(map[string]any)
	metric.Data["if_in_octets"] = float64(100)
	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric}

	m2 := skogul.Metric{}
	m2.Metadata = make(map[string]any)
	m2.Metadata["systemId"] = "xxx1"
	ti2, err := time.Parse(time.RFC3339, "2025-01-01T12:05:00.000+02:00")
	if err != nil {
		t.Errorf("Failed to parse time %v", err)
		return
	} else {
		m2.Time = &ti2
	}
	m2.Metadata["if_name"] = "SNMPv2-SMI::enterprises.2011.5.25.31.1.1.3.1.8.67305550"
	m2.Data = make(map[string]any)
	m2.Data["if_in_octets"] = float64(400)
	c2 := skogul.Container{}

	c2.Metrics = []*skogul.Metric{&m2}

	var ratefields [][]string
	json.Unmarshal([]byte(ratefields_json), &ratefields)
	conv := transformer.Rate{}
	conv.RateFields = ratefields
	var bitfields []string
	json.Unmarshal([]byte(bitfields_json), &bitfields)
	conv.BitFields = bitfields
	var qosfields []string
	json.Unmarshal([]byte(qosfields_json), &qosfields)
	conv.QosFields = qosfields
	conv.SystemId = "systemId"
	conv.InterfaceId = "if_name"
	conv.QoSClass = "queue_number"

	err = conv.Verify()
	if err != nil {
		t.Errorf("Rate Verify() returned non-nil err: %v", err)
		return
	}

	err = conv.Transform(&c)
	if err != nil {
		t.Errorf("Rate returned non-nil err: %v", err)
	}

	err = conv.Transform(&c2)
	if err != nil {
		t.Errorf("Rate returned non-nil err: %v", err)
	}

	t.Logf("Container after transform:\n%v", c2)

	v, ok := c2.Metrics[0].Data["if_in_bps"].(float64)
	if !ok {
		t.Fatal("Failed to get 'if_in_bps' field")
	}
	expect := 8.0 // 300 bytes in 300 sec = 1 byte 8 bit
	if v != expect {
		t.Errorf("Failed to compute rate. Got: %f. Expected: %f", v, expect)
	}
}

func TestIfInPkts(t *testing.T) {

	metric := skogul.Metric{}
	metric.Metadata = make(map[string]any)
	metric.Metadata["systemId"] = "xxx1"
	ti, err := time.Parse(time.RFC3339, "2025-01-01T12:00:00.000+02:00")
	if err != nil {
		t.Errorf("Failed to parse time %v", err)
		return
	} else {
		metric.Time = &ti
	}
	metric.Metadata["if_name"] = "SNMPv2-SMI::enterprises.2011.5.25.31.1.1.3.1.8.67305550"
	metric.Data = make(map[string]any)
	metric.Data["if_in_pkts"] = float64(100)
	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric}

	m2 := skogul.Metric{}
	m2.Metadata = make(map[string]any)
	m2.Metadata["systemId"] = "xxx1"
	ti2, err := time.Parse(time.RFC3339, "2025-01-01T12:05:00.000+02:00")
	if err != nil {
		t.Errorf("Failed to parse time %v", err)
		return
	} else {
		m2.Time = &ti2
	}
	m2.Metadata["if_name"] = "SNMPv2-SMI::enterprises.2011.5.25.31.1.1.3.1.8.67305550"
	m2.Data = make(map[string]any)
	m2.Data["if_in_pkts"] = float64(400)
	c2 := skogul.Container{}

	c2.Metrics = []*skogul.Metric{&m2}

	var ratefields [][]string
	json.Unmarshal([]byte(ratefields_json), &ratefields)
	conv := transformer.Rate{}
	conv.RateFields = ratefields
	var bitfields []string
	json.Unmarshal([]byte(bitfields_json), &bitfields)
	conv.BitFields = bitfields
	var qosfields []string
	json.Unmarshal([]byte(qosfields_json), &qosfields)
	conv.QosFields = qosfields
	conv.SystemId = "systemId"
	conv.InterfaceId = "if_name"
	conv.QoSClass = "queue_number"

	err = conv.Verify()
	if err != nil {
		t.Errorf("Rate Verify() returned non-nil err: %v", err)
		return
	}

	err = conv.Transform(&c)
	if err != nil {
		t.Errorf("Rate returned non-nil err: %v", err)
	}

	err = conv.Transform(&c2)
	if err != nil {
		t.Errorf("Rate returned non-nil err: %v", err)
	}

	t.Logf("Container after transform:\n%v", c2)

	v, ok := c2.Metrics[0].Data["if_in_pps"].(float64)
	if !ok {
		t.Fatal("Failed to get 'if_in_pps' field")
	}
	expect := 1.0
	if v != expect {
		t.Errorf("Failed to compute rate. Got: %f. Expected: %f", v, expect)
	}
}

func TestOutQueueBytes(t *testing.T) {

	metric := skogul.Metric{}
	metric.Metadata = make(map[string]any)
	metric.Metadata["systemId"] = "xxx1"
	ti, err := time.Parse(time.RFC3339, "2025-01-01T12:00:00.000+02:00")
	if err != nil {
		t.Errorf("Failed to parse time %v", err)
		return
	} else {
		metric.Time = &ti
	}
	metric.Metadata["if_name"] = "SNMPv2-SMI::enterprises.2011.5.25.31.1.1.3.1.8.67305550"
	metric.Metadata["queue_number"] = float64(3)
	metric.Data = make(map[string]any)
	metric.Data["if_out_queue_bytes"] = float64(100)
	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric}

	m2 := skogul.Metric{}
	m2.Metadata = make(map[string]any)
	m2.Metadata["systemId"] = "xxx1"
	ti2, err := time.Parse(time.RFC3339, "2025-01-01T12:05:00.000+02:00")
	if err != nil {
		t.Errorf("Failed to parse time %v", err)
		return
	} else {
		m2.Time = &ti2
	}
	m2.Metadata["if_name"] = "SNMPv2-SMI::enterprises.2011.5.25.31.1.1.3.1.8.67305550"
	m2.Metadata["queue_number"] = float64(3)
	m2.Data = make(map[string]any)
	m2.Data["if_out_queue_bytes"] = float64(400)
	c2 := skogul.Container{}

	c2.Metrics = []*skogul.Metric{&m2}

	var ratefields [][]string
	json.Unmarshal([]byte(ratefields_json), &ratefields)
	conv := transformer.Rate{}
	conv.RateFields = ratefields
	var bitfields []string
	json.Unmarshal([]byte(bitfields_json), &bitfields)
	conv.BitFields = bitfields
	var qosfields []string
	json.Unmarshal([]byte(qosfields_json), &qosfields)
	conv.QosFields = qosfields
	conv.SystemId = "systemId"
	conv.InterfaceId = "if_name"
	conv.QoSClass = "queue_number"

	err = conv.Verify()
	if err != nil {
		t.Errorf("Rate Verify() returned non-nil err: %v", err)
		return
	}

	err = conv.Transform(&c)
	if err != nil {
		t.Errorf("Rate returned non-nil err: %v", err)
	}

	err = conv.Transform(&c2)
	if err != nil {
		t.Errorf("Rate returned non-nil err: %v", err)
	}

	t.Logf("Container after transform:\n%v", c2)

	v, ok := c2.Metrics[0].Data["if_out_queue_bps"].(float64)
	if !ok {
		t.Fatal("Failed to get 'if_out_queue_bps' field")
	}
	expect := 8.0 // 300 bytes in 300 sec = 1 byte 8 bit
	if v != expect {
		t.Errorf("Failed to compute rate. Got: %f. Expected: %f", v, expect)
	}
}

func TestOutQueueBytesMultiClass(t *testing.T) {

	metric := skogul.Metric{}
	metric.Metadata = make(map[string]any)
	metric.Metadata["systemId"] = "xxx1"
	ti, err := time.Parse(time.RFC3339, "2025-01-01T12:00:00.000+02:00")
	if err != nil {
		t.Errorf("Failed to parse time %v", err)
		return
	} else {
		metric.Time = &ti
	}
	metric.Metadata["if_name"] = "SNMPv2-SMI::enterprises.2011.5.25.31.1.1.3.1.8.67305550"
	metric.Metadata["queue_number"] = float64(3)
	metric.Data = make(map[string]any)
	metric.Data["if_out_queue_bytes"] = float64(100)
	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric}

	// Extra metric but in different class
	metric2 := skogul.Metric{}
	metric2.Metadata = make(map[string]any)
	metric2.Metadata["systemId"] = "xxx1"
	tib, err := time.Parse(time.RFC3339, "2025-01-01T12:03:00.000+02:00")
	if err != nil {
		t.Errorf("Failed to parse time %v", err)
		return
	} else {
		metric2.Time = &tib
	}
	metric2.Metadata["if_name"] = "SNMPv2-SMI::enterprises.2011.5.25.31.1.1.3.1.8.67305550"
	metric2.Metadata["queue_number"] = float64(1)
	metric2.Data = make(map[string]any)
	metric2.Data["if_out_queue_bytes"] = float64(00)
	cb := skogul.Container{}
	cb.Metrics = []*skogul.Metric{&metric2}

	m2 := skogul.Metric{}
	m2.Metadata = make(map[string]any)
	m2.Metadata["systemId"] = "xxx1"
	ti2, err := time.Parse(time.RFC3339, "2025-01-01T12:05:00.000+02:00")
	if err != nil {
		t.Errorf("Failed to parse time %v", err)
		return
	} else {
		m2.Time = &ti2
	}
	m2.Metadata["if_name"] = "SNMPv2-SMI::enterprises.2011.5.25.31.1.1.3.1.8.67305550"
	m2.Metadata["queue_number"] = float64(3)
	m2.Data = make(map[string]any)
	m2.Data["if_out_queue_bytes"] = float64(400)
	c2 := skogul.Container{}

	c2.Metrics = []*skogul.Metric{&m2}

	var ratefields [][]string
	json.Unmarshal([]byte(ratefields_json), &ratefields)
	conv := transformer.Rate{}
	conv.RateFields = ratefields
	var bitfields []string
	json.Unmarshal([]byte(bitfields_json), &bitfields)
	conv.BitFields = bitfields
	var qosfields []string
	json.Unmarshal([]byte(qosfields_json), &qosfields)
	conv.QosFields = qosfields
	conv.SystemId = "systemId"
	conv.InterfaceId = "if_name"
	conv.QoSClass = "queue_number"

	err = conv.Verify()
	if err != nil {
		t.Errorf("Rate Verify() returned non-nil err: %v", err)
		return
	}

	err = conv.Transform(&c)
	if err != nil {
		t.Errorf("Rate returned non-nil err: %v", err)
	}

	err = conv.Transform(&cb)
	if err != nil {
		t.Errorf("Rate returned non-nil err: %v", err)
	}

	err = conv.Transform(&c2)
	if err != nil {
		t.Errorf("Rate returned non-nil err: %v", err)
	}

	t.Logf("Container after transform:\n%v", c2)

	v, ok := c2.Metrics[0].Data["if_out_queue_bps"].(float64)
	if !ok {
		t.Fatal("Failed to get 'if_out_queue_bps' field")
	}
	expect := 8.0 // 300 bytes in 300 sec = 1 byte 8 bit
	if v != expect {
		t.Errorf("Failed to compute rate. Got: %f. Expected: %f", v, expect)
	}
}
