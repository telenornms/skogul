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
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"github.com/KristianLyng/skogul"
	pb "github.com/KristianLyng/skogul/gen"
)

// ProtoBuf parses a byte string-representation of a Container
type ProtoBuf struct{}

// Parse accepts a byte slice of protobuf data and marshals it into a container
func (x ProtoBuf) Parse(b []byte) (skogul.Container, error) {
	log.Println("Parsing some protobuf")
	container := skogul.Container{}

	parsedProtoBuf, err := parseTelemetryStream(b)

	if err != nil {
		log.Printf("Failed to parse protocol buffer (err: %s)", err)
		return container, err
	}

	log.Printf("Parsed protocol buffer: %v", parsedProtoBuf)

	// protobuf timestamp is milliseconds, unix timestamp is seconds
	var protobufTimestamp time.Time = time.Unix(int64(*parsedProtoBuf.Timestamp/1000), 0)

	metric := skogul.Metric{
		Time:     &protobufTimestamp,
		Metadata: createMetadata(parsedProtoBuf),
		Data:     createData(parsedProtoBuf),
	}

	if metric.Metadata == nil || metric.Data == nil {
		return container, errors.New("Metric metadata or data was nil; aborting")
	}

	containerMetrics := make([]*skogul.Metric, 1)
	containerMetrics[0] = &metric

	container.Metrics = containerMetrics

	return container, err
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

const telemetryMetadataSystemID = "systemId"
const telemetryMetadataSensorName = "sensorName"

// createMetadata extracts the fields containing metadata from the protocol buffer
// and stores them in a string-interface map to be consumed at a later stage.
func createMetadata(telemetry *pb.TelemetryStream) map[string]interface{} {
	var metadata = make(map[string]interface{})

	metadata[telemetryMetadataSystemID] = telemetry.GetSystemId()
	metadata[telemetryMetadataSensorName] = telemetry.GetSensorName()

	return metadata
}

/* createData creates a string-interface map of skogul.Metric type Data
by first marshalling the raw protobuf data into json and then parsing
it back in to a string-interface map.
@ToDo: Make this cheaper
*/
func createData(telemetry *pb.TelemetryStream) map[string]interface{} {
	pbjsonmarshaler := jsonpb.Marshaler{}

	// pbjsoned, err := pbjsonmarshaler.MarshalToString(telemetry)
	// log.Printf("Protocol buffer as json string (err: %v): %v", err, pbjsoned)

	var out bytes.Buffer
	err := pbjsonmarshaler.Marshal(&out, telemetry)

	if err != nil {
		log.Printf("Marshalling protocol buffer data to JSON failed: %s", err)
		return nil
	}

	var metrics map[string]interface{}
	err = json.Unmarshal(out.Bytes(), &metrics)

	if err != nil {
		log.Printf("Unmarshalling JSON data to string/interface map failed: %s", err)
		return nil
	}

	return metrics
}
