package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/telenornms/skogul"
)

// InfluxDB provides a byte sequence parser for the InfluxDB Line Protocol
// https://docs.influxdata.com/influxdb/v1.7/write_protocols/line_protocol_tutorial/
type InfluxDB struct{}

// InfluxDBLineProtocol is a struct with the same data types as defined in the InfluxDB
// Line Protocol; namely the measurement name, a set of tags, a set of fields and a timestamp.
type InfluxDBLineProtocol struct {
	measurement string
	tags        map[string]interface{}
	fields      map[string]interface{}
	timestamp   time.Time
}

var influxLogger = skogul.Logger("parser", "influxdb")

// Parse marshals a byte sequence of InfluxDB line protocol into a skogul container
func (influxdb InfluxDB) Parse(bytes []byte) (*skogul.Container, error) {
	s := string(bytes)

	// Do we receive data with \r\n?
	lines := strings.Split(s, "\n")

	container := skogul.Container{
		Metrics: make([]*skogul.Metric, len(lines)),
	}

	errors := make([]skogul.Error, 0)

	for i, line := range lines {
		influxLine := InfluxDBLineProtocol{}
		if err := influxLine.ParseLine(line); err != nil {
			// skogul.Error.Source shows index in list and actual line that failed parsing.
			// Skip either? Both?
			errors = append(errors, skogul.Error{Source: fmt.Sprintf("%d-'%s'", i, line), Reason: "Failed to parse influx line", Next: err})
			influxLogger.WithError(err).Error("Failed to parse influx line protocol")
			continue
		}

		container.Metrics[i] = influxLine.Metric()
	}

	if len(errors) > 0 {
		return &container, skogul.Error{
			Source: "parser-influxdb",
			Reason: fmt.Sprintf("One or more influxdb line protocol parse failures. Returning %d successful parses and skipping %d errors.", len(container.Metrics), len(errors)),
		}
	}

	return &container, nil
}

// ParseLine parses a single line into an internal InfluxDBLineProtocol
func (line *InfluxDBLineProtocol) ParseLine(s string) error {
	sections := strings.Split(s, " ")

	if len(sections) < 2 || len(sections) > 3 {
		return skogul.Error{
			Source: "parser-influxdb",
			Reason: fmt.Sprintf("Invalid number of sections in influxdb line protocol line, expected '2' or '3', but got '%d'", len(sections))}
	}

	if len(sections) == 2 {
		// Create own timestamp if it doesn't exist in the source line
		line.timestamp = skogul.Now()
	} else {
		nsTime, err := strconv.ParseInt(sections[2], 0, 64)
		if err != nil {
			return skogul.Error{Source: "parser-influxdb", Reason: "Failed to parse time for influxdb line protocol", Next: err}
		}
		line.timestamp = time.Unix(0, nsTime)
	}

	line.measurement = strings.Split(sections[0], ",")[0]

	tags := sections[0][len(line.measurement)+1:]

	line.tags = make(map[string]interface{})
	for _, tag := range strings.Split(tags, ",") {
		tagValue := strings.Split(tag, "=")
		line.tags[tagValue[0]] = tagValue[1]
	}

	line.fields = make(map[string]interface{})
	for _, field := range strings.Split(sections[1], ",") {
		fieldValue := strings.Split(field, "=")
		line.fields[fieldValue[0]] = fieldValue[1]
	}

	return nil
}

// Metric converts an internal InfluxDBLineProtocol struct to a skogul.Metric
func (line *InfluxDBLineProtocol) Metric() *skogul.Metric {
	metric := skogul.Metric{
		Time:     &line.timestamp,
		Metadata: line.tags,
		Data:     line.fields,
	}

	return &metric
}
