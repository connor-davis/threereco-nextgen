package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

var ProductSchema = openapi3.NewSchema().
	WithProperties(properties.ProductProperties).
	WithProperty(
		"materials",
		MaterialArraySchema.Value,
	).
	WithProperty(
		"modifiedBy",
		ModifiedByUserSchema.Value,
	).
	WithRequired([]string{
		"id",
		"name",
		"value",
		"materials",
		"modifiedById",
		"createdAt",
		"updatedAt",
	}).NewRef()

var ProductArraySchema = openapi3.NewArraySchema().
	WithItems(ProductSchema.Value).NewRef()

var CreateProductPayloadSchema = openapi3.NewSchema().WithProperties(properties.CreateProductPayloadProperties).
	WithProperty(
		"materials",
		openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema().WithFormat("uuid")),
	).
	WithRequired([]string{
		"name",
		"value",
	}).NewRef()

var UpdateProductPayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateProductPayloadProperties).
	WithProperty(
		"materials",
		openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema().WithFormat("uuid")),
	).NewRef()
