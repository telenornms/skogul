package parser

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/gogo/protobuf/proto"
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/gen/usp"
)

type P struct {
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

func (p *P) initParserStatistics() {
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
func (p *P) Parse(b []byte) (*skogul.Container, error) {
	p.once.Do(p.initParserStatistics)

	if b == nil {
		atomic.AddUint64(&p.stats.NilData, 1)
	}

	record := p.getUspRecord(b)
	recordMetric := skogul.Metric{
		Time:     nil,
		Metadata: p.createRecordMetadata(record),
		Data:     p.createRecordData(record),
	}

	if recordMetric.Data == nil || recordMetric.Metadata == nil {
		atomic.AddUint64(&p.stats.NilData, 1)
		return nil, errors.New("Metric metadata or data was nil; aborting")
	}

	container := skogul.Container{}
	container.Metrics = make([]*skogul.Metric, 1)
	container.Metrics[0] = &recordMetric
	return &container, nil
}

// Unmarshals []byte into a protoc generated struct and returns it
func (p *P) getUspRecord(d []byte) *usp.Record {
	unmarshaledMessage := &usp.Record{}
	if err := proto.Unmarshal(d, unmarshaledMessage); err != nil {
		atomic.AddUint64(&p.stats.ParseErrors, 1)
	}
	return unmarshaledMessage
}

// Unmarshals []byte consisting of the record payload into
// a protoc generated struct and returns it
func (p *P) getRecordMsgPayload(payload []byte) *usp.Msg {
	msgPayload := &usp.Msg{}

	if err := proto.Unmarshal(payload, msgPayload); err != nil {
		atomic.AddUint64(&p.stats.ParseErrors, 1)
	}

	return msgPayload
}

// Creates a map[string]interface{} of the metadata for skogul.Metric
func (p *P) createRecordMetadata(h *usp.Record) map[string]interface{} {
	var d = make(map[string]interface{})

	d["from_id"] = h.GetFromId()
	d["to_id"] = h.GetToId()
	d["payload_security"] = h.GetPayloadSecurity()
	d["sender_cert"] = h.GetSenderCert()
	d["version"] = h.GetVersion()
	d["mac"] = h.GetMacSignature()
	return d
}

// Creates a map[string]interface{} of the record payload for skogul.Metric
func (p *P) createRecordData(t *usp.Record) map[string]interface{} {
	var jsonMap = make(map[string]interface{})
	payload := p.getRecordMsgPayload(t.GetNoSessionContext().GetPayload())

	jsonMap["event"] = payload.GetBody().GetRequest().GetNotify().GetEvent().GetObjPath()
	jsonMap["event_type"] = payload.GetBody().GetRequest().GetNotify().GetEvent().GetEventName()
	jsonMap["subscription_id"] = payload.GetBody().GetRequest().GetNotify().GetSubscriptionId()
	jsonMap["event_parameters"] = payload.GetBody().GetRequest().GetNotify().GetEvent().GetParams()["Data"]

	return jsonMap
}
