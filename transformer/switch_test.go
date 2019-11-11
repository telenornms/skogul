/*
 * skogul, switch transformer tests
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

package transformer_test

import (
	"encoding/json"
	"testing"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/transformer"
)

const testMetadataA = `{ "sensor": "a" }`
const testMetadataB = `{ "sensor": "b" }`
const testData = `{ "bannable_field": "someValue", "removable_field": "someOtherValue", "data": 42 }`

func generateContainer() skogul.Container {
	metric := skogul.Metric{}
	json.Unmarshal([]byte(testMetadataA), &metric.Metadata)
	json.Unmarshal([]byte(testData), &metric.Data)

	container := skogul.Container{
		Metrics: []*skogul.Metric{&metric},
	}

	return container
}

// func generateTransformer() skogul.Transformer {
// 	ban := transformer.Metadata{
// 		Ban: []string{"a"},
// 	}

// 	return ban
// }

func TestSwitch1(t *testing.T) {
	removeTransformer := transformer.Data{
		Remove: []string{"removable_field"},
	}

	case1 := transformer.Case{
		When:         "sensor",
		Is:           "a",
		Transformers: []string{""},
	}

	config := transformer.Switch{
		Cases: []transformer.Case{case1},
	}

}
