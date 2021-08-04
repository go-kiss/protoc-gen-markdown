package main

import (
	"fmt"
	// "os"
	// "strconv"
	"strings"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	// "github.com/k0kubun/pp/v3"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type twirp struct{}

func newGenerator() *twirp {
	t := &twirp{}

	return t
}

func (t *twirp) Generate(plugin *protogen.Plugin) error {
	for _, f := range plugin.Files {
		if len(f.Services) == 0 {
			continue
		}

		fname := f.GeneratedFilenamePrefix + ".md"
		gf := plugin.NewGeneratedFile(fname, f.GoImportPath)

		for _, s := range f.Services {
			for _, m := range s.Methods {
				gf.P(t.jsDocForMessage(m.Input))
				// gf.P(t.generateJsDocForMessage(m.Output))
			}
		}
	}
	return nil
}

func (t *twirp) scalarDefaultValue(field *protogen.Field) string {
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

func (t *twirp) jsDocForField(field *protogen.Field) string {
	js := field.Comments.Leading.String()
	js += `"` + string(field.Desc.Name()) + `":`

	var vv string
	var vt string
	if field.Desc.IsMap() {
		vf := field.Message.Fields[1]
		if m := vf.Message; m != nil {
			vv = t.jsDocForMessage(m)
			vt = string(vf.Message.Desc.FullName())
		} else {
			vv = t.scalarDefaultValue(vf)
			vt = vf.Desc.Kind().String()
		}
		kf := field.Desc.MapKey()
		vv = fmt.Sprintf("{\n\"%s\":%s}", kf.Default().String(), vv)
		vt = fmt.Sprintf("%s,%s", kf.Kind().String(), vt)
	} else if field.Message != nil {
		vv = t.jsDocForMessage(field.Message)
		vt = string(field.Message.Desc.Name())
	} else if field.Enum != nil {
		vv = `"ENUM"`
	} else {
		vv = t.scalarDefaultValue(field)
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

func (t *twirp) jsDocForMessage(m *protogen.Message) string {
	js := "{\n"

	for _, field := range m.Fields {
		js += t.jsDocForField(field)
	}

	js += "}"
	options := jsbeautifier.DefaultOptions()
	js, _ = jsbeautifier.Beautify(&js, options)

	return js
}

// func (t *twirp) generateDoc() {
// 	options := jsbeautifier.DefaultOptions()
// 	t.P("# ", t.name)
// 	t.P()
// 	comments := strings.Split(t.comments.Leading, "\n")
// 	for _, value := range comments {
// 		t.P(value, "  ")
// 	}
// 	t.P()
// 	for _, api := range t.apis {
// 		anchor := strings.Replace(api.Path, "/", "", -1)
// 		anchor = strings.Replace(anchor, ".", "", -1)
// 		anchor = strings.ToLower(anchor)
//
// 		t.P(fmt.Sprintf("- [%s](#%s)", api.Path, anchor))
// 	}
//
// 	t.P()
//
// 	for _, api := range t.apis {
// 		t.P("## ", api.Path)
// 		t.P()
// 		t.P(api.Doc)
// 		t.P()
// 		t.P("### Method")
// 		t.P()
// 		t.P(api.Method)
// 		t.P()
// 		t.P("### Request")
// 		t.P("```javascript")
// 		code, _ := jsbeautifier.Beautify(&api.Input, options)
// 		t.P(code)
// 		t.P("```")
// 		t.P()
// 		t.P("### Reply")
// 		t.P("```javascript")
// 		code, _ = jsbeautifier.Beautify(&api.Output, options)
// 		t.P(code)
// 		t.P("```")
// 	}
// }
