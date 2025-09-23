package schemas

import "github.com/getkin/kin-openapi/openapi3"

var RoleSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"id": {
				Value: openapi3.NewUUIDSchema(),
			},
			"name": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
			},
			"description": {
				Value: openapi3.NewStringSchema().WithFormat("text"),
			},
			"permissions": {
				Value: openapi3.NewArraySchema().WithItems(
					openapi3.NewStringSchema().
						WithPattern(`^(\*|[a-zA-Z0-9]+(\.(\*|[a-zA-Z0-9]+))*)$`),
				),
			},
			"default": {
				Value: openapi3.NewBoolSchema(),
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
			"description",
			"permissions",
			"default",
			"createdAt",
			"updatedAt",
		},
	},
}

var RolesSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/Role",
		},
	},
}

var AssignRoleSchema = &openapi3.SchemaRef{
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

var AssignRolesSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/AssignRole",
		},
	},
}

var CreateRoleSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"name": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
			},
			"description": {
				Value: openapi3.NewStringSchema().WithFormat("text"),
			},
			"permissions": {
				Value: openapi3.NewArraySchema().WithItems(
					openapi3.NewStringSchema().
						WithPattern(`^(\*|[a-zA-Z0-9]+(\.(\*|[a-zA-Z0-9]+))*)$`),
				),
			},
		},
		Required: []string{
			"name",
			"description",
			"permissions",
		},
	},
}

var UpdateRoleSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"name": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3).WithNullable(),
			},
			"description": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithNullable(),
			},
			"permissions": {
				Value: openapi3.NewArraySchema().WithItems(
					openapi3.NewStringSchema().WithPattern(`^(\*|[a-zA-Z0-9]+(\.(\*|[a-zA-Z0-9]+))*)$`),
				).WithNullable(),
			},
		},
		Required: []string{
			"name",
			"description",
			"permissions",
		},
	},
}
