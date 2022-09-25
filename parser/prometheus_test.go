/*
 * skogul, Prometheus parser test
 *
 * Copyright (c) 2022 Telenor Norge AS
 * Author(s):
 *  - Roshini Narasimha Raghavan <roshiragavi@gmail.com>
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
	"os"
	"testing"

	"github.com/telenornms/skogul/parser"
)

func TestPrometheus(t *testing.T) {
	b, err := os.ReadFile("./testdata/prometheus_testdata")

	if err != nil {
		t.Errorf("Failed to read test data file: %v", err)
		return
	}
	// convert byte array to json
	p := parser.PROMETHEUS{}

	container, err := p.Parse(b)

	if err != nil {
		t.Logf("Error while parsing the data. %v", err)
		t.FailNow()
	}

	//expectedTime := 1608520832877
	metricKey1 := "dialer_name"
	metricValue1 := "federate"
	metricKey2 := "instance"
	metricValue2 := "localhost:9090"
	metricKey3 := "job"
	metricValue3 := "prometheus"
	var dataValue float64
	dataValue = 1
	dataKey := "net_conntrack_dialer_conn_attempted_total"

	if container == nil || container.Metrics == nil || len(container.Metrics) == 0 {
		t.Logf("Expected parsed prometheus to return a container with at least 1 metric. Container: %v", container.Describe())
		t.FailNow()
	}
	if container.Metrics[0].Metadata[metricKey1] != metricValue1 {
		t.Logf("Expected parsed prometheus to return a metadata field value")
		t.FailNow()
	}

	if container.Metrics[0].Metadata[metricKey2] != metricValue2 {
		t.Logf("Expected parsed prometheus to return a metadata field value")
		t.FailNow()
	}

	if container.Metrics[0].Metadata[metricKey3] != metricValue3 {
		t.Logf("Expected parsed prometheus to return a metadata field value")
		t.FailNow()
	}
	if container.Metrics[0].Data[dataKey] != dataValue {
		t.Logf("Expected parsed prometheus to return a metadata field value %v", container.Metrics[0].Data[dataKey])
		t.FailNow()
	}
	t.Logf("time: %v", container.Metrics[0].Time)
	t.Logf("container: %s", container.Describe())
}
