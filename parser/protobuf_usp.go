package parser

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/gogo/protobuf/proto"
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/gen/usp"
)

type ProtoBuffer struct {
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

func (p *ProtoBuffer) initParserStatistics() {
	p.stats = &statistics{
		Received:              0,
		ParseErrors:           0,
		FailedToJsonMarshal:   0,
		FailedToJsonUnmarshal: 0,
		NilData:               0,
		Parsed:                0,
	}
}

func (p *ProtoBuffer) Parse(b []byte) (*skogul.Container, error) {
	p.once.Do(p.initParserStatistics)

	record := p.getUspRecord(b)
	metric := skogul.Metric{
		Time:     nil,
		Metadata: p.createMetadata(record),
		Data:     p.createData(record),
	}

	if metric.Data == nil || metric.Metadata == nil {
		return nil, errors.New("Metric metadata or data was nil; aborting")
	}

	container := skogul.Container{}
	container.Metrics = make([]*skogul.Metric, 1)
	container.Metrics[0] = &metric

	return &container, nil
}

func (p *ProtoBuffer) getUspRecord(d []byte) *usp.Record {
	unmarshaledMessage := &usp.Record{}
	if err := proto.Unmarshal(d, unmarshaledMessage); err != nil {
		atomic.AddUint64(&p.stats.ParseErrors, 1)
	}
	return unmarshaledMessage
}

func (p *ProtoBuffer) createMetadata(h *usp.Record) map[string]interface{} {
	var d = make(map[string]interface{})
	/*
		d["message_id"] = h.GetHeader().GetMsgId()
		d["message_type"] = h.GetHeader().GetMsgType()

		d["request_type"] = h.GetBody().GetRequest().GetReqType()
	*/
	return d
}

func (p *ProtoBuffer) createData(t *usp.Record) map[string]interface{} {
	var d = make(map[string]interface{})

	// check what body is sent
	//d["body"] = t.GetBody().GetMsgBody()
	return d
}
