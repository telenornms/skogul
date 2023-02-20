/*
 * skogul, kafka producer/sender
 *
 * Copyright (c) 2023 Telenor Norge AS
 * Author(s):
 *  - Kamil Oracz <kamil.oracz@telenor.no>
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

package sender

import (
	"testing"

	"github.com/telenornms/skogul"
)

func createContainer() *skogul.Container {
	meta := make(map[string]interface{})
	meta["foo"] = "bar"
	data := make(map[string]interface{})
	data["baz"] = "qux"

	metric := skogul.Metric{
		Metadata: meta,
		Data:     data,
	}
	metrics := make([]*skogul.Metric, 0)
	metrics = append(metrics, &metric)

	return &skogul.Container{
		Metrics: metrics,
	}
}

func TestRabbitmq(t *testing.T) {
	data := createContainer()

	r := Rabbitmq{
		Username: "guest",
		Password: "guest",
		Queue:    "test-queue",
	}

	err := r.Send(data)

	if err != nil {
		t.Error(err)
	}
}

func TestRabbitmqTonsOfMessages(t *testing.T) {
	data := createContainer()

	r := Rabbitmq{
		Username: "guest",
		Password: "guest",
		Queue:    "test-queue",
	}

	i := 0
	for i < 100000 {
		err := r.Send(data)

		if err != nil {
			t.Error(err)
		}
		i++
	}
}
