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
		fields := make([]field, len(md.GetMessageFields()))
		for i, fd := range md.GetMessageFields() {
			typeName := fd.GetTypeName()
			if typeName == "" {
				typeName = fd.GetType().String()
			}

			fields[i] = field{
				Name:  fd.GetName(),
				Type:  typeName,
				Doc:   fd.GetComments().GetLeading(),
				Note:  fd.GetComments().GetTrailing(),
				Label: fd.GetLabel(),
			}
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
	Name  string
	Type  string
	Note  string
	Doc   string
	Label descriptor.FieldDescriptorProto_Label
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
		api.Doc = md.GetComments().GetLeading()

		inputType := md.GetInputType()[1:] // trim leading dot
		api.Request = t.messages[inputType]

		outputType := md.GetOutputType()[1:] // trim leading dot
		api.Reply = t.messages[outputType]

		t.apis = append(t.apis, &api)
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
		if field.isRepeated() {
			v = "[\"\",\"\"]"
		} else {
			v = "\"\""
		}
		vt = "string"
	} else if field.Type == "TYPE_DOUBLE" || field.Type == "TYPE_FLOAT" {
		if field.isRepeated() {
			v = "[0.0, 0.0]"
		} else {
			v = "0.0"
		}
		vt = "float"
	} else if field.Type == "TYPE_BOOL" {
		if field.isRepeated() {
			v = "[false, false]"
		} else {
			v = "false"
		}
		vt = "bool"
	} else if field.Type == "TYPE_INT64" || field.Type == "TYPE_UINT64" {
		if field.isRepeated() {
			v = "[\"0\", \"0\"]"
		} else {
			v = "\"0\""
		}
		vt = "string(int64)"
	} else if field.Type == "TYPE_INT32" || field.Type == "TYPE_UINT32" {
		if field.isRepeated() {
			v = "[0, 0]"
		} else {
			v = "0"
		}
		vt = "int"
	} else if field.Type[0] == '.' {
		m := t.messages[field.Type[1:]]
		v = t.generateJsDocForMessage(m)
		if field.isRepeated() {
			doc := "// type:<list>"
			if field.Note != "" {
				doc = " " + field.Note
			}
			v = "[" + doc + "\n" + v + "]"
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
	comment := strings.Split(t.comments.Leading, "\n")
	for _, value := range comment {
		t.P("### ", value)
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
