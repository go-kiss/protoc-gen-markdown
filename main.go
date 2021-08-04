package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func main() {
	g := markdown{}

	var flags flag.FlagSet

	flags.StringVar(&g.Prefix, "prefix", "/", "API path prefix")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(g.Generate)
}

type markdown struct {
	Prefix string
}

func (md *markdown) Generate(plugin *protogen.Plugin) error {
	for _, f := range plugin.Files {
		if len(f.Services) == 0 {
			continue
		}

		fname := f.GeneratedFilenamePrefix + ".md"
		t := plugin.NewGeneratedFile(fname, f.GoImportPath)

		printDoc := func(c protogen.Comments) {
			t.P(string(c))
		}

		for _, s := range f.Services {
			t.P("# ", s.Desc.Name())
			t.P()
			printDoc(s.Comments.Leading)

			for _, m := range s.Methods {
				name := string(m.Desc.FullName())
				api := md.api(name)
				anchor := md.anchor(api)

				t.P(fmt.Sprintf("- [%s](#%s)", api, anchor))
			}
			t.P()
			for _, m := range s.Methods {
				n := string(m.Desc.FullName())
				t.P("## ", md.api(n))
				t.P()
				printDoc(m.Comments.Leading)
				t.P()
				t.P("### Request")
				t.P("```javascript")
				t.P(md.jsDocForMessage(m.Input))
				t.P("```")
				t.P()
				t.P("### Reply")
				t.P("```javascript")
				t.P(md.jsDocForMessage(m.Output))
				t.P("```")
			}
		}

		t.P()
	}
	return nil
}

func (md *markdown) api(s string) string {
	i := strings.LastIndex(s, ".")

	prefix := strings.Trim(md.Prefix, "/")
	if prefix != "" {
		prefix = "/" + prefix
	}

	return prefix + "/" + s[:i] + "/" + s[i+1:]
}

func (md *markdown) anchor(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "/", "")
	return s
}

func (md *markdown) scalarDefaultValue(field *protogen.Field) string {
	switch field.Desc.Kind() {
	case protoreflect.StringKind, protoreflect.BytesKind:
		return `""`
	case protoreflect.Fixed64Kind, protoreflect.Int64Kind,
		protoreflect.Sfixed64Kind, protoreflect.Sint64Kind,
		protoreflect.Uint64Kind:
		return `"0"`
	case protoreflect.DoubleKind, protoreflect.FloatKind:
		return `0.0`
	case protoreflect.BoolKind:
		return "false"
	default:
		return "0"
	}
}

func (md *markdown) jsDocForField(field *protogen.Field) string {
	js := field.Comments.Leading.String()
	js += string(field.Desc.Name()) + ":"

	var vv string
	var vt string
	if field.Desc.IsMap() {
		vf := field.Message.Fields[1]
		if m := vf.Message; m != nil {
			vv = md.jsDocForMessage(m)
			vt = string(vf.Message.Desc.FullName())
		} else {
			vv = md.scalarDefaultValue(vf)
			vt = vf.Desc.Kind().String()
		}
		kf := field.Desc.MapKey()
		vv = fmt.Sprintf("{\n\"%s\":%s}", kf.Default().String(), vv)
		vt = fmt.Sprintf("%s,%s", kf.Kind().String(), vt)
	} else if field.Message != nil {
		vv = md.jsDocForMessage(field.Message)
		vt = string(field.Message.Desc.Name())
	} else if field.Enum != nil {
		vv = `"ENUM"`
	} else {
		vv = md.scalarDefaultValue(field)
		vt = field.Desc.Kind().String()
	}

	if field.Desc.IsList() {
		js += "[\n" + vv + "]" + fmt.Sprintf(", // list<%s>", vt)
	} else if field.Desc.IsMap() {
		js += vv + fmt.Sprintf(", // map<%s>", vt)
	} else {
		js += vv + fmt.Sprintf(", // type<%s>", vt)
	}

	if t := string(field.Comments.Trailing); len(t) > 0 {
		js += ", " + strings.TrimLeft(t, " ")
	} else {
		js += "\n"
	}

	return js
}

func (md *markdown) jsDocForMessage(m *protogen.Message) string {
	js := "{\n"

	for _, field := range m.Fields {
		js += md.jsDocForField(field)
	}

	js += "}"
	options := jsbeautifier.DefaultOptions()
	js, _ = jsbeautifier.Beautify(&js, options)

	return js
}
