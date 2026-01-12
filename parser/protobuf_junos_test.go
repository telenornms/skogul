/*
 * skogul, test protobuf parser
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

package parser_test

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/telenornms/skogul"
	junos_protobuf_telemetry "github.com/telenornms/skogul/gen/junos/telemetry"
	"github.com/telenornms/skogul/parser"
	"google.golang.org/protobuf/proto"
)

/*
bit patterns for 32-bit infinity. See ../common.go for details.
*/

const (
	uvneginf = 0b11111111100000000000000000000000
)

type failer interface {
	Fatalf(format string, args ...any)
	Helper()
}

func readProtobufFile(t failer, file string) []byte {
	t.Helper()
	b := make([]byte, 9000)
	f, err := os.Open(file)
	if err != nil {
		t.Fatalf("unable to open protobuf packet file: %v", err)
	}
	defer f.Close()
	n, err := f.Read(b)
	if err != nil {
		t.Fatalf("unable to read protobuf packet file: %v", err)
	}
	if n == 0 {
		t.Fatalf("read 0 bytes from protobuf packet file....")
	}
	return b[0:n]
}

func TestProtoBuf(t *testing.T) {
	b := readProtobufFile(t, "testdata/protobuf-packet.bin")
	x := parser.ProtoBuf{}
	c, err := x.Parse(b)
	if err != nil {
		t.Errorf("ProtoBuf.Parse(b) failed: %s", err)
	}
	if c == nil {
		t.Errorf("ProtoBuf.Parse(b) returned nil-container")
	}
}

func BenchmarkProtoBufParse(b *testing.B) {
	by := readProtobufFile(b, "testdata/protobuf-packet.bin")
	x := parser.ProtoBuf{}

	b.ReportAllocs()
	for b.Loop() {
		_, err := x.Parse(by)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

// BenchmarkProtoBufUnmarshal benchmarks just the protobuf unmarshal step
func BenchmarkProtoBufUnmarshal(b *testing.B) {
	by := readProtobufFile(b, "testdata/protobuf-packet.bin")

	b.ReportAllocs()
	for b.Loop() {
		telemetrystream := &junos_protobuf_telemetry.TelemetryStream{}
		if err := proto.Unmarshal(by, telemetrystream); err != nil {
			b.Fatalf("Unmarshal failed: %v", err)
		}
	}
}

// BenchmarkProtoBufFullPipeline benchmarks the complete parsing pipeline
// including protobuf unmarshal, extension extraction, JSON conversion
func BenchmarkProtoBufFullPipeline(b *testing.B) {
	// Generate a synthetic telemetry stream for consistent benchmarking
	telemetry := generateOpticsDiag(-40)
	bytes, err := proto.Marshal(&telemetry)
	if err != nil {
		b.Fatalf("Failed to marshal test data: %v", err)
	}

	x := parser.ProtoBuf{}

	b.ReportAllocs()
	for b.Loop() {
		c, err := x.Parse(bytes)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
		if c == nil || len(c.Metrics) == 0 {
			b.Fatal("Parse returned empty container")
		}
	}
}

// BenchmarkProtoBufMemoryFootprint measures memory allocations
func BenchmarkProtoBufMemoryFootprint(b *testing.B) {
	by := readProtobufFile(b, "testdata/protobuf-packet.bin")
	x := parser.ProtoBuf{}

	b.Run("SmallMessage", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = x.Parse(by)
		}
	})

	b.Run("LargeMessage", func(b *testing.B) {
		// Generate a larger message with more optics data
		val := float32(-40)
		eps := junos_protobuf_telemetry.EnterpriseSensors{}
		juniperNetworksSensors := junos_protobuf_telemetry.JuniperNetworksSensors{}
		proto.SetExtension(&eps, junos_protobuf_telemetry.E_JuniperNetworks, &juniperNetworksSensors)

		// Create 100 optics diag entries
		opticsDiags := make([]*junos_protobuf_telemetry.OpticsInfos, 100)
		for j := range 100 {
			ifName := fmt.Sprintf("ge-%d/0/%d", j/10, j%10)
			opticsDiags[j] = &junos_protobuf_telemetry.OpticsInfos{
				IfName: &ifName,
				OpticsDiagStats: &junos_protobuf_telemetry.OpticsDiagStats{
					OpticsLaneDiagStats: []*junos_protobuf_telemetry.OpticsDiagLaneStats{
						{LaneLaserReceiverPowerDbm: &val},
					},
				},
			}
		}

		optics := junos_protobuf_telemetry.Optics{OpticsDiag: opticsDiags}
		proto.SetExtension(&juniperNetworksSensors, junos_protobuf_telemetry.E_JnprOpticsExt, &optics)
		telemetry := generateJunosTelemetryStream("large-test", &eps)

		largeBytes, err := proto.Marshal(&telemetry)
		if err != nil {
			b.Fatalf("Failed to marshal large message: %v", err)
		}

		b.ReportAllocs()
		b.SetBytes(int64(len(largeBytes)))
		for b.Loop() {
			_, _ = x.Parse(largeBytes)
		}
	})
}

