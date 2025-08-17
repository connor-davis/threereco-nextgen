package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

var MaterialSchema = openapi3.NewSchema().
	WithProperties(properties.MaterialProperties).
	WithProperty(
		"modifiedBy",
		ModifiedByUserSchema.Value,
	).
	WithRequired([]string{
		"id",
		"name",
		"gwCode",
		"carbonFactor",
		"products",
		"modifiedById",
		"createdAt",
		"updatedAt",
	}).NewRef()

var MaterialArraySchema = openapi3.NewArraySchema().
	WithItems(MaterialSchema.Value).NewRef()

var CreateMaterialPayloadSchema = openapi3.NewSchema().WithProperties(properties.CreateMaterialPayloadProperties).WithRequired([]string{
	"name",
	"gwCode",
	"carbonFactor",
}).NewRef()

var UpdateMaterialPayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateMaterialPayloadProperties).NewRef()
