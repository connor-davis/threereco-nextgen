package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

// BankDetailsSchema defines the OpenAPI schema for a bank details resource.
//
// Required fields:
//   - id            : Unique identifier for the bank details record.
//   - accountHolder : Name of the individual or entity that owns the account.
//   - accountNumber : The bank account number (format per regional standards).
//   - bankName      : Name of the financial institution.
//   - branchCode    : Branch / routing / sort / transit code identifying the bank branch.
//   - createdAt     : ISO 8601 timestamp for when the record was created.
//   - updatedAt     : ISO 8601 timestamp for the last update to the record.
//
// The schema aggregates predefined BankDetailsProperties and enforces
// the above fields as required for validation in API requests and responses.
var BankDetailsSchema = openapi3.NewSchema().
	WithProperties(properties.BankDetailsProperties).
	WithRequired([]string{
		"id",
		"accountHolder",
		"accountNumber",
		"bankName",
		"branchCode",
		"createdAt",
		"updatedAt",
	}).NewRef()

// BankDetailsArraySchema defines the reusable OpenAPI schema reference for an array of bank detail objects,
// where each element conforms to BankDetailsSchema. It is used wherever a list of bank account details is
// required in the API specification (e.g., responses returning multiple bank records or request bodies that
// submit batches). By referencing the item schema (BankDetailsSchema.Value) it ensures consistency, enables
// validation of each element, and reduces duplication in the generated OpenAPI document.
var BankDetailsArraySchema = openapi3.NewArraySchema().
	WithItems(BankDetailsSchema.Value).NewRef()

var CreateBankDetailsPayloadSchema = openapi3.NewSchema().
	WithProperties(properties.CreateBankDetailsPayloadProperties).
	WithRequired([]string{
		"accountHolder",
		"accountNumber",
		"bankName",
		"branchCode",
	}).NewRef()

var UpdateBankDetailsPayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateBankDetailsPayloadProperties).NewRef()
