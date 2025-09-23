package schemas

import "github.com/getkin/kin-openapi/openapi3"

var BusinessSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"id": {
				Value: openapi3.NewUUIDSchema(),
			},
			"name": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
			},
			"ownerId": {
				Value: openapi3.NewUUIDSchema(),
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
			"users": {
				Value: &openapi3.Schema{
					Type: openapi3.NewArraySchema().Type,
					Items: &openapi3.SchemaRef{
						Ref: "#/components/schemas/User",
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
			"ownerId",
			"roles",
			"users",
			"address",
			"bankDetails",
			"createdAt",
			"updatedAt",
		},
	},
}

var BusinessesSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/Business",
		},
	},
}

var AssignBusinessSchema = &openapi3.SchemaRef{
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

var AssignBusinessesSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/AssignBusiness",
		},
	},
}

var CreateBusinessSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"name": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
			},
			"ownerId": {
				Value: openapi3.NewUUIDSchema(),
			},
			"address": {
				Ref: "#/components/schemas/Address",
			},
			"bankDetails": {
				Ref: "#/components/schemas/BankDetails",
			},
		},
		Required: []string{
			"name",
			"ownerId",
			"address",
			"bankDetails",
		},
	},
}

var UpdateBusinessSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"name": {
				Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
			},
			"ownerId": {
				Value: openapi3.NewUUIDSchema(),
			},
			"address": {
				Ref: "#/components/schemas/Address",
			},
			"bankDetails": {
				Ref: "#/components/schemas/BankDetails",
			},
		},
		Required: []string{
			"name",
			"ownerId",
			"address",
			"bankDetails",
		},
	},
}
