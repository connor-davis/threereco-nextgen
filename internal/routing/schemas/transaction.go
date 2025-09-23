package schemas

import "github.com/getkin/kin-openapi/openapi3"

var TransactionSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"id": {
				Value: openapi3.NewUUIDSchema(),
			},
			"materials": {
				Ref: "#/components/schemas/TransactionMaterials",
			},
			"seller": {
				Ref: "#/components/schemas/User",
			},
			"buyer": {
				Ref: "#/components/schemas/Business",
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
			"materials",
			"seller",
			"buyer",
			"createdAt",
			"updatedAt",
		},
	},
}

var TransactionsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/Transaction",
		},
	},
}

var CreateTransactionSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"sellerId": {
				Value: openapi3.NewUUIDSchema(),
			},
			"buyerId": {
				Value: openapi3.NewUUIDSchema(),
			},
			"materials": {
				Ref: "#/components/schemas/AssignTransactionMaterials",
			},
		},
		Required: []string{
			"sellerId",
			"buyerId",
			"materials",
		},
	},
}

var UpdateTransactionSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"sellerId": {
				Value: openapi3.NewUUIDSchema().WithNullable(),
			},
			"buyerId": {
				Value: openapi3.NewUUIDSchema().WithNullable(),
			},
		},
		Required: []string{
			"sellerId",
			"buyerId",
		},
	},
}

var TransactionMaterialSchema = &openapi3.SchemaRef{
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
			"weight": {
				Value: openapi3.NewFloat64Schema().WithMin(0),
			},
			"value": {
				Value: openapi3.NewFloat64Schema().WithMin(0),
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
			"material",
			"weight",
			"value",
			"createdAt",
			"updatedAt",
		},
	},
}

var TransactionMaterialsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/TransactionMaterial",
		},
	},
}

var AssignTransactionMaterialSchema = &openapi3.SchemaRef{
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
			"weight": {
				Value: openapi3.NewFloat64Schema().WithMin(0),
			},
			"value": {
				Value: openapi3.NewFloat64Schema().WithMin(0),
			},
		},
		Required: []string{
			"name",
			"gwCode",
			"carbonFactor",
			"weight",
			"value",
		},
	},
}

var AssignTransactionMaterialsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/AssignTransactionMaterial",
		},
	},
}

var CreateTransactionMaterialSchema = &openapi3.SchemaRef{
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
			"weight": {
				Value: openapi3.NewFloat64Schema().WithMin(0),
			},
			"value": {
				Value: openapi3.NewFloat64Schema().WithMin(0),
			},
		},
		Required: []string{
			"name",
			"gwCode",
			"carbonFactor",
			"weight",
			"value",
		},
	},
}

var UpdateTransactionMaterialSchema = &openapi3.SchemaRef{
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
			"weight": {
				Value: openapi3.NewFloat64Schema().WithMin(0).WithNullable(),
			},
			"value": {
				Value: openapi3.NewFloat64Schema().WithMin(0).WithNullable(),
			},
		},
		Required: []string{
			"name",
			"gwCode",
			"carbonFactor",
			"weight",
			"value",
		},
	},
}
