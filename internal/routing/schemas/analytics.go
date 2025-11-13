package schemas

import "github.com/getkin/kin-openapi/openapi3"

var TopBusinessSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"name": {
				Value: openapi3.NewStringSchema(),
			},
			"weight": {
				Value: openapi3.NewFloat64Schema(),
			},
		},
		Required: []string{
			"name",
			"weight",
		},
	},
}

var MaterialTypeSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"type": {
				Value: openapi3.NewStringSchema(),
			},
			"weight": {
				Value: openapi3.NewFloat64Schema(),
			},
		},
		Required: []string{
			"type",
			"weight",
		},
	},
}
