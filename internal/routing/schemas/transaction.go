package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

var TransactionSchema = openapi3.NewSchema().
	WithProperties(properties.TransactionProperties).
	WithProperty(
		"products",
		ProductSchema.Value,
	).
	WithProperty(
		"modifiedBy",
		ModifiedByUserSchema.Value,
	).
	WithRequired([]string{
		"id",
		"type",
		"weight",
		"amount",
		"sellerAccepted",
		"sellerDeclined",
		"sellerId",
		"sellerType",
		"buyerId",
		"buyerType",
		"products",
		"modifedById",
		"createdAt",
		"updatedAt",
	}).NewRef()

var TransactionArraySchema = openapi3.NewArraySchema().
	WithItems(TransactionSchema.Value).NewRef()

var CreateTransactionPayloadSchema = openapi3.NewSchema().WithProperties(properties.CreateTransactionPayloadProperties).
	WithProperty(
		"products",
		openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema().WithFormat("uuid")),
	).NewRef()

var UpdateTransactionPayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateTransactionPayloadProperties).
	WithProperty(
		"products",
		openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema().WithFormat("uuid")),
	).NewRef()
