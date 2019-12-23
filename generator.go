package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"
)

type twirp struct {
	params commandLineParams

	// Output buffer that holds the bytes we want to write out for a single file.
	// Gets reset after working on a file.
	output *bytes.Buffer

	// Map of all proto messages
	messages map[string]*message

	// List of all APIs
	apis []*api

	// List of all service comments
	comments *protokit.Comment

	// Service name
	name string
}

func newGenerator(params commandLineParams) *twirp {
	t := &twirp{
		params:   params,
		messages: map[string]*message{},
		apis:     []*api{},
		output:   bytes.NewBuffer(nil),
	}

	return t
}

func (t *twirp) Generate(in *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	resp := new(plugin.CodeGeneratorResponse)

	t.scanAllMessages(in, resp)
	t.GenerateMarkdown(in, resp)

	return resp
}

// P forwards to g.gen.P, which prints output.
func (t *twirp) P(args ...string) {
	for _, v := range args {
		t.output.WriteString(v)
	}
	t.output.WriteByte('\n')
}

func (t *twirp) scanAllMessages(req *plugin.CodeGeneratorRequest, resp *plugin.CodeGeneratorResponse) {
	descriptors := protokit.ParseCodeGenRequest(req)

	for _, d := range descriptors {
		t.scanMessages(d)
	}
}

func (t *twirp) GenerateMarkdown(req *plugin.CodeGeneratorRequest, resp *plugin.CodeGeneratorResponse) {
	descriptors := protokit.ParseCodeGenRequest(req)

	for _, d := range descriptors {
		for _, sd := range d.GetServices() {
			t.scanService(sd)
			t.name = *sd.Name
			for _, api := range t.apis {
				api.Input = t.generateJsDocForMessage(api.Request)
				api.Output = t.generateJsDocForMessage(api.Reply)
			}

			t.generateDoc()

			name := strings.Replace(d.GetName(), ".proto", ".md", 1)
			resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(name),
				Content: proto.String(t.output.String()),
			})
		}
	}
}

func (t *twirp) scanMessages(d *protokit.FileDescriptor) {
	for _, md := range d.GetMessages() {
		t.scanMessage(md)
	}
}

func (t *twirp) scanMessage(md *protokit.Descriptor) {
	for _, smd := range md.GetMessages() {
		t.scanMessage(smd)
	}

	{
		fields := make([]field, len(md.GetMessageFields()))

		maps := make(map[string]*descriptor.DescriptorProto)
		for _, t := range md.NestedType {
			if t.Options.GetMapEntry() {
				pkg := md.GetPackage()
				name := fmt.Sprintf(".%s.%s.%s", pkg, md.GetName(), t.GetName())
				maps[name] = t
			}
		}

		for i, fd := range md.GetMessageFields() {
			typeName := fd.GetTypeName()
			if typeName == "" {
				typeName = fd.GetType().String()
			}

			f := field{
				Name:  fd.GetName(),
				Type:  typeName,
				Doc:   fd.GetComments().GetLeading(),
				Note:  fd.GetComments().GetTrailing(),
				Label: fd.GetLabel(),
			}

			if m, ok := maps[f.Type]; ok {
				for _, ff := range m.GetField() {
					switch ff.GetName() {
					case "key":
						f.KeyType = ff.GetType().String()
					case "value":
						typeName := ff.GetTypeName()
						if typeName == "" {
							typeName = ff.GetType().String()
						}
						f.Type = typeName
					}
				}
				f.Label = 0
			}
			fields[i] = f
		}

		t.messages[md.GetFullName()] = &message{
			Name:   md.GetName(),
			Doc:    md.GetComments().GetTrailing(),
			Fields: fields,
		}
	}
}

type message struct {
	Name   string
	Fields []field
	Label  descriptor.FieldDescriptorProto_Label
	Doc    string
}

type field struct {
	Name    string
	Type    string
	KeyType string
	Note    string
	Doc     string
	Label   descriptor.FieldDescriptorProto_Label
}

