package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

// RoleSchema defines the OpenAPI schema for a Role object.
// It includes the following properties:
//   - users: an array of user objects, each defined by UserProperties.
//   - organizations: an array of organization objects, each defined by OrganizationProperties.
//   - modifiedBy: a user object representing the last user who modified the role, defined by UserProperties.
//
// The base properties for the role are provided by RoleProperties.
var RoleSchema = openapi3.NewSchema().
	WithProperties(properties.RoleProperties).
	WithProperty(
		"modifiedBy",
		ModifiedByUserSchema.Value,
	).
	WithRequired([]string{
		"id",
		"name",
		"description",
		"permissions",
		"modifiedById",
		"createdAt",
		"updatedAt",
	}).NewRef()

var AvailablePermissionSchema = openapi3.NewSchema().
	WithProperties(properties.AvailablePermissionProperties).
	WithRequired([]string{
		"value",
		"description",
	}).
	NewRef()

var AvailablePermissionGroupSchema = openapi3.NewSchema().
	WithProperties(properties.AvailablePermissionsGroupProperties).
	WithProperty(
		"permissions",
		openapi3.NewArraySchema().WithItems(AvailablePermissionSchema.Value),
	).
	WithRequired([]string{
		"name",
		"permissions",
	}).
	NewRef()

// RoleArraySchema defines an OpenAPI array schema for roles, where each item in the array
// conforms to the RoleSchema specification. This schema can be used to describe API endpoints
// that accept or return a list of roles.
var RoleArraySchema = openapi3.NewArraySchema().
	WithItems(RoleSchema.Value).NewRef()

var AvailablePermissionsGroupArraySchema = openapi3.NewArraySchema().
	WithItems(AvailablePermissionGroupSchema.Value).NewRef()

var CreateRolePayloadSchema = openapi3.NewSchema().WithProperties(properties.CreateRolePayloadProperties).NewRef()

var UpdateRolePayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateRolePayloadProperties).NewRef()
