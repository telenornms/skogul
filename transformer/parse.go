package transformer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/parser"
)

var parseLog = skogul.Logger("transformer", "parser")

type Parse struct {
	Parser              string `doc:"Name of skogul parser to use. Note: has to be auto-initialisable."` // Feature request: Re-use parser unmarshal logic here
	Source              string `doc:"Name of data field to apply parser to"`
	Destination         string `doc:"Field to create with the parsed data (default: [Source]_data)"`
	DestinationMetadata string `doc:"Field to create with the parsed metadata (default: [Source]_metadata)"`
	Keep                bool   `doc:"Keep the unparsed value (default: false)"`
	Append              bool   `doc:"Insert the values directly in the existing container, instead of creating new fields (destination and destination_metadata). (default: false)"`
	once                sync.Once
	parser              skogul.Parser
	ok                  bool
}

func (p *Parse) init() {
	p.ok = false
	if p.Destination == "" {
		p.Destination = fmt.Sprintf("%s_data", p.Source)
	}
	if p.DestinationMetadata == "" {
		p.DestinationMetadata = fmt.Sprintf("%s_metadata", p.Source)
	}

	parseLog.Debug("Trying to find parser to initialise")
	for k, pars := range parser.Auto {
		if strings.ToLower(k) == strings.ToLower(p.Parser) {
			parseLog.Debugf("Found parser '%s', using it", k)
			p.parser = pars.Alloc().(skogul.Parser)
			p.ok = true
			break
		}
	}
	if p.parser == nil {
		parseLog.Errorf("Failed to find parser with name '%s'", p.Parser)
		p.ok = false
		return
	}
}

func (p *Parse) Transform(c *skogul.Container) error {
	p.once.Do(func() {
		p.init()
	})

	if !p.ok {
		return skogul.Error{Reason: "Parser transformer not initialized", Source: "transformer-parser"}
	}

	for i, metric := range c.Metrics {
		val, ok := metric.Data[p.Source].(string)
		if !ok {
			parseLog.Error("Failed to cast value to string")
			continue
		}
		b := []byte(val)
		parsed, err := p.parse(b)
		if err != nil {
			return err
		}
		if len(parsed.Metrics) == 0 {
			return skogul.Error{Reason: "Parsed 0 metrics", Source: "transformer-parser"}
		}

		if !p.Keep {
			delete(c.Metrics[i].Data, p.Source)
		}

		if p.Append {
			// Insert directly into existing container
			for field, val := range parsed.Metrics[0].Metadata {
				c.Metrics[i].Metadata[field] = val
			}
			for field, val := range parsed.Metrics[0].Data {
				c.Metrics[i].Data[field] = val
			}
		} else {
			// Add new fields
			data := make(map[string]interface{})
			metadata := make(map[string]interface{})

			for field, val := range parsed.Metrics[0].Metadata {
				metadata[field] = val
			}
			for field, val := range parsed.Metrics[0].Data {
				data[field] = val
			}
			c.Metrics[i].Data[p.DestinationMetadata] = metadata
			c.Metrics[i].Data[p.Destination] = data
		}
	}
	return nil
}

func (p *Parse) parse(b []byte) (*skogul.Container, error) {

	parsed, err := p.parser.Parse(b)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func (p *Parse) Verify() error {
	if p.Source == "" {
		return skogul.Error{Reason: "Missing source field", Source: "transformer-parser"}
	}
	if (p.Destination != "" || p.DestinationMetadata != "") && p.Append {
		parseLog.Warn("Destination or DestinationMetadata configured at the same time as Append - Append takes precedence, and Destination(Metadata) will be ignored.")
	}
	return nil
}
