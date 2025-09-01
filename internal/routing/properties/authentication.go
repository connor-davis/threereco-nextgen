package properties

import "github.com/getkin/kin-openapi/openapi3"

var LoginPayloadProperties = map[string]*openapi3.Schema{
	"emailOrPhone": openapi3.NewStringSchema().WithFormat("email"),
	"password":     openapi3.NewStringSchema().WithMinLength(6).WithMaxLength(100),
}

var SignUpPayloadProperties = map[string]*openapi3.Schema{
	"name":     openapi3.NewStringSchema().WithMinLength(2).WithMaxLength(100),
	"email":    openapi3.NewStringSchema().WithFormat("email"),
	"phone":    openapi3.NewStringSchema().WithMinLength(10).WithMaxLength(15),
	"password": openapi3.NewStringSchema().WithMinLength(6).WithMaxLength(100),
}
