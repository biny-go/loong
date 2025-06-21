package add

import (
	"bytes"
	"strings"
	"text/template"
)

const protoTemplate = `
syntax = "proto3";

package {{.Package}};

option go_package = "{{.GoPackage}}";
option java_multiple_files = true;
option java_package = "{{.JavaPackage}}";

import "google/api/annotations.proto";

service {{.Service}} {
	rpc Save{{.Service}} (Save{{.Service}}Request) returns (Save{{.Service}}Reply) {
		option (google.api.http) = {
			post: "/{{.ServicePath}}"
			body: "*"
		};
	};
	rpc Delete{{.Service}} (Delete{{.Service}}Request) returns (Delete{{.Service}}Reply) {
		option (google.api.http) = {
			delete: "/{{.ServicePath}}"
		};
	};
	rpc Get{{.Service}} (Get{{.Service}}Request) returns (Get{{.Service}}Reply) {
		option (google.api.http) = {
			get: "/{{.ServicePath}}"
		};
	};
	rpc List{{.Service}} (List{{.Service}}Request) returns (List{{.Service}}Reply) {
		option (google.api.http) = {
			post: "/{{.ServicePath}}/list"
			body: "*"
		};
	};
}

message Save{{.Service}}Request {}
message Save{{.Service}}Reply {}

message Delete{{.Service}}Request {}
message Delete{{.Service}}Reply {}

message Get{{.Service}}Request {}
message Get{{.Service}}Reply {}

message List{{.Service}}Request {}
message List{{.Service}}Reply {}
`

func (p *Proto) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("proto").Parse(strings.TrimSpace(protoTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
