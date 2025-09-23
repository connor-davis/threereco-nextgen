package schemas

import "github.com/getkin/kin-openapi/openapi3"

var PermissionSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"label": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
			},
			"value": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
			},
			"description": {
				Value: openapi3.NewStringSchema().WithFormat("text"),
			},
		},
		Required: []string{
			"label",
			"value",
			"description",
		},
	},
}

var PermissionsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/Permission",
		},
	},
}

var PermissionGroupSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"name": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
			},
			"permissions": {
				Ref: "#/components/schemas/Permissions",
			},
			"subGroups": {
				Value: &openapi3.Schema{
					Type: openapi3.NewArraySchema().Type,
					Items: &openapi3.SchemaRef{
						Ref: "#/components/schemas/PermissionGroup",
					},
				},
			},
		},
		Required: []string{
			"name",
			"permissions",
			"subGroups",
		},
	},
}

var PermissionGroupsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/PermissionGroup",
		},
	},
}
