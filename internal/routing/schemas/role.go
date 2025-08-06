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
var RoleSchema = openapi3.NewSchema().WithProperties(properties.RoleProperties).
	WithProperty(
		"users",
		openapi3.NewArraySchema().
			WithItems(openapi3.NewObjectSchema().
				WithProperties(properties.UserProperties)),
	).
	WithProperty(
		"organizations",
		openapi3.NewArraySchema().
			WithItems(openapi3.NewObjectSchema().
				WithProperties(properties.OrganizationProperties)),
	).
	WithProperty(
		"modifiedBy",
		openapi3.NewObjectSchema().
			WithProperties(properties.UserProperties),
	).NewRef()

// RoleArraySchema defines an OpenAPI array schema for roles, where each item in the array
// conforms to the RoleSchema specification. This schema can be used to describe API endpoints
// that accept or return a list of roles.
var RoleArraySchema = openapi3.NewArraySchema().
	WithItems(RoleSchema.Value).NewRef()
