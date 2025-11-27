package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/gogo/protobuf/proto"
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/gen/usp"
)

type USPParser struct {
	once  sync.Once
	stats *statistics
}

type statistics struct {
	Received              uint64 // Received parse calls
	ParseErrors           uint64 // Failure to parse the bytes using the protobuf definitions provided
	FailedToJsonMarshal   uint64 // Failed to marshal protobuf data to json (this might fail if the data is not representable in JSON, such as the value '-Inf' as float64)
	FailedToJsonUnmarshal uint64 // Failed to marshal JSON data back into skogul.Metric
	NilData               uint64 // Parsed protobuf contains no data/metadata
	Parsed                uint64 // Successful parses
}

func (p *USPParser) initParserStatistics() {
	p.stats = &statistics{
		Received:              0,
		ParseErrors:           0,
		FailedToJsonMarshal:   0,
		FailedToJsonUnmarshal: 0,
		NilData:               0,
		Parsed:                0,
	}
}

// Parse accepts a byte slice of protobuf data and marshals it into a container
func (p *USPParser) Parse(b []byte) (*skogul.Container, error) {
	p.once.Do(p.initParserStatistics)
	atomic.AddUint64(&p.stats.Received, 1)

	if b == nil {
		atomic.AddUint64(&p.stats.NilData, 1)
		return nil, errors.New("nil byte slice provided")
	}

	record, err := p.getUspRecord(b)
	if err != nil {
		atomic.AddUint64(&p.stats.ParseErrors, 1)
		return nil, fmt.Errorf("failed to parse protocol buffer: %w", err)
	}

	recordData, err := p.createRecordData(record)
	if err != nil {
		atomic.AddUint64(&p.stats.ParseErrors, 1)
		return nil, fmt.Errorf("failed to create data: %w", err)
	}

	metadata := p.createRecordMetadata(record, recordData)

	eventData, ok := recordData["event_data"].(string)
	if !ok {
		atomic.AddUint64(&p.stats.FailedToJsonUnmarshal, 1)
		return nil, errors.New("event_data is missing or not a string")
	}

	json, err := p.extractJSON(eventData)
	if err != nil {
		atomic.AddUint64(&p.stats.FailedToJsonUnmarshal, 1)
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	recordMetric := skogul.Metric{
		Time:     nil,
		Metadata: metadata,
		Data:     json,
	}

	if recordMetric.Data == nil || recordMetric.Metadata == nil {
		atomic.AddUint64(&p.stats.NilData, 1)
		return nil, errors.New("metric metadata or data was nil; aborting")
	}

	container := skogul.Container{}
	container.Metrics = make([]*skogul.Metric, 1)
	container.Metrics[0] = &recordMetric

	atomic.AddUint64(&p.stats.Parsed, 1)
	return &container, nil
}

// getUspRecord Unmarshals []byte into a protoc generated struct
func (p *USPParser) getUspRecord(d []byte) (*usp.Record, error) {
	unmarshaledMessage := &usp.Record{}
	if err := proto.Unmarshal(d, unmarshaledMessage); err != nil {
		atomic.AddUint64(&p.stats.ParseErrors, 1)
		return nil, fmt.Errorf("failed to unmarshal protocol buffer: %w", err)
	}
	return unmarshaledMessage, nil
}

/*
getRecordMsgPayload unmarshals []byte consisting of the record payload into
a protoc generated struct
*/
func (p *USPParser) getRecordMsgPayload(payload []byte) (*usp.Msg, error) {
	msgPayload := &usp.Msg{}
	if err := proto.Unmarshal(payload, msgPayload); err != nil {
		atomic.AddUint64(&p.stats.ParseErrors, 1)
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return msgPayload, nil
}

// createRecordMetadata creates a map[string]interface{} of the metadata for skogul.Metric
func (p *USPParser) createRecordMetadata(h *usp.Record, xh map[string]any) map[string]any {
	d := make(map[string]any)

	d["event"] = xh["event"]
	d["event_type"] = xh["event_type"]
	d["subscription_id"] = xh["subscription_id"]
	d["from_id"] = h.GetFromId()
	d["to_id"] = h.GetToId()
	d["payload_security"] = h.GetPayloadSecurity()
	d["sender_cert"] = h.GetSenderCert()
	d["version"] = h.GetVersion()
	d["mac"] = h.GetMacSignature()
	return d
}

// extractJSON unmarshals event parameters to json
func (p *USPParser) extractJSON(s string) (map[string]any, error) {
	input := []byte(s)

	var d map[string]any

	if err := json.Unmarshal(input, &d); err != nil {
		return nil, err
	}

	return d, nil
}

// createRecordData creates a map[string]interface{} of the record payload for skogul.Metric
func (p *USPParser) createRecordData(t *usp.Record) (map[string]any, error) {
	jsonMap := make(map[string]any)
	payload, err := p.getRecordMsgPayload(t.GetNoSessionContext().GetPayload())
	if err != nil {
		return nil, err
	}

	// Check if request contains the Notify event. (It could be a different event by mistake)
	reqType := payload.Body.GetRequest().GetReqType()
	if _, ok := reqType.(*usp.Request_Notify); !ok {
		return nil, fmt.Errorf("request does not contain a Notify event, got %T", reqType)
	}

	jsonMap["event"] = payload.GetBody().GetRequest().GetNotify().GetEvent().GetObjPath()
	jsonMap["event_type"] = payload.GetBody().GetRequest().GetNotify().GetEvent().GetEventName()
	jsonMap["subscription_id"] = payload.GetBody().GetRequest().GetNotify().GetSubscriptionId()
	jsonMap["event_data"] = payload.GetBody().GetRequest().GetNotify().GetEvent().GetParams()["Data"]

	return jsonMap, nil
}

// GetStats prepares a skogul metric with stats for the USP parser.
func (p *USPParser) GetStats() *skogul.Metric {
	now := skogul.Now()
	metric := skogul.Metric{
		Time:     &now,
		Metadata: make(map[string]any),
		Data:     make(map[string]any),
	}
	metric.Metadata["component"] = "parser"
	metric.Metadata["type"] = "usp"
	metric.Metadata["identity"] = skogul.Identity[p]

	p.once.Do(p.initParserStatistics)

	metric.Data["received"] = p.stats.Received
	metric.Data["parse_errors"] = p.stats.ParseErrors
	metric.Data["failed_to_json_marshal"] = p.stats.FailedToJsonMarshal
	metric.Data["failed_to_json_unmarshal"] = p.stats.FailedToJsonUnmarshal
	metric.Data["nil_data"] = p.stats.NilData
	metric.Data["parsed"] = p.stats.Parsed
	return &metric
}
