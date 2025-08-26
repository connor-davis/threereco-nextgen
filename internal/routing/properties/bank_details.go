package properties

import "github.com/getkin/kin-openapi/openapi3"

// BankDetailsProperties defines the OpenAPI schema components for a BankDetails resource.
// It maps logical field names to their corresponding OpenAPI 3 schema definitions:
//   - id: UUID string uniquely identifying the bank details record.
//   - accountNumber: Non-empty string representing the bank account number.
//   - accountHolder: Non-empty string with the full name of the account holder.
//   - bankName: Non-empty string naming the financial institution.
//   - branchCode: Non-empty string identifying the specific bank branch (e.g., routing / sort / branch code).
//   - createdAt: RFC 3339 timestamp indicating when the record was created.
//   - updatedAt: RFC 3339 timestamp indicating the last time the record was modified.
// Use this map when constructing OpenAPI component schemas or assembling request/response
// models to ensure consistency across routing, validation, and documentation layers.
var BankDetailsProperties = map[string]*openapi3.Schema{
	"id":            openapi3.NewStringSchema().WithFormat("uuid"),
	"accountNumber": openapi3.NewStringSchema().WithMinLength(1),
	"accountHolder": openapi3.NewStringSchema().WithMinLength(1),
	"bankName":      openapi3.NewStringSchema().WithMinLength(1),
	"branchCode":    openapi3.NewStringSchema().WithMinLength(1),
	"createdAt":     openapi3.NewDateTimeSchema(),
	"updatedAt":     openapi3.NewDateTimeSchema(),
}

var CreateBankDetailsPayloadProperties = map[string]*openapi3.Schema{
	"accountNumber": openapi3.NewStringSchema().WithMinLength(1),
	"accountHolder": openapi3.NewStringSchema().WithMinLength(1),
	"bankName":      openapi3.NewStringSchema().WithMinLength(1),
	"branchCode":    openapi3.NewStringSchema().WithMinLength(1),
}

var UpdateBankDetailsPayloadProperties = map[string]*openapi3.Schema{
	"accountNumber": openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
	"accountHolder": openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
	"bankName":      openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
	"branchCode":    openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
}
