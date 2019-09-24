package config

import (
	"fmt"
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/sender"
	"reflect"
	"strings"
	"unicode"
)

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
