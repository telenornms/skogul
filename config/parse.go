/*
 * skogul, configuration parsing
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

/*
Package config handles Skogul configuration parsing.
*/
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"unicode"

	//"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/sender"
)

type Sender struct {
	Type   string `json:"type,omitempty"`
	Next   []string
	Sender skogul.Sender
}
type Receiver struct {
	Type     string `json:"type,omitempty"`
	Receiver skogul.Receiver
}
type Config struct {
	Receivers map[string]Receiver
	Senders   map[string]Sender
}

type typeAndNextOnly struct {
	Type string
	Next []string
}

func (s *Sender) UnmarshalJSON(b []byte) error {
	to := typeAndNextOnly{}
	if err := json.Unmarshal(b, &to); err != nil {
		return err
	}
	if sender.Auto[to.Type] == nil {
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Unknown sender %v", to.Type)}
	}
	if sender.Auto[to.Type].Alloc == nil {
		return skogul.Error{Source: "config parser", Reason: fmt.Sprintf("Bad sender %v", to.Type)}
	}
	ns := sender.Auto[to.Type].Alloc()
	if err := json.Unmarshal(b, &ns); err != nil {
		return skogul.Error{Source: "config parser", Reason: "Failed marshalling", Next: err}
	}
	s.Sender = ns
	s.Next = to.Next
	return nil
}

type fieldDoc struct {
	Doc     string
	Example string
	Type    string
}

type SenderHelp struct {
	Name   string
	Doc    string
	Fields map[string]fieldDoc
}

func HelpSender(s string) (SenderHelp, error) {
	if sender.Auto[s] == nil {
		return SenderHelp{}, skogul.Error{Source: "config parser", Reason: "No such sender"}
	}
	sh := SenderHelp{}
	sh.Name = s
	sh.Doc = sender.Auto[s].Help
	sh.Fields = make(map[string]fieldDoc)
	news := sender.Auto[s].Alloc()
	st := reflect.TypeOf(news)
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if !unicode.IsUpper(rune(field.Name[0])) {
			continue
		}
		fielddoc := fieldDoc{}
		fielddoc.Type = fmt.Sprintf("%v", field.Type.Kind())
		if doc, ok := field.Tag.Lookup("doc"); ok {
			fielddoc.Doc = doc
			if ex, ok := field.Tag.Lookup("example"); ok {
				fielddoc.Example = fmt.Sprintf("Example: %s", ex)
			}
		}
		sh.Fields[field.Name] = fielddoc
	}
	return sh, nil
}

const helpWidth = 66

/*
Print a table of scheme | desc, wrapping the description at helpWidth.

E.g. assuming small helpWidth value:

Without prettyPrint:

foo:// | A very long line will be wrapped

With:

foo:// | A very long
       | line will
       | be wrapped

We wrap at word boundaries to avoid splitting words.
*/
func prettyPrint(scheme string, desc string) {
	fmt.Printf("%11s |", scheme)
	fields := strings.Fields(desc)
	l := 0
	for _, w := range fields {
		if (l + len(w)) > helpWidth {
			l = 0
			fmt.Printf("\n%11s |", "")
		}
		fmt.Printf(" %s", w)
		l += len(w) + 1
	}
	fmt.Printf("\n")
}

func (sh SenderHelp) Print() {
	fmt.Printf("%s - %s\n", sh.Name, sh.Doc)
	fmt.Printf("Variables:\n")
	for n, f := range sh.Fields {
		prettyPrint(n, fmt.Sprintf("Type: %s", f.Type))
		if f.Doc != "" {
			prettyPrint("", f.Doc)
		}
		if f.Example != "" {
			prettyPrint("", "")
			prettyPrint("", f.Example)
			prettyPrint("", "")
		}
	}
}

func File(f string) (*Config, error) {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, skogul.Error{Source: "config parser", Reason: "Failed to read config file", Next: err}
	}

	rawc := Config{}
	err = json.Unmarshal(dat, &rawc)
	if err != nil {
		return nil, skogul.Error{Source: "config parser", Reason: "Unable to parse JSON config", Next: err}
	}

	for _, s := range rawc.Senders {
		for _, next := range s.Next {
			realNext := rawc.Senders[next]
			sendN, ok := s.Sender.(skogul.SenderNext)
			if ok {
				sendN.Next(realNext.Sender)
			}
		}
	}

	return &rawc, nil
}
