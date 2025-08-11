package properties

import "github.com/getkin/kin-openapi/openapi3"

var ProductProperties = map[string]*openapi3.Schema{
	"id":        openapi3.NewStringSchema().WithFormat("uuid"),
	"name":      openapi3.NewStringSchema(),
	"value":     openapi3.NewFloat64Schema().WithMin(0.0),
	"createdAt": openapi3.NewDateTimeSchema(),
	"updatedAt": openapi3.NewDateTimeSchema(),
}

var CreateProductPayloadProperties = map[string]*openapi3.Schema{
	"name":  openapi3.NewStringSchema(),
	"value": openapi3.NewFloat64Schema().WithMin(0.0),
}

var UpdateProductPayloadProperties = map[string]*openapi3.Schema{
	"name":  openapi3.NewStringSchema().WithNullable(),
	"value": openapi3.NewFloat64Schema().WithMin(0.0).WithNullable(),
}
