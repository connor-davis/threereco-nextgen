package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

// AddressSchema defines the OpenAPI schema reference for an Address object.
// It includes standard address attributes plus metadata timestamps.
// Required fields:
//   - id: Unique identifier for the address record
//   - lineOne: Primary street address line
//   - city: City or locality
//   - state: State, province, or region
//   - zip: Postal or ZIP code
//   - country: ISO country designation
//   - createdAt: RFC 3339 timestamp when the record was created
//   - updatedAt: RFC 3339 timestamp when the record was last modified
//
// Additional optional fields (if present in AddressProperties) may include
// secondary address lines, coordinates, or validation metadata, but are not
// mandated by this schema. The schema is generated via openapi3 utilities
// and exposed as a reusable reference for route and component registration.
var AddressSchema = openapi3.NewSchema().
	WithProperties(properties.AddressProperties).
	WithRequired([]string{
		"id",
		"lineOne",
		"city",
		"state",
		"postalCode",
		"country",
		"createdAt",
		"updatedAt",
	}).NewRef()

// AddressArraySchema is an OpenAPI schema reference representing an array of Address objects.
// It wraps a base array schema whose item type is the underlying AddressSchema, allowing
// reuse anywhere a list (collection) of addresses is required in the API specification.
var AddressArraySchema = openapi3.NewArraySchema().
	WithItems(AddressSchema.Value).NewRef()

var CreateAddressPayloadSchema = openapi3.NewSchema().
	WithProperties(properties.CreateAddressPayloadProperties).
	WithRequired([]string{
		"lineOne",
		"city",
		"state",
		"postalCode",
		"country",
	}).NewRef()

var UpdateAddressPayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateAddressPayloadProperties).NewRef()