func generateJunosTelemetryStream(sensorName string, eps *junos_protobuf_telemetry.EnterpriseSensors) junos_protobuf_telemetry.TelemetryStream {
	systemId := "localhost"
	now := uint64(time.Now().Unix())
	componentId := uint32(1)
	subComponentId := uint32(2)

	return junos_protobuf_telemetry.TelemetryStream{
		SystemId:       &systemId,
		Timestamp:      &now,
		ComponentId:    &componentId,
		SubComponentId: &subComponentId,
		SensorName:     &sensorName,
		Enterprise:     eps,
		// Should this be used ?  Ietf:       (*junos_protobuf_telemetry.IETFSensors)(&juniperNetworksSensors),
	}
}

func generateOpticsDiag(val float32) junos_protobuf_telemetry.TelemetryStream {
	eps := junos_protobuf_telemetry.EnterpriseSensors{}
	juniperNetworksSensors := junos_protobuf_telemetry.JuniperNetworksSensors{}
	proto.SetExtension(&eps, junos_protobuf_telemetry.E_JuniperNetworks, &juniperNetworksSensors)

	ifName := "ge-1/0/1"
	optics := junos_protobuf_telemetry.Optics{
		OpticsDiag: []*junos_protobuf_telemetry.OpticsInfos{
			{
				IfName: &ifName,
				OpticsDiagStats: &junos_protobuf_telemetry.OpticsDiagStats{
					OpticsLaneDiagStats: []*junos_protobuf_telemetry.OpticsDiagLaneStats{
						{
							LaneLaserReceiverPowerDbm: &val,
						},
					},
				},
			},
		},
	}
	proto.SetExtension(&juniperNetworksSensors, junos_protobuf_telemetry.E_JnprOpticsExt, &optics)
	return generateJunosTelemetryStream("foo", &eps)
}

func parseDiagStatsResp(data map[string]any, key string) any {
	opticsDiag, ok := data["Optics_diag"].([]any)
	if !ok {
		fmt.Printf("failed to cast")
	}
	foo, ok := opticsDiag[0].(map[string]any)
	if !ok {
		fmt.Printf("failed to cast 2")
	}
	opticsDiagStats := foo["optics_diag_stats"].(map[string]any)

	opticsLaneDiagStats := opticsDiagStats["optics_lane_diag_stats"].([]any)

	bar, ok := opticsLaneDiagStats[0].(map[string]any)
	if !ok {
		fmt.Printf("failed to cast 3")
	}

	return bar[key]
}

func TestParseJunosProtobufTelemetryStreamOptics(t *testing.T) {
	expected := float64(-40)
	telemetry := generateOpticsDiag(-40)

	bytes, err := proto.Marshal(&telemetry)
	if err != nil {
		t.Errorf("Failed to marshal protobuf message to bytes: %v", err)
		return
	}
	if bytes == nil {
		t.Error("Bytes marshalling resulted in nil")
		return
	}

	protobufParser := parser.ProtoBuf{}
	c, err := protobufParser.Parse(bytes)
	if err != nil {
		t.Errorf("Failed to parse optics diag lane stats protobuf data, err: %v", err)
	}
	if c == nil {
		t.Error("Protobuf parse returned nil-container")
		return
	}
	if len(c.Metrics) == 0 {
		t.Error("Protobuf parse returned container with no metrics")
		return
	}

	got := parseDiagStatsResp(c.Metrics[0].Data, "lane_laser_receiver_power_dbm")
	if got != expected {
		t.Errorf("Expected lane_laser_receiver_power_dbm to be %T(%v), but got %T(%v)", expected, expected, got, got)
	}
}

func TestParseJunosProtobufTelemetryStreamOpticsNegativeInf(t *testing.T) {
	val := float32(math.Float32frombits(uvneginf))
	telemetry := generateOpticsDiag(val)

	bytes, err := proto.Marshal(&telemetry)
	if !skogul.IsInf(val, -1) {
		t.Errorf("Trying to test negative infinity value, but value isn't interpreted as infinity. Value: %f", val)
		return
	}
	if err != nil {
		t.Errorf("Failed to marshal protobuf message to bytes: %v", err)
		return
	}
	if bytes == nil {
		t.Error("Bytes marshalling resulted in nil")
		return
	}

	protobufParser := parser.ProtoBuf{}
	_, err = protobufParser.Parse(bytes)
	if err != nil {
		t.Errorf("Expected parsing -Inf values to NOT return an error, ref issue #194 which should now be ... resolved.")
		return
	}
}
