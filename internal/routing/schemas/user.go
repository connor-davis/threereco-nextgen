package schemas

import "github.com/getkin/kin-openapi/openapi3"

var UserSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"id": {
				Value: openapi3.NewUUIDSchema(),
			},
			"name": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
			},
			"username": {
				Value: openapi3.NewStringSchema().
					WithPattern(`(?:\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b)|(?:\b(?:\+?\d{1,3}[-.\s]?)?(?:\(?\d{2,4}\)?[-.\s]?)?\d{3,4}[-.\s]?\d{3,4}\b)`),
			},
			"mfaEnabled": {
				Value: openapi3.NewBoolSchema(),
			},
			"mfaVerified": {
				Value: openapi3.NewBoolSchema(),
			},
			"permissions": {
				Value: openapi3.NewArraySchema().WithItems(
					openapi3.NewStringSchema().
						WithPattern(`^(\*|[a-zA-Z0-9]+(\.(\*|[a-zA-Z0-9]+))*)$`),
				),
			},
			"type": {
				Value: openapi3.NewStringSchema().WithEnum("system", "collector", "business"),
			},
			"businessId": {
				Value: openapi3.NewUUIDSchema().WithNullable(),
			},
			"address": {
				Ref: "#/components/schemas/Address",
			},
			"bankDetails": {
				Ref: "#/components/schemas/BankDetails",
			},
			"roles": {
				Ref: "#/components/schemas/Roles",
			},
			"businesses": {
				Value: &openapi3.Schema{
					Type: openapi3.NewArraySchema().Type,
					Items: &openapi3.SchemaRef{
						Ref: "#/components/schemas/Business",
					},
				},
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
			"username",
			"mfaEnabled",
			"mfaVerified",
			"permissions",
			"type",
			"businessId",
			"roles",
			"businesses",
			"address",
			"bankDetails",
			"createdAt",
			"updatedAt",
		},
	},
}

var UsersSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/User",
		},
	},
}

var AssignUserSchema = &openapi3.SchemaRef{
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

var AssignUsersSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/AssignUser",
		},
	},
}

var CreateUserSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"name": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
			},
			"username": {
				Value: openapi3.NewStringSchema().
					WithPattern(`(?:\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b)|(?:\b(?:\+?\d{1,3}[-.\s]?)?(?:\(?\d{2,4}\)?[-.\s]?)?\d{3,4}[-.\s]?\d{3,4}\b)`),
			},
			"permissions": {
				Value: openapi3.NewArraySchema().WithItems(
					openapi3.NewStringSchema().
						WithPattern(`^(\*|[a-zA-Z0-9]+(\.(\*|[a-zA-Z0-9]+))*)$`),
				),
			},
			"type": {
				Value: openapi3.NewStringSchema().WithEnum("system", "collector", "business").WithDefault("collector"),
			},
			"businessId": {
				Value: openapi3.NewUUIDSchema(),
			},
			"address": {
				Ref: "#/components/schemas/Address",
			},
			"bankDetails": {
				Ref: "#/components/schemas/BankDetails",
			},
			"roles": {
				Ref: "#/components/schemas/AssignRoles",
			},
			"businesses": {
				Ref: "#/components/schemas/AssignBusinesses",
			},
		},
		Required: []string{
			"name",
			"username",
			"permissions",
			"type",
			"address",
			"bankDetails",
			"roles",
			"businesses",
		},
	},
}

var UpdateUserSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"name": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3).WithNullable(),
			},
			"username": {
				Value: openapi3.NewStringSchema().
					WithPattern(`(?:\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b)|(?:\b(?:\+?\d{1,3}[-.\s]?)?(?:\(?\d{2,4}\)?[-.\s]?)?\d{3,4}[-.\s]?\d{3,4}\b)`).
					WithNullable(),
			},
			"permissions": {
				Value: openapi3.NewArraySchema().WithItems(
					openapi3.NewStringSchema().WithPattern(`^(\*|[a-zA-Z0-9]+(\.(\*|[a-zA-Z0-9]+))*)$`),
				).WithNullable(),
			},
			"type": {
				Value: openapi3.NewStringSchema().WithEnum("system", "collector", "business").WithDefault("collector").WithNullable(),
			},
			"businessId": {
				Value: openapi3.NewUUIDSchema().WithNullable(),
			},
			"address": {
				Ref: "#/components/schemas/Address",
			},
			"bankDetails": {
				Ref: "#/components/schemas/BankDetails",
			},
			"roles": {
				Ref: "#/components/schemas/AssignRoles",
			},
			"businesses": {
				Ref: "#/components/schemas/AssignBusinesses",
			},
		},
		Required: []string{
			"name",
			"username",
			"permissions",
			"type",
			"address",
			"bankDetails",
			"roles",
			"businesses",
		},
	},
}
