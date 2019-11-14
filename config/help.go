package config

import (
	"fmt"
	"github.com/telenornms/skogul"
	"reflect"
	"unicode"
)

// FieldDoc is a structured representation of the documentation of a single
// field in a struct, used for modules
type FieldDoc struct {
	Doc     string
	Example string
	Type    string
}

// Help is the relevant help for a single module
type Help struct {
	Name        string
	Aliases     string
	Doc         string
	Fields      map[string]FieldDoc
	CustomTypes map[string]map[string]FieldDoc
}

// HelpModule looks up help for a module in the specified module map. It
// also fetches documentation for the struct fields, using reflection.
func HelpModule(mmap skogul.ModuleMap, mod string) (Help, error) {
	if mmap[mod] == nil {
		return Help{}, skogul.Error{Source: "config parser", Reason: "No such module"}
	}
	mh := Help{}
	mh.Name = mod
	mh.Doc = mmap[mod].Help
	for _, alias := range mmap[mod].Aliases {
		mh.Aliases = fmt.Sprintf("%s %s", alias, mh.Aliases)
	}
	mh.Fields, _ = getFieldDoc(mmap[mod].Alloc())
	mh.CustomTypes = make(map[string]map[string]FieldDoc)
	for _, extra := range mmap[mod].Extras {
		d, name := getFieldDoc(extra)
		mh.CustomTypes[name] = d
	}
	return mh, nil
}

func getFieldDoc(d interface{}) (map[string]FieldDoc, string) {
	fields := make(map[string]FieldDoc)
	st := reflect.TypeOf(d)
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if !unicode.IsUpper(rune(field.Name[0])) {
			continue
		}
		fielddoc := FieldDoc{}
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
		fields[field.Name] = fielddoc
	}
	return fields, st.Name()
}
