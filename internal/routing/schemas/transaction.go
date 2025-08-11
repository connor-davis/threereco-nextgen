package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

var TransactionSchema = openapi3.NewSchema().WithProperties(properties.TransactionProperties).
	WithProperty(
		"products",
		openapi3.NewArraySchema().WithItems(ProductSchema.Value),
	).NewRef()

var TransactionArraySchema = openapi3.NewArraySchema().
	WithItems(TransactionSchema.Value).NewRef()

var CreateTransactionPayloadSchema = openapi3.NewSchema().WithProperties(properties.CreateTransactionPayloadProperties).NewRef()

var UpdateTransactionPayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateTransactionPayloadProperties).
	WithProperty(
		"products",
		openapi3.NewArraySchema().WithItems(ProductSchema.Value),
	).NewRef()
