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

type Prometheus struct{}

func (data Prometheus) Parse(b []byte) (*skogul.Container, error) {
	reader := bytes.NewBuffer(b)
	var parser expfmt.TextParser
	// parse prometheus metrics
	mf, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, err
	}
	container := skogul.Container{
		Metrics: make([]*skogul.Metric, 0, len(mf)),
	}
	tmpMetric := make([]skogul.Metric, len(mf))
	metadataDict := make(map[string]interface{})
	dataDict := make(map[string]interface{})
	var tm time.Time
	indexCounter := 0
	for k, v := range mf {
		for _, i := range v.GetMetric() {
			for _, l := range i.GetLabel() {
				metadataDict[l.GetName()] = l.GetValue()
			}
			// convert int64 timestamp to time.Time
			tm = time.UnixMilli(i.GetTimestampMs())
			if !tm.IsZero() {
				tmpMetric[indexCounter].Time = &tm
			} else {
				tm = skogul.Now()
				tmpMetric[indexCounter].Time = &tm
			}
			Metadatastr, _ := json.Marshal(metadataDict)
			err := json.Unmarshal(Metadatastr, &tmpMetric[indexCounter].Metadata)
			if err != nil {
				return nil, err
			}
			// we do not need tmp to iteratate to get to the value. The library offers GetUntyped().Value call to directly get the value.
			dataDict[k] = i.GetUntyped().Value
			dataDictstr, _ := json.Marshal(dataDict)
			err1 := json.Unmarshal(dataDictstr, &tmpMetric[indexCounter].Data)
			if err1 != nil {
				return nil, err
			}
			container.Metrics = append(container.Metrics, &tmpMetric[indexCounter])
			// clean up the old values of the dictionary so that they don't get carried to the next iteration.
			metadataDict = make(map[string]interface{})
			dataDict = make(map[string]interface{})
			indexCounter++
		}
	}
	return &container, err
}
