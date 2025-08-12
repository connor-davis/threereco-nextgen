package properties

import "github.com/getkin/kin-openapi/openapi3"

// RoleProperties defines the OpenAPI schema properties for a Role object.
// Each key in the map corresponds to a property of the Role, with its value specifying
// the OpenAPI schema for that property:
//   - "id":           UUID string uniquely identifying the role.
//   - "name":         Non-empty string representing the name of the role.
//   - "description":  Nullable string providing a description of the role.
//   - "permissions":  Array of strings specifying permissions associated with the role.
//   - "modifiedById": UUID string identifying the user who last modified the role.
//   - "createdAt":    DateTime indicating when the role was created.
//   - "updatedAt":    DateTime indicating when the role was last updated.
var RoleProperties = map[string]*openapi3.Schema{
	"id":           openapi3.NewStringSchema().WithFormat("uuid"),
	"name":         openapi3.NewStringSchema().WithMinLength(1),
	"description":  openapi3.NewStringSchema(),
	"permissions":  openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema()).WithDefault([]string{}),
	"modifiedById": openapi3.NewStringSchema().WithFormat("uuid"),
	"createdAt":    openapi3.NewDateTimeSchema(),
	"updatedAt":    openapi3.NewDateTimeSchema(),
}

var AvailablePermissionProperties = map[string]*openapi3.Schema{
	"value":       openapi3.NewStringSchema(),
	"description": openapi3.NewStringSchema(),
}

var AvailablePermissionsGroupProperties = map[string]*openapi3.Schema{
	"name": openapi3.NewStringSchema().WithMinLength(1),
}

var CreateRolePayloadProperties = map[string]*openapi3.Schema{
	"name":        openapi3.NewStringSchema().WithMinLength(1),
	"description": openapi3.NewStringSchema().WithNullable(),
	"permissions": openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema()).WithDefault([]string{}),
}

var UpdateRolePayloadProperties = map[string]*openapi3.Schema{
	"name":        openapi3.NewStringSchema().WithMinLength(1).WithNullable(),
	"description": openapi3.NewStringSchema().WithNullable(),
	"permissions": openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema()).WithDefault([]string{}),
}
