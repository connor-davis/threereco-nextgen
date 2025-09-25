package schemas

import "github.com/getkin/kin-openapi/openapi3"

var CollectionSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"id": {
				Value: openapi3.NewUUIDSchema(),
			},
			"materials": {
				Ref: "#/components/schemas/CollectionMaterials",
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

var CollectionsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/Collection",
		},
	},
}

var CreateCollectionSchema = &openapi3.SchemaRef{
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
				Ref: "#/components/schemas/AssignCollectionMaterials",
			},
			"createdAt": {
				Value: openapi3.NewDateTimeSchema(),
			},
		},
		Required: []string{
			"sellerId",
			"buyerId",
			"materials",
			"createdAt",
		},
	},
}

var UpdateCollectionSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"sellerId": {
				Value: openapi3.NewUUIDSchema().WithNullable(),
			},
			"buyerId": {
				Value: openapi3.NewUUIDSchema().WithNullable(),
			},
			"createdAt": {
				Value: openapi3.NewDateTimeSchema().WithNullable(),
			},
		},
		Required: []string{
			"sellerId",
			"buyerId",
			"createdAt",
		},
	},
}

var CollectionMaterialSchema = &openapi3.SchemaRef{
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

var CollectionMaterialsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/CollectionMaterial",
		},
	},
}

var AssignCollectionMaterialSchema = &openapi3.SchemaRef{
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

var AssignCollectionMaterialsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewArraySchema().Type,
		Items: &openapi3.SchemaRef{
			Ref: "#/components/schemas/AssignCollectionMaterial",
		},
	},
}

var CreateCollectionMaterialSchema = &openapi3.SchemaRef{
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

var UpdateCollectionMaterialSchema = &openapi3.SchemaRef{
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
