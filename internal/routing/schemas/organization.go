package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

// OrganizationSchema defines the OpenAPI schema for an organization entity.
// It includes the following properties:
//   - owner: An object representing the owner of the organization, with user properties.
//   - users: An array of user objects associated with the organization.
//   - roles: An array of role objects assigned within the organization.
//   - modifiedBy: An object representing the user who last modified the organization.
//
// The schema utilizes predefined property sets from the properties package for organization, user, and role details.
var OrganizationSchema = openapi3.NewSchema().WithProperties(properties.OrganizationProperties).
	WithProperty(
		"owner",
		openapi3.NewSchema().
			WithProperties(properties.UserProperties).
			WithRequired([]string{
				"id",
				"email",
				"mfaEnabled",
				"mfaVerified",
				"tags",
				"modifiedById",
				"createdAt",
				"updatedAt",
			}).NewRef().Value,
	).
	WithProperty(
		"modifiedBy",
		ModifiedByUserSchema.Value,
	).
	WithRequired([]string{
		"id",
		"name",
		"domain",
		"ownerId",
		"owner",
		"modifiedById",
		"createdAt",
		"updatedAt",
	}).NewRef()

// OrganizationArraySchema defines an OpenAPI schema reference for an array of Organization objects.
// It is constructed using the OrganizationSchema as the item type, allowing validation and documentation
// of API endpoints that return or accept lists of organizations.
var OrganizationArraySchema = openapi3.NewArraySchema().
	WithItems(OrganizationSchema.Value).NewRef()

var CreateOrganizationPayloadSchema = openapi3.NewSchema().WithProperties(properties.CreateOrganizationPayloadProperties).NewRef()

var UpdateOrganizationPayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateOrganizationPayloadProperties).NewRef()
