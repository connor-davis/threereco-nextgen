package bodies

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
)

var CreateUserPayloadBody = &openapi3.RequestBodyRef{
	Value: openapi3.NewRequestBody().WithJSONSchemaRef(schemas.CreateUserPayloadSchema).WithRequired(true),
}
