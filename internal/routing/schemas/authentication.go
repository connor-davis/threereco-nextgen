package schemas

import "github.com/getkin/kin-openapi/openapi3"

var LoginPayloadSchema = &openapi3.RequestBodyRef{
	Value: &openapi3.RequestBody{
		Description: "Login payload",
		Content: openapi3.Content{
			"application/json": openapi3.NewMediaType().
				WithSchema(&openapi3.Schema{
					Type: openapi3.NewObjectSchema().Type,
					Properties: map[string]*openapi3.SchemaRef{
						"username": {
							Value: openapi3.NewStringSchema().
								WithPattern(`^(?:\+?1)?[2-9]\d{2}[2-9](?!11)\d{6}$|^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).
								WithDefault(nil),
						},
						"password": {
							Value: openapi3.NewStringSchema().
								WithMinLength(8),
						},
					},
					Required: []string{
						"username",
						"password",
					},
				}),
		},
		Required: true,
	},
}

var RegisterPayloadSchema = &openapi3.RequestBodyRef{
	Value: &openapi3.RequestBody{
		Description: "Register payload",
		Content: openapi3.Content{
			"application/json": openapi3.NewMediaType().
				WithSchema(&openapi3.Schema{
					Type: openapi3.NewObjectSchema().Type,
					Properties: map[string]*openapi3.SchemaRef{
						"name": {
							Value: openapi3.NewStringSchema().WithFormat("text").WithMin(3),
						},
						"username": {
							Value: openapi3.NewStringSchema().
								WithPattern(`^(?:\+?1)?[2-9]\d{2}[2-9](?!11)\d{6}$|^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).
								WithDefault(nil),
						},
						"password": {
							Value: openapi3.NewStringSchema().
								WithMinLength(8),
						},
						"type": {
							Value: openapi3.NewStringSchema().WithEnum("business", "collector"),
						},
					},
					Required: []string{
						"name",
						"username",
						"password",
						"type",
					},
				}),
		},
		Required: true,
	},
}

var VerifyMfaPayloadSchema = &openapi3.RequestBodyRef{
	Value: &openapi3.RequestBody{
		Description: "Verify MFA payload",
		Content: openapi3.Content{
			"application/json": openapi3.NewMediaType().
				WithSchema(&openapi3.Schema{
					Type: openapi3.NewObjectSchema().Type,
					Properties: map[string]*openapi3.SchemaRef{
						"code": {
							Value: openapi3.NewStringSchema().WithMinLength(6).WithMaxLength(6),
						},
					},
					Required: []string{
						"code",
					},
				}),
		},
		Required: true,
	},
}
