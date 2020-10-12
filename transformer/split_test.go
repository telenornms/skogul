/*
 * skogul, split transformer tests
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
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

func TestSplit(t *testing.T) {
	var c skogul.Container
	testData := `
	{
		"metrics": [
		{
			"data": {
				"data": [
				{
					"splitField": "key1",
					"data": "yes"
				},
				{
					"splitField": "key2",
					"data": "yes also"
				}
				]
			}

		},
		{
			"data": {
				"data": "bad"
			}
		},
		{
			"data": {
				"data": [
				{
					"splitField": "key3",
					"data": "2yes"
				},
				{
					"splitField": "key4",
					"data": "2yes also"
				}
				]
			}
		}
		]
	}
	`
	if err := json.Unmarshal([]byte(testData), &c); err != nil {
		t.Error(err)
		return
	}

	split_path := "data"
	metadata := transformer.Split{
		Field: []string{split_path},
	}

	if err := metadata.Transform(&c); err != nil {
		t.Error(err)
		return
	}

	if len(c.Metrics) != 5 {
		t.Errorf(`Expected c.Metrics to be of len %d but got %d`, 5, len(c.Metrics))
		return
	}

	// Verify that the data is not the same in the two objects as it might differ
	if c.Metrics[0].Data["data"] != "yes" {
		t.Errorf(`Expected Metrics Data to contain key of val '%s' but got '%s'`, "yes", c.Metrics[0].Data["data"])
	}

	if c.Metrics[1].Data["data"] != "yes also" {
		t.Errorf(`Expected Metrics Data to contain key of val '%s' but got '%s'`, "yes also", c.Metrics[1].Data["data"])
	}
	if c.Metrics[2].Data["data"] != "bad" {
		t.Errorf(`Expected Metrics Data to contain key of val '%s' but got '%s'`, "bad", c.Metrics[2].Data["data"])
	}
	if c.Metrics[3].Data["data"] != "2yes" {
		t.Errorf(`Expected Metrics Data to contain key of val '%s' but got '%s'`, "2yes", c.Metrics[3].Data["data"])
	}
	if c.Metrics[4].Data["data"] != "2yes also" {
		t.Errorf(`Expected Metrics Data to contain key of val '%s' but got '%s'`, "2yes also", c.Metrics[4].Data["data"])
	}
}
