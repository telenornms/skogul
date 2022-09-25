/*
 * skogul,  prometheus parser
 *
 * Copyright (c) 2022 Telenor Norge AS
 * Author(s):
 *  - Roshini Narasimha Raghavan <roshiragavi@gmail.com>
 *  - Kristian Lyngstøl <kly@kly.no>
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
	"time"

	"github.com/prometheus/common/expfmt"
	"github.com/telenornms/skogul"
)

type PROMETHEUS struct{}

func (data PROMETHEUS) Parse(b []byte) (*skogul.Container, error) {
	reader := bytes.NewBuffer(b)
	var parser expfmt.TextParser
	// parse prometheus metrics 
	mf, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, err
	}
	container := skogul.Container{}
	var tmpMetric skogul.Metric
	metadataDict := make(map[string]interface{})
	dataDict := make(map[string]interface{})
	var Time time.Time 

	for k, v := range mf {
		for _, i := range v.GetMetric() {
			container.Metrics = make([]*skogul.Metric, 0, len(v.GetMetric()))
			for _, l := range i.GetLabel() {
				metadataDict[l.GetName()] = l.GetValue()	
			}
			dataDict[k] = i.GetUntyped()
			// convert int64 timestamp to time.Time 
			Time = time.Unix(i.GetTimestampMs(), 0)
			if !Time.IsZero() {
				tmpMetric.Time = &Time 
			} else {
				Time = time.Now()
				tmpMetric.Time = &Time 
			}
			Metadatastr, _ := json.Marshal(metadataDict)
			dataDictstr, _ := json.Marshal(dataDict) 
			err := json.Unmarshal(Metadatastr, &tmpMetric.Metadata)
			if err != nil {
				return nil, err
			}
			err1 := json.Unmarshal(dataDictstr, &tmpMetric.Data)
			if err1 != nil {
				return nil, err
			}
		}
		container.Metrics = append(container.Metrics, &tmpMetric) 
	}
	return &container, err 
}
