package bodies

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
)

var LoginPayloadBody = &openapi3.RequestBodyRef{
	Value: openapi3.NewRequestBody().WithJSONSchemaRef(schemas.LoginPayloadSchema).WithRequired(true),
}

var SignUpPayloadBody = &openapi3.RequestBodyRef{
	Value: openapi3.NewRequestBody().WithJSONSchemaRef(schemas.SignUpPayloadSchema).WithRequired(true),
}
