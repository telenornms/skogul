/*
 * skogul, protocol buffers parser
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.com>
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

package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/telenornms/skogul"
	pb "github.com/telenornms/skogul/gen/junos/telemetry"
	"github.com/telenornms/skogul/stats"
)

var pbLog = skogul.Logger("parser", "protobuf")

// ProtoBuf parses a byte string-representation of a Container
type ProtoBuf struct {
	once   sync.Once
	stats  *protobufStats
	ticker *time.Ticker
}

type protobufStats struct {
	Received                     uint64 // Received parse calls
	ParseErrors                  uint64 // Failure to parse the bytes using the protobuf definitions provided
	MissingExtension             uint64 // Missing Protobuf extension
	FailedToCastToJuniperMessage uint64 // We assumed it was a Juniper TelemetryStream message, but it failed to cast to it.
	FailedToJsonMarshal          uint64 // Failed to marshal protobuf data to json (this might fail if the data is not representable in JSON, such as the value '-Inf' as float64)
	FailedToJsonUnmarshal        uint64 // Failed to marshal JSON data back into skogul.Metric
	NilData                      uint64 // Parsed protobuf contains no data/metadata
	Parsed                       uint64 // Successful parses
}

// Parse accepts a byte slice of protobuf data and marshals it into a container
func (x *ProtoBuf) Parse(b []byte) (*skogul.Container, error) {
	x.once.Do(func() {
		// XXX: This doesn't start until the first message is parsed. But it's probably fine for now.
		x.stats = &protobufStats{
			Received:                     0,
			ParseErrors:                  0,
			MissingExtension:             0,
			FailedToCastToJuniperMessage: 0,
			FailedToJsonMarshal:          0,
			FailedToJsonUnmarshal:        0,
			NilData:                      0,
			Parsed:                       0,
		}
		x.ticker = time.NewTicker(stats.DefaultInterval)
		go x.sendStats()
	})
	atomic.AddUint64(&x.stats.Received, 1)
	parsedProtoBuf, err := parseTelemetryStream(b)

	if err != nil {
		atomic.AddUint64(&x.stats.ParseErrors, 1)
		return nil, fmt.Errorf("Failed to parse protocol buffer (err: %s)", err)
	}

	protobufTimestamp := time.Unix(int64(*parsedProtoBuf.Timestamp/1000), int64(*parsedProtoBuf.Timestamp%1000)*1000000)

	metric := skogul.Metric{
		Time:     &protobufTimestamp,
		Metadata: x.createMetadata(parsedProtoBuf),
		Data:     x.createData(parsedProtoBuf),
	}

	if metric.Metadata == nil || metric.Data == nil {
		atomic.AddUint64(&x.stats.NilData, 1)
		return nil, errors.New("Metric metadata or data was nil; aborting")
	}

	container := skogul.Container{}
	container.Metrics = make([]*skogul.Metric, 1)
	container.Metrics[0] = &metric

	atomic.AddUint64(&x.stats.Parsed, 1)
	return &container, err
}

// parseTelemetryStream parses a protocol buffer with the Juniper TelemetryStream
// protobuf definitions
func parseTelemetryStream(protobuffer []byte) (*pb.TelemetryStream, error) {
	telemetrystream := &pb.TelemetryStream{}
	if err := proto.Unmarshal(protobuffer, telemetrystream); err != nil {
		// @ToDo: Consider what to do if failing to unmarshal the protobuf
		// Reasons: Invalid proto spec, invalid data, invalid version of proto spec (?)
		// not necessary to return here if we dont log or anything
		return telemetrystream, err
	}

	return telemetrystream, nil
}

// createMetadata extracts the fields containing metadata from the protocol buffer
// and stores them in a string-interface map to be consumed at a later stage.
func (x *ProtoBuf) createMetadata(telemetry *pb.TelemetryStream) map[string]interface{} {
	var metadata = make(map[string]interface{})

	metadata["systemId"] = telemetry.GetSystemId()
	metadata["sensorName"] = telemetry.GetSensorName()
	metadata["componentId"] = telemetry.GetComponentId()
	metadata["subComponentId"] = telemetry.GetSubComponentId()
	return metadata
}

/* createData creates a string-interface map of skogul.Metric type Data
by first marshalling the protobuf message into json and then parsing
it back in to a string-interface map.
*/
func (x *ProtoBuf) createData(telemetry *pb.TelemetryStream) map[string]interface{} {
	var err error
	defer func() {
		if err != nil {
			systemID := telemetry.GetSystemId()
			sensorName := telemetry.GetSensorName()
			pbLog.Printf("Failed to read protobuf telemetry data. SystemID: %v SensorName: %v", systemID, sensorName)
		}
	}()

	extension, err := proto.GetExtension(telemetry.GetEnterprise(), pb.E_JuniperNetworks)
	if err != nil {
		atomic.AddUint64(&x.stats.MissingExtension, 1)
		pbLog.Debug("Failed to get Juniper protobuf extension, is this really a Juniper protobuf message?")
		return nil
	}

	enterpriseExtension, ok := extension.(proto.Message)
	if !ok {
		atomic.AddUint64(&x.stats.FailedToCastToJuniperMessage, 1)
		pbLog.Debug("Failed to cast to juniper message")
		return nil
	}

	registeredExtensions := proto.RegisteredExtensions(enterpriseExtension)

	var regextensions []*proto.ExtensionDesc
	for _, ext := range registeredExtensions {
		regextensions = append(regextensions, ext)
	}

	availableExtensions, err := proto.GetExtensions(enterpriseExtension, regextensions)

	var jsonMessage []byte
	found := false
	for _, ext := range availableExtensions {
		if ext == nil {
			continue
		}

		if found {
			pbLog.Debug("Multiple extensions found, don't know what to do!")
			return nil
		}

		messageOnly, ok := ext.(proto.Message)
		if !ok {
			pbLog.Debugf("Failed to cast to message: %v", ext)
			return nil
		}

		jsonMessage, err = json.Marshal(messageOnly)
		if err != nil {
			atomic.AddUint64(&x.stats.FailedToJsonMarshal, 1)
			pbLog.WithError(err).Error("Failed to marshal data to JSON")
			return nil
		}

		found = true
	}

	var metrics map[string]interface{}
	if err = json.Unmarshal(jsonMessage, &metrics); err != nil {
		atomic.AddUint64(&x.stats.FailedToJsonUnmarshal, 1)
		pbLog.WithError(err).Debug("Unmarshalling JSON data to string/interface map failed")
		return nil
	}

	delete(metrics, "timestamp")
	delete(metrics, "sensorName")
	delete(metrics, "componentId")
	delete(metrics, "subComponentId")

	atomic.AddUint64(&x.stats.Parsed, 1)
	return metrics
}

