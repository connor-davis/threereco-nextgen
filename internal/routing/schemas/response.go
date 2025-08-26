package schemas

import "github.com/getkin/kin-openapi/openapi3"

// SuccessResponseSchema defines a reusable OpenAPI schema reference for successful responses.
// It includes two properties:
//   - "items": a oneOf schema, typically used for responses containing multiple items.
//   - "item": a oneOf schema, typically used for responses containing a single item.
// This schema can be used to standardize the structure of success responses in the API.
var SuccessResponseSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"items": openapi3.NewArraySchema().WithItems(
		openapi3.NewAnyOfSchema(
			UserSchema.Value,
			OrganizationSchema.Value,
			RoleSchema.Value,
			AuditLogSchema.Value,
			MaterialSchema.Value,
			ProductSchema.Value,
			TransactionSchema.Value,
			NotificationSchema.Value,
			AvailablePermissionGroupSchema.Value,
			AddressSchema.Value,
			BankDetailsSchema.Value,
			openapi3.NewStringSchema(),
		),
	),
	"item": openapi3.NewAnyOfSchema(
		UserSchema.Value,
		OrganizationSchema.Value,
		RoleSchema.Value,
		AuditLogSchema.Value,
		MaterialSchema.Value,
		ProductSchema.Value,
		TransactionSchema.Value,
		NotificationSchema.Value,
		AvailablePermissionGroupSchema.Value,
		AddressSchema.Value,
		BankDetailsSchema.Value,
	),
	"pageDetails": openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
		"count":        openapi3.NewIntegerSchema().WithMin(0),
		"nextPage":     openapi3.NewIntegerSchema().WithMin(2),
		"previousPage": openapi3.NewIntegerSchema().WithMin(1),
		"currentPage":  openapi3.NewIntegerSchema().WithMin(1),
		"pages":        openapi3.NewIntegerSchema().WithMin(0),
	}).WithNullable(),
}).NewRef()

// ErrorResponseSchema defines the OpenAPI schema for error responses.
// The schema includes two string properties:
//   - "error": A brief description of the error type, with a default value "Bad Request".
//   - "message": A detailed message explaining the error, with a default value
//     "The request could not be understood or was missing required parameters."
var ErrorResponseSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"error":   openapi3.NewStringSchema().WithDefault("Bad Request"),
	"message": openapi3.NewStringSchema().WithDefault("The request could not be understood or was missing required parameters."),
}).NewRef()
