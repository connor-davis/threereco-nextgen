package properties

import "github.com/getkin/kin-openapi/openapi3"

// AuditLogProperties defines the OpenAPI schema properties for an audit log entry.
// The map keys represent the property names, and the values specify their respective schema definitions:
//   - "id":            UUID string identifying the audit log entry.
//   - "tableName":     Name of the database table affected (non-empty string).
//   - "operationType": Type of operation performed (e.g., "CREATE", "UPDATE", "DELETE"; non-empty string).
//   - "objectId":      UUID string identifying the affected object.
//   - "data":          Serialized data related to the operation (non-empty string).
//   - "userId":        UUID string identifying the user who performed the operation.
//   - "createdAt":     Timestamp when the audit log entry was created.
//   - "updatedAt":     Timestamp when the audit log entry was last updated.
var AuditLogProperties = map[string]*openapi3.Schema{
	"id":            openapi3.NewStringSchema().WithFormat("uuid"),
	"tableName":     openapi3.NewStringSchema().WithMinLength(1),
	"operationType": openapi3.NewStringSchema().WithMinLength(1),
	"objectId":      openapi3.NewStringSchema().WithFormat("uuid"),
	"data":          openapi3.NewObjectSchema(),
	"userId":        openapi3.NewStringSchema().WithFormat("uuid"),
	"createdAt":     openapi3.NewDateTimeSchema(),
	"updatedAt":     openapi3.NewDateTimeSchema(),
}
