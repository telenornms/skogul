/*
 * skogul, marshaling functions
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
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

package skogul

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

/*
SenderMap is a list of all referenced senders. This is used during
configuration loading and should not be used afterwards. However,
it needs to be exported so skogul.config can reach it, and it
needs to be outside of skogul.config to avoid circular dependencies.
*/
var SenderMap []*SenderRef

// HandlerMap keeps track of which named handlers exists. A configuration
// engine needs to iterate over this and back-fill the real handlers.
var HandlerMap []*HandlerRef

// TransformerMap keeps track of the named transformers.
var TransformerMap []*TransformerRef

/*
UnmarshalJSON will unmarshal a sender reference by creating a
SenderRef object and putting it on the SenderMap list. The
configuration system in question needs to iterate over SenderMap
after it has completed the first pass of configuration
*/
func (sr *SenderRef) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	sr.Name = s
	sr.S = nil
	SenderMap = append(SenderMap, sr)
	return nil
}

// MarshalJSON for a reference just prints the name
func (sr *SenderRef) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", sr.Name)), nil
}

// MarshalJSON just returns the Name of the handler reference.
func (hr *HandlerRef) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", hr.Name)), nil
}

// MarshalJSON just returns the Name of the transformer reference.
func (tr *TransformerRef) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", tr.Name)), nil
}

// UnmarshalJSON will create an entry on the HandlerMap for the parsed
// handler reference, so the real handler can be substituted later.
func (hr *HandlerRef) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	hr.Name = s
	hr.H = nil
	HandlerMap = append(HandlerMap, hr)
	return nil
}

// UnmarshalJSON will create an entry on the TransformerMap for the parsed
// transformer reference, so the real transformer can be substituted later.
func (tr *TransformerRef) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	tr.Name = s
	tr.T = nil
	TransformerMap = append(TransformerMap, tr)
	return nil
}

// MarshalJSON provides JSON marshalling for Duration
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON provides JSON unmarshalling for Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
