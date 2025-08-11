package bodies

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
)

var CreateMaterialPayloadBody = &openapi3.RequestBodyRef{
	Value: openapi3.NewRequestBody().WithJSONSchemaRef(schemas.CreateMaterialPayloadSchema).WithRequired(true),
}

var UpdateMaterialPayloadBody = &openapi3.RequestBodyRef{
	Value: openapi3.NewRequestBody().WithJSONSchemaRef(schemas.UpdateMaterialPayloadSchema).WithRequired(true),
}
