package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

var MaterialSchema = openapi3.NewSchema().WithProperties(properties.MaterialProperties).NewRef()

var MaterialArraySchema = openapi3.NewArraySchema().
	WithItems(MaterialSchema.Value).NewRef()

var CreateMaterialPayloadSchema = openapi3.NewSchema().WithProperties(properties.CreateMaterialPayloadProperties).NewRef()

var UpdateMaterialPayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateMaterialPayloadProperties).NewRef()
