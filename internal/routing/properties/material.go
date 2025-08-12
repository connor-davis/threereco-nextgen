package properties

import "github.com/getkin/kin-openapi/openapi3"

var MaterialProperties = map[string]*openapi3.Schema{
	"id":           openapi3.NewStringSchema().WithFormat("uuid"),
	"name":         openapi3.NewStringSchema(),
	"gwCode":       openapi3.NewStringSchema(),
	"carbonFactor": openapi3.NewStringSchema(),
	"modifiedById": openapi3.NewStringSchema().WithFormat("uuid"),
	"createdAt":    openapi3.NewDateTimeSchema(),
	"updatedAt":    openapi3.NewDateTimeSchema(),
}

var CreateMaterialPayloadProperties = map[string]*openapi3.Schema{
	"name":         openapi3.NewStringSchema(),
	"gwCode":       openapi3.NewStringSchema(),
	"carbonFactor": openapi3.NewStringSchema(),
}

var UpdateMaterialPayloadProperties = map[string]*openapi3.Schema{
	"name":         openapi3.NewStringSchema().WithNullable(),
	"gwCode":       openapi3.NewStringSchema().WithNullable(),
	"carbonFactor": openapi3.NewStringSchema().WithNullable(),
}
