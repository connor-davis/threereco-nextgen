package bodies

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
)

var CreateTransactionPayloadBody = &openapi3.RequestBodyRef{
	Value: openapi3.NewRequestBody().WithJSONSchemaRef(schemas.CreateTransactionPayloadSchema).WithRequired(true),
}

var UpdateTransactionPayloadBody = &openapi3.RequestBodyRef{
	Value: openapi3.NewRequestBody().WithJSONSchemaRef(schemas.UpdateTransactionPayloadSchema).WithRequired(true),
}
