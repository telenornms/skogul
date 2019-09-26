package config

import (
	"fmt"
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"reflect"
	"unicode"
)

// fieldDoc is a structured representation of the documentation of a single
// field in a struct, used for both senders and receivers (and more?)
type fieldDoc struct {
	Doc     string
	Example string
	Type    string
}

// Help is the relevant help for a single sender/receiver
type Help struct {
	Name    string
	Aliases string
	Doc     string
	Fields  map[string]fieldDoc
}

// HelpSender looks up documentation for a named sender and provides a
// help-structure. Should probably be merged with HelpReceiver somewhat.
func HelpSender(s string) (Help, error) {
	if sender.Auto[s] == nil {
		return Help{}, skogul.Error{Source: "config parser", Reason: "No such sender"}
	}
	sh := Help{}
	sh.Name = s
	sh.Doc = sender.Auto[s].Help
	for _, alias := range sender.Auto[s].Aliases {
		sh.Aliases = fmt.Sprintf("%s %s", alias, sh.Aliases)
	}
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
		t := fmt.Sprintf("%v", field.Type.Kind())
		typeName := field.Type.Name()
		typeString := field.Type.String()
		if typeName != "" {
			t = typeName
		} else if typeString != "" {
			t = typeString
		}
		fielddoc.Type = fmt.Sprintf("%s", t)
		if doc, ok := field.Tag.Lookup("doc"); ok {
			fielddoc.Doc = doc
			if ex, ok := field.Tag.Lookup("example"); ok {
				fielddoc.Example = ex
			}
		}
		sh.Fields[field.Name] = fielddoc
	}
	return sh, nil
}

// HelpReceiver looks up documentation for a named receiver.
func HelpReceiver(r string) (Help, error) {
	if receiver.Auto[r] == nil {
		return Help{}, skogul.Error{Source: "config parser", Reason: "No such receiver"}
	}
	sh := Help{}
	sh.Name = r
	sh.Doc = receiver.Auto[r].Help
	for _, alias := range receiver.Auto[r].Aliases {
		sh.Aliases = fmt.Sprintf("%s %s", alias, sh.Aliases)
	}
	sh.Fields = make(map[string]fieldDoc)
	news := receiver.Auto[r].Alloc()
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
		t := fmt.Sprintf("%v", field.Type.Kind())
		typeName := field.Type.Name()
		typeString := field.Type.String()
		if typeName != "" {
			t = typeName
		} else if typeString != "" {
			t = typeString
		}
		fielddoc.Type = fmt.Sprintf("%s", t)
		if doc, ok := field.Tag.Lookup("doc"); ok {
			fielddoc.Doc = doc
			if ex, ok := field.Tag.Lookup("example"); ok {
				fielddoc.Example = ex
			}
		}
		sh.Fields[field.Name] = fielddoc
	}
	return sh, nil
}
