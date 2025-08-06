package properties

import "github.com/getkin/kin-openapi/openapi3"

// UserProperties defines the OpenAPI schema properties for a User object.
// Each key represents a user attribute and maps to its corresponding OpenAPI schema definition.
// The properties include:
//   - "id": Unique identifier for the user (UUID format).
//   - "email": User's email address (email format, required).
//   - "image": URL or reference to the user's profile image (nullable).
//   - "name": User's full name (nullable).
//   - "phone": User's phone number (nullable).
//   - "jobTitle": User's job title (nullable).
//   - "mfaEnabled": Indicates if multi-factor authentication is enabled (default: false).
//   - "mfaVerified": Indicates if multi-factor authentication is verified (default: false).
//   - "primaryOrganizationId": Identifier for the user's primary organization (UUID format, nullable).
//   - "modifiedById": Identifier for the user who last modified this record (UUID format).
//   - "createdAt": Timestamp when the user was created (date-time format).
//   - "updatedAt": Timestamp when the user was last updated (date-time format).
var UserProperties = map[string]*openapi3.Schema{
	"id":                    openapi3.NewStringSchema().WithFormat("uuid"),
	"email":                 openapi3.NewStringSchema().WithFormat("email").WithMinLength(1),
	"image":                 openapi3.NewStringSchema().WithNullable(),
	"name":                  openapi3.NewStringSchema().WithNullable(),
	"phone":                 openapi3.NewStringSchema().WithNullable(),
	"jobTitle":              openapi3.NewStringSchema().WithNullable(),
	"mfaEnabled":            openapi3.NewBoolSchema().WithDefault(false),
	"mfaVerified":           openapi3.NewBoolSchema().WithDefault(false),
	"primaryOrganizationId": openapi3.NewStringSchema().WithFormat("uuid").WithNullable(),
	"modifiedById":          openapi3.NewStringSchema().WithFormat("uuid"),
	"createdAt":             openapi3.NewDateTimeSchema(),
	"updatedAt":             openapi3.NewDateTimeSchema(),
}

var CreateUserPayloadProperties = map[string]*openapi3.Schema{
	"email":    openapi3.NewStringSchema().WithFormat("email").WithMinLength(1),
	"password": openapi3.NewStringSchema().WithMinLength(8),
	"name":     openapi3.NewStringSchema().WithNullable(),
	"phone":    openapi3.NewStringSchema().WithNullable(),
	"jobTitle": openapi3.NewStringSchema().WithNullable(),
	"roles":    openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema().WithFormat("uuid")),
}
