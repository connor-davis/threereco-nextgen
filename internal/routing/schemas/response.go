package schemas

import (
	"github.com/getkin/kin-openapi/openapi3"
)

var ErrorSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"error": {
				Value: openapi3.NewStringSchema().WithFormat("text"),
			},
			"message": {
				Value: openapi3.NewStringSchema().WithFormat("text"),
			},
		},
		Required: []string{
			"error",
			"message",
		},
	},
}

var SuccessSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		OneOf: []*openapi3.SchemaRef{
			{
				Value: &openapi3.Schema{
					Type: openapi3.NewObjectSchema().Type,
					Properties: map[string]*openapi3.SchemaRef{
						"item": {
							Value: &openapi3.Schema{
								AnyOf: []*openapi3.SchemaRef{
									UserSchema,
									RoleSchema,
									MaterialSchema,
									CollectionSchema,
									CollectionMaterialSchema,
									TransactionSchema,
									TransactionMaterialSchema,
									BusinessSchema,
								},
							},
						},
						"items": {
							Value: &openapi3.Schema{
								Type: openapi3.NewArraySchema().Type,
								Items: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										AnyOf: []*openapi3.SchemaRef{
											UserSchema,
											RoleSchema,
											MaterialSchema,
											CollectionSchema,
											CollectionMaterialSchema,
											TransactionSchema,
											TransactionMaterialSchema,
											BusinessSchema,
											PermissionGroupSchema,
										},
									},
								},
							},
						},
						"pagination": {
							Ref: "#/components/schemas/Pagination",
						},
					},
					Required: []string{},
				},
			},
		},
	},
}
