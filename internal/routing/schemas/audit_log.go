package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

// AuditLogSchema defines the OpenAPI schema for an audit log entry.
// It includes the standard audit log properties as specified in properties.AuditLogProperties,
// and additionally embeds a "user" property, which is an object described by properties.UserProperties.
// This schema is intended for use in API documentation and validation of audit log data structures.
var AuditLogSchema = openapi3.NewSchema().WithProperties(properties.AuditLogProperties).
	WithProperty(
		"user",
		UserSchema.Value,
	).
	WithRequired([]string{
		"id",
		"tableName",
		"objectId",
		"operationType",
		"data",
		"organizationId",
		"userId",
		"createdAt",
		"updatedAt",
	}).NewRef()

// AuditLogArraySchema defines an OpenAPI array schema for audit log entries.
// Each item in the array conforms to the AuditLogSchema specification.
// This schema is typically used to describe API responses that return multiple audit log records.
var AuditLogArraySchema = openapi3.NewArraySchema().
	WithItems(AuditLogSchema.Value).NewRef()
