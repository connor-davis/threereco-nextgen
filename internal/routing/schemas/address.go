package schemas

import "github.com/getkin/kin-openapi/openapi3"

var AddressSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"lineOne": {
				Value: openapi3.NewStringSchema(),
			},
			"lineTwo": {
				Value: openapi3.NewStringSchema().WithNullable(),
			},
			"city": {
				Value: openapi3.NewStringSchema(),
			},
			"province": {
				Value: openapi3.NewStringSchema(),
			},
			"country": {
				Value: openapi3.NewStringSchema(),
			},
			"zipCode": {
				Value: openapi3.NewStringSchema(),
			},
		},
		Required: []string{
			"lineOne",
			"lineTwo",
			"city",
			"province",
			"country",
			"zipCode",
		},
		Nullable: true,
	},
}
