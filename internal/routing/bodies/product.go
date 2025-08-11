package bodies

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
)

var CreateProductPayloadBody = &openapi3.RequestBodyRef{
	Value: openapi3.NewRequestBody().WithJSONSchemaRef(schemas.CreateProductPayloadSchema).WithRequired(true),
}

var UpdateProductPayloadBody = &openapi3.RequestBodyRef{
	Value: openapi3.NewRequestBody().WithJSONSchemaRef(schemas.UpdateProductPayloadSchema).WithRequired(true),
}
