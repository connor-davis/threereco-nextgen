package schemas

import "github.com/getkin/kin-openapi/openapi3"

var MaterialSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"id": {
				Value: openapi3.NewUUIDSchema(),
			},
			"name": {
				Value: openapi3.NewStringSchema(),
			},
			"gwCode": {
				Value: openapi3.NewStringSchema(),
			},
			"carbonFactor": {
				Value: openapi3.NewStringSchema(),
			},
			"createdAt": {
				Value: openapi3.NewDateTimeSchema(),
			},
			"updatedAt": {
				Value: openapi3.NewDateTimeSchema(),
			},
		},
		Required: []string{
			"id",
			"name",
			"gwCode",
			"carbonFactor",
			"createdAt",
			"updatedAt",
		},
	},
}

var MaterialsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/Material",
		},
	},
}

var AssignMaterialSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"id": {
				Value: openapi3.NewUUIDSchema(),
			},
		},
		Required: []string{
			"id",
		},
	},
}

var AssignMaterialsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/AssignMaterial",
		},
	},
}

var CreateMaterialSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"name": {
				Value: openapi3.NewStringSchema(),
			},
			"gwCode": {
				Value: openapi3.NewStringSchema(),
			},
			"carbonFactor": {
				Value: openapi3.NewStringSchema(),
			},
		},
		Required: []string{
			"name",
			"gwCode",
			"carbonFactor",
		},
	},
}

var UpdateMaterialSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"name": {
				Value: openapi3.NewStringSchema().WithNullable(),
			},
			"gwCode": {
				Value: openapi3.NewStringSchema().WithNullable(),
			},
			"carbonFactor": {
				Value: openapi3.NewStringSchema().WithNullable(),
			},
		},
		Required: []string{
			"name",
			"gwCode",
			"carbonFactor",
		},
	},
}
