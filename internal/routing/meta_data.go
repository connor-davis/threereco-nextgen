package routing

import "github.com/getkin/kin-openapi/openapi3"

// OpenAPIMetadata encapsulates metadata for an OpenAPI operation.
// It includes summary and description fields, tags for categorization,
// parameter references, request body information, and response definitions.
type OpenAPIMetadata struct {
	Summary     string
	Description string
	Tags        []string
	Parameters  []*openapi3.ParameterRef
	RequestBody *openapi3.RequestBodyRef
	Responses   *openapi3.Responses
}
