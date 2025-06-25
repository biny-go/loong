package server

import (
	"bytes"
	"html/template"
	"strings"
)

//nolint:lll
var controllerTemplate = `
{{- /* delete empty line */ -}}
package Controller

import (
	{{- if .UseContext }}
	"context"
	{{- end }}
	{{- if .UseIO }}
	"io"
	{{- end }}

	pb "{{ .Package }}"
	{{- if .GoogleEmpty }}
	"google.golang.org/protobuf/types/known/emptypb"
	{{- end }}

	common "github.com/biny-go/loong/proto/common"
	svc "{{ getFirstPart .Package }}/service"
)

type {{ .Service }}Controller struct {
	pb.Unimplemented{{ .Service }}Server
	s svc.{{ .Service }}Service

}

func New{{ .Service }}Controller() *{{ .Service }}Controller {
	return &{{ .Service }}Controller{}
}

{{- $s1 := "google.protobuf.Empty" }}
{{ range .Methods }}
{{- if eq .Type 1 }}
func (s *{{ .Service }}Controller) {{ .Name }}(ctx context.Context, req {{ if eq .Request $s1 }}*emptypb.Empty{{ else if hasPrefix .Request "common." }}*{{ .Request }}{{ else }}*pb.{{ .Request }}{{ end }}) ({{ if eq .Reply $s1 }}*emptypb.Empty{{ else if hasPrefix .Reply "common." }}*{{ .Reply }}{{ else }}*pb.{{ .Reply }}{{ end }}, error) {
	return {{ if eq .Reply $s1 }}&emptypb.Empty{}{{ else if hasPrefix .Reply "common." }}&{{ .Reply }}{}{{ else }}&pb.{{ .Reply }}{}{{ end }}, nil
}

{{- else if eq .Type 2 }}
func (s *{{ .Service }}Controller) {{ .Name }}(conn pb.{{ .Service }}_{{ .Name }}Server) error {
	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		
		err = conn.Send({{ if hasPrefix .Reply "common." }}&{{ .Reply }}{}{{ else }}&pb.{{ .Reply }}{}{{ end }})
		if err != nil {
			return err
		}
	}
}

{{- else if eq .Type 3 }}
func (s *{{ .Service }}Controller) {{ .Name }}(conn pb.{{ .Service }}_{{ .Name }}Server) error {
	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return conn.SendAndClose({{ if hasPrefix .Reply "common." }}&{{ .Reply }}{}{{ else }}&pb.{{ .Reply }}{}{{ end }})
		}
		if err != nil {
			return err
		}
	}
}

{{- else if eq .Type 4 }}
func (s *{{ .Service }}Controller) {{ .Name }}(req {{ if eq .Request $s1 }}*emptypb.Empty {{ else if hasPrefix .Request "common." }}*{{ .Request }}
{{ else }}*pb.{{ .Request }}{{ end }}, conn pb.{{ .Service }}_{{ .Name }}Server) error {
	for {
		err := conn.Send({{ if hasPrefix .Reply "common." }}&{{ .Reply }}{}{{ else }}&pb.{{ .Reply }}{}{{ end }})
		if err != nil {
			return err
		}
	}
}

{{- end }}
{{- end }}
`

type MethodType uint8

const (
	unaryType          MethodType = 1
	twoWayStreamsType  MethodType = 2
	requestStreamsType MethodType = 3
	returnsStreamsType MethodType = 4
)

// Service is a proto service.
type Service struct {
	Package     string
	Service     string
	Methods     []*Method
	GoogleEmpty bool

	UseIO      bool
	UseContext bool
}

// Method is a proto method.
type Method struct {
	Service string
	Name    string
	Request string
	Reply   string

	// type: unary or stream
	Type MethodType
}

func (s *Service) execute() ([]byte, error) {
	const empty = "google.protobuf.Empty"
	// 新增模板函数
	funcMap := template.FuncMap{
		"hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"getFirstPart": func(str string) string {
			parts := strings.SplitN(str, "/", 2)
			if len(parts) > 0 {
				return parts[0]
			}
			return ""
		},
	}
	buf := new(bytes.Buffer)
	for _, method := range s.Methods {
		if (method.Type == unaryType && (method.Request == empty || method.Reply == empty)) ||
			(method.Type == returnsStreamsType && method.Request == empty) {
			s.GoogleEmpty = true
		}
		if method.Type == twoWayStreamsType || method.Type == requestStreamsType {
			s.UseIO = true
		}
		if method.Type == unaryType {
			s.UseContext = true
		}
	}
	tmpl, err := template.New("service").Funcs(funcMap).Parse(controllerTemplate)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
