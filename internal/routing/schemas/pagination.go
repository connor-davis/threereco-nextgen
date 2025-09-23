package schemas

import "github.com/getkin/kin-openapi/openapi3"

var PaginationSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"count": {
				Value: openapi3.NewIntegerSchema().WithDefault(0),
			},
			"pages": {
				Value: openapi3.NewIntegerSchema().WithDefault(1),
			},
			"pageSize": {
				Value: openapi3.NewIntegerSchema().WithDefault(10),
			},
			"currentPage": {
				Value: openapi3.NewIntegerSchema().WithDefault(1),
			},
			"nextPage": {
				Value: openapi3.NewIntegerSchema().WithDefault(1),
			},
			"previousPage": {
				Value: openapi3.NewIntegerSchema().WithDefault(1),
			},
		},
		Required: []string{
			"count",
			"pages",
			"pageSize",
			"currentPage",
			"nextPage",
			"previousPage",
		},
	},
}