// GetStats prepares a skogul metric with stats
// for the protobuf parser.
func (x *ProtoBuf) GetStats() *skogul.Metric {
	now := skogul.Now()
	metric := skogul.Metric{
		Time:     &now,
		Metadata: make(map[string]interface{}),
		Data:     make(map[string]interface{}),
	}
	metric.Metadata["component"] = "parser"
	metric.Metadata["type"] = "protobuf"
	metric.Metadata["identity"] = skogul.Identity[x]

	metric.Data["received"] = x.stats.Received
	metric.Data["parse_errors"] = x.stats.ParseErrors
	metric.Data["missing_protobuf_extension"] = x.stats.MissingExtension
	metric.Data["failed_to_cast_to_juniper_message"] = x.stats.FailedToCastToJuniperMessage
	metric.Data["nil_data"] = x.stats.NilData
	metric.Data["failed_to_json_marshal"] = x.stats.FailedToJsonMarshal
	metric.Data["failed_to_json_unmarshal"] = x.stats.FailedToJsonUnmarshal
	metric.Data["parsed"] = x.stats.Parsed
	return &metric
}

// sendStats sets up a forever-running loop which sends stats
// to the global skogul stats channel at the configured interval.
func (x *ProtoBuf) sendStats() {
	for range x.ticker.C {
		pbLog.Trace("sending stats")
		stats.Chan <- x.GetStats()
	}
}
