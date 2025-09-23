package schemas

import "github.com/getkin/kin-openapi/openapi3"

var QuerySchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"page": {
				Value: openapi3.NewIntegerSchema().WithDefault(1),
			},
			"pageSize": {
				Value: openapi3.NewIntegerSchema().WithDefault(10),
			},
			"searchTerm": {
				Value: openapi3.NewStringSchema().WithDefault(""),
			},
		},
		Required: []string{
			"page",
			"pageSize",
			"searchTerm",
		},
	},
}
