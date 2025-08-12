package properties

import "github.com/getkin/kin-openapi/openapi3"

var TransactionProperties = map[string]*openapi3.Schema{
	"id":             openapi3.NewStringSchema().WithFormat("uuid"),
	"type":           openapi3.NewStringSchema(),
	"weight":         openapi3.NewFloat64Schema().WithMin(0.0),
	"amount":         openapi3.NewFloat64Schema().WithMin(0.0),
	"sellerAccepted": openapi3.NewBoolSchema(),
	"sellerDeclined": openapi3.NewBoolSchema(),
	"sellerId":       openapi3.NewStringSchema().WithFormat("uuid"),
	"sellerType":     openapi3.NewStringSchema(),
	"buyerId":        openapi3.NewStringSchema().WithFormat("uuid"),
	"buyerType":      openapi3.NewStringSchema(),
	"modifiedById":   openapi3.NewStringSchema().WithFormat("uuid"),
	"createdAt":      openapi3.NewDateTimeSchema(),
	"updatedAt":      openapi3.NewDateTimeSchema(),
}

var CreateTransactionPayloadProperties = map[string]*openapi3.Schema{
	"type":     openapi3.NewStringSchema(),
	"weight":   openapi3.NewFloat64Schema().WithMin(0.0),
	"amount":   openapi3.NewFloat64Schema().WithMin(0.0),
	"sellerId": openapi3.NewStringSchema().WithFormat("uuid"),
	"buyerId":  openapi3.NewStringSchema().WithFormat("uuid"),
}

var UpdateTransactionPayloadProperties = map[string]*openapi3.Schema{
	"type":           openapi3.NewStringSchema().WithNullable(),
	"weight":         openapi3.NewFloat64Schema().WithMin(0.0).WithNullable(),
	"amount":         openapi3.NewFloat64Schema().WithMin(0.0).WithNullable(),
	"sellerId":       openapi3.NewStringSchema().WithFormat("uuid").WithNullable(),
	"buyerId":        openapi3.NewStringSchema().WithFormat("uuid").WithNullable(),
	"sellerAccepted": openapi3.NewBoolSchema().WithNullable(),
	"sellerDeclined": openapi3.NewBoolSchema().WithNullable(),
}
