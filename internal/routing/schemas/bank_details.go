package schemas

import "github.com/getkin/kin-openapi/openapi3"

var BankDetailsSchema = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.NewObjectSchema().Type,
		Properties: map[string]*openapi3.SchemaRef{
			"accountHolder": {
				Value: openapi3.NewStringSchema(),
			},
			"accountNumber": {
				Value: openapi3.NewStringSchema(),
			},
			"bankName": {
				Value: openapi3.NewStringSchema(),
			},
			"branchCode": {
				Value: openapi3.NewStringSchema(),
			},
		},
		Required: []string{
			"accountHolder",
			"accountNumber",
			"bankName",
			"branchCode",
		},
		Nullable: true,
	},
}
