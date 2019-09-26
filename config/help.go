package config

import (
	"fmt"
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"reflect"
	"strings"
	"unicode"
)

// Console width :D
const helpWidth = 66

// fieldDoc is a structured representation of the documentation of a single
// field in a struct, used for both senders and receivers (and more?)
type fieldDoc struct {
	Doc     string
	Example string
	Type    string
}

// Help is the relevant help for a single sender/receiver
type Help struct {
	Name   string
	Doc    string
	Fields map[string]fieldDoc
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
				fielddoc.Example = fmt.Sprintf("Example: %s", ex)
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
				fielddoc.Example = fmt.Sprintf("Example: %s", ex)
			}
		}
		sh.Fields[field.Name] = fielddoc
	}
	return sh, nil
}

/*
PrettyPrint is used to print a table with a header and wrapping the
description to fit a terminal nicely. Uses helpWidth to determine the size
of the "terminal".

Without PrettyPrint:

foo    | A very long line will be wrapped

With:

foo    | A very long
       | line will
       | be wrapped

We wrap at word boundaries to avoid splitting words.
*/
func PrettyPrint(scheme string, desc string) {
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

// Print uses PrettyPrint to output help for a sender or receiver.
func (sh Help) Print() {
	fmt.Printf("%s - %s\n", sh.Name, sh.Doc)
	fmt.Printf("Variables:\n")
	for n, f := range sh.Fields {
		t := ""
		t = fmt.Sprintf("[%s] ", f.Type)
		d := fmt.Sprintf("%s%s", t, f.Doc)
		PrettyPrint(n, d)
		if f.Example != "" {
			PrettyPrint("", "")
			PrettyPrint("", f.Example)
			PrettyPrint("", "")
		}
	}
}
