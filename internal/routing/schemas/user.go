package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

// UserSchema defines the OpenAPI schema for a User object.
// It includes the following properties:
//   - organizations: an array of Organization objects, each defined by OrganizationProperties.
//   - roles: an array of Role objects, each defined by RoleProperties.
//   - modifiedBy: a User object representing the user who last modified this record, defined by UserProperties.
//
// The base properties for the User object are provided by UserProperties.
var UserSchema = openapi3.NewSchema().WithProperties(properties.UserProperties).
	WithProperty(
		"organizations",
		openapi3.NewArraySchema().
			WithItems(openapi3.NewObjectSchema().
				WithProperties(properties.OrganizationProperties)),
	).
	WithProperty(
		"roles",
		openapi3.NewArraySchema().
			WithItems(openapi3.NewObjectSchema().
				WithProperties(properties.RoleProperties)),
	).
	WithProperty(
		"modifiedBy",
		openapi3.NewObjectSchema().
			WithProperties(properties.UserProperties),
	).NewRef()

// UserArraySchema defines an OpenAPI array schema for a collection of User objects.
// Each item in the array conforms to the UserSchema specification.
// This schema can be used to describe API responses or request bodies that return or accept
// multiple users in a single payload.
var UserArraySchema = openapi3.NewArraySchema().
	WithItems(UserSchema.Value).NewRef()

var CreateUserPayloadSchema = openapi3.NewSchema().WithProperties(properties.CreateUserPayloadProperties).NewRef()

var UpdateUserPayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateUserPayloadProperties).NewRef()
