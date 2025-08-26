package properties

import "github.com/getkin/kin-openapi/openapi3"

// AddressProperties defines the reusable OpenAPI schema components for an Address resource.
// It maps each field name to an openapi3.Schema describing validation rules and formats.
// Fields:
//   - id: string, format uuid (unique identifier for the address)
//   - lineOne: string, required (minimum length 1)
//   - lineTwo: string, nullable (optional second address line)
//   - city: string, required (minimum length 1)
//   - state: string, required (minimum length 1)
//   - zip: integer (postal code; stored as an integer rather than string)
//   - country: string, required (minimum length 1; expected to hold ISO country name or code)
//   - createdAt: string, date-time format (timestamp when the address was created)
//   - updatedAt: string, date-time format (timestamp of the last update)
//
// Note: Required vs optional enforcement occurs at the schema usage site (e.g., when composing object schemas)
// by listing required property names; this map only provides the individual property schemas.
var AddressProperties = map[string]*openapi3.Schema{
	"id":        openapi3.NewStringSchema().WithFormat("uuid"),
	"lineOne":   openapi3.NewStringSchema().WithMinLength(1),
	"lineTwo":   openapi3.NewStringSchema().WithNullable(),
	"city":      openapi3.NewStringSchema().WithMinLength(1),
	"state":     openapi3.NewStringSchema().WithMinLength(1),
	"zip":       openapi3.NewIntegerSchema(),
	"country":   openapi3.NewStringSchema().WithMinLength(1),
	"createdAt": openapi3.NewDateTimeSchema(),
	"updatedAt": openapi3.NewDateTimeSchema(),
}

var CreateAddressPayloadProperties = map[string]*openapi3.Schema{
	"lineOne": openapi3.NewStringSchema().WithMinLength(1),
	"lineTwo": openapi3.NewStringSchema().WithNullable(),
	"city":    openapi3.NewStringSchema().WithMinLength(1),
	"state":   openapi3.NewStringSchema().WithMinLength(1),
	"zip":     openapi3.NewIntegerSchema(),
	"country": openapi3.NewStringSchema().WithMinLength(1),
}

var UpdateAddressPayloadProperties = map[string]*openapi3.Schema{
	"lineOne": openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
	"lineTwo": openapi3.NewStringSchema().WithNullable(),
	"city":    openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
	"state":   openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
	"zip":     openapi3.NewIntegerSchema().WithNullable(),
	"country": openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
}
