package properties

import "github.com/getkin/kin-openapi/openapi3"

// OrganizationProperties defines the OpenAPI schema properties for an Organization object.
// Each property specifies its type, format, and validation constraints:
//   - "id": UUID string identifying the organization.
//   - "name": Non-empty string representing the organization's name.
//   - "domain": Non-empty string for the organization's domain.
//   - "ownerId": UUID string identifying the owner of the organization.
//   - "modifiedById": UUID string identifying the user who last modified the organization.
//   - "createdAt": DateTime indicating when the organization was created.
//   - "updatedAt": DateTime indicating when the organization was last updated.
var OrganizationProperties = map[string]*openapi3.Schema{
	"id":           openapi3.NewStringSchema().WithFormat("uuid"),
	"name":         openapi3.NewStringSchema().WithMinLength(1),
	"domain":       openapi3.NewStringSchema().WithMinLength(1),
	"ownerId":      openapi3.NewStringSchema().WithFormat("uuid"),
	"modifiedById": openapi3.NewStringSchema().WithFormat("uuid"),
	"createdAt":    openapi3.NewDateTimeSchema(),
	"updatedAt":    openapi3.NewDateTimeSchema(),
}

var CreateOrganizationPayloadProperties = map[string]*openapi3.Schema{
	"name":    openapi3.NewStringSchema().WithMinLength(1),
	"domain":  openapi3.NewStringSchema().WithMinLength(1),
	"ownerId": openapi3.NewStringSchema().WithFormat("uuid"),
}

var UpdateOrganizationPayloadProperties = map[string]*openapi3.Schema{
	"name":    openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
	"domain":  openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
	"ownerId": openapi3.NewStringSchema().WithFormat("uuid").WithNullable(),
}
