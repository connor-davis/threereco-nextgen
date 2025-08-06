package properties

import "github.com/getkin/kin-openapi/openapi3"

var LoginPayloadProperties = map[string]*openapi3.Schema{
	"Email":    openapi3.NewStringSchema().WithFormat("email"),
	"Password": openapi3.NewStringSchema().WithMinLength(6).WithMaxLength(100),
}