func (f field) isRepeated() bool {
	return f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

type api struct {
	Method  string
	Path    string
	Doc     string
	Request *message
	Reply   *message
	Input   string
	Output  string
}

func (t *twirp) scanService(d *protokit.ServiceDescriptor) {
	t.comments = d.Comments
	for _, md := range d.GetMethods() {
		api := api{}

		api.Method = "POST"
		api.Path = t.params.pathPrefix + "/" + d.GetFullName() + "/" + md.GetName()
		doc := md.GetComments().GetLeading()
		// 支持文档换行
		api.Doc = strings.Replace(doc, "\n", "\n\n", -1)

		inputType := md.GetInputType()[1:] // trim leading dot
		api.Request = t.messages[inputType]

		outputType := md.GetOutputType()[1:] // trim leading dot
		api.Reply = t.messages[outputType]

		t.apis = append(t.apis, &api)
	}
}

func getType(t string) string {
	switch t {
	case "TYPE_STRING":
		return "string"
	case "TYPE_DOUBLE", "TYPE_FLOAT":
		return "float"
	case "TYPE_BOOL":
		return "bool"
	case "TYPE_INT64", "TYPE_UINT64", "TYPE_INT32", "TYPE_UINT32":
		return "int"
	default:
		return t
	}
}

func getTypeValue(t string) string {
	switch t {
	case "TYPE_STRING":
		return ""
	case "TYPE_DOUBLE", "TYPE_FLOAT":
		return "0.0"
	case "TYPE_BOOL":
		return "false"
	case "TYPE_INT64", "TYPE_UINT64", "TYPE_INT32", "TYPE_UINT32":
		return "0"
	default:
		return ""
	}
}

func (t *twirp) generateJsDocForField(field field) string {
	var js string
	var v, vt string
	disableDoc := false

	if field.Doc != "" {
		for _, line := range strings.Split(field.Doc, "\n") {
			js += "// " + line + "\n"
		}
	}

	if field.Type == "TYPE_STRING" {
		vt = "string"
		if field.isRepeated() {
			v = `["",""]`
		} else if field.KeyType != "" {
			v = fmt.Sprintf(`{"%s":""}`, getTypeValue(field.KeyType))
			vt = fmt.Sprintf("map<%s,string>", getType(field.KeyType))
		} else {
			v = `""`
		}
	} else if field.Type == "TYPE_DOUBLE" || field.Type == "TYPE_FLOAT" {
		vt = "float"
		if field.isRepeated() {
			v = "[0.0, 0.0]"
		} else if field.KeyType != "" {
			v = fmt.Sprintf(`{"%s":0.0}`, getTypeValue(field.KeyType))
			vt = fmt.Sprintf("map<%s,float>", getType(field.KeyType))
		} else {
			v = "0.0"
		}
	} else if field.Type == "TYPE_BOOL" {
		vt = "bool"
		if field.isRepeated() {
			v = "[false, false]"
		} else if field.KeyType != "" {
			v = fmt.Sprintf(`{"%s":false}`, getTypeValue(field.KeyType))
			vt = fmt.Sprintf("map<%s,bool>", getType(field.KeyType))
		} else {
			v = "false"
		}
	} else if field.Type == "TYPE_INT64" || field.Type == "TYPE_UINT64" {
		vt = "string(int64)"
		if field.isRepeated() {
			v = `["0", "0"]`
		} else if field.KeyType != "" {
			v = fmt.Sprintf(`{"%s":"0"}`, getTypeValue(field.KeyType))
			vt = fmt.Sprintf("map<%s,string<int64>>", getType(field.KeyType))
		} else {
			v = `"0"`
		}
	} else if field.Type == "TYPE_INT32" || field.Type == "TYPE_UINT32" {
		vt = "int"
		if field.isRepeated() {
			v = "[0, 0]"
		} else if field.KeyType != "" {
			v = fmt.Sprintf(`{"%s":0}`, getTypeValue(field.KeyType))
			vt = fmt.Sprintf("map<%s,int>", getType(field.KeyType))
		} else {
			v = "0"
		}
	} else if field.Type[0] == '.' {
		m := t.messages[field.Type[1:]]
		v = t.generateJsDocForMessage(m)
		if field.isRepeated() {
			doc := fmt.Sprintf("// type:<list<%s>>", m.Name)
			if field.Note != "" {
				doc = " " + field.Note
			}
			v = "[" + doc + "\n" + v + "]"
		} else if field.KeyType != "" {
			doc := fmt.Sprintf("// type:<map<%s,%s>>", getType(field.KeyType), m.Name)
			if field.Note != "" {
				doc = " " + field.Note
			}
			v = fmt.Sprintf("{%s\n\"%s\":%s}", doc, getTypeValue(field.KeyType), v)
		}
		disableDoc = true
	} else {
		v = "UNKNOWN"
	}

	if disableDoc {
		js += fmt.Sprintf("%s: %s,", field.Name, v)
	} else {
		js += fmt.Sprintf("%s: %s, // type:<%s>", field.Name, v, vt)
		if field.Note != "" {
			js = js + ", " + field.Note
		}
	}
	js = strings.Trim(js, " ")

	js += "\n"

	return js
}

func (t *twirp) generateJsDocForMessage(m *message) string {
	var js string
	js += "{\n"

	for _, field := range m.Fields {
		js += t.generateJsDocForField(field)
	}

	js += "}"

	return js
}

func (t *twirp) generateDoc() {
	options := jsbeautifier.DefaultOptions()
	t.P("# ", t.name)
	t.P()
	comments := strings.Split(t.comments.Leading, "\n")
	for _, value := range comments {
		t.P(value, "  ")
	}
	t.P()
	for _, api := range t.apis {
		anchor := strings.Replace(api.Path, "/", "", -1)
		anchor = strings.Replace(anchor, ".", "", -1)
		anchor = strings.ToLower(anchor)

		t.P(fmt.Sprintf("- [%s](#%s)", api.Path, anchor))
	}

	t.P()

	for _, api := range t.apis {
		t.P("## ", api.Path)
		t.P()
		t.P(api.Doc)
		t.P()
		t.P("### Method")
		t.P()
		t.P(api.Method)
		t.P()
		t.P("### Request")
		t.P("```javascript")
		code, _ := jsbeautifier.Beautify(&api.Input, options)
		t.P(code)
		t.P("```")
		t.P()
		t.P("### Reply")
		t.P("```javascript")
		code, _ = jsbeautifier.Beautify(&api.Output, options)
		t.P(code)
		t.P("```")
	}
}
