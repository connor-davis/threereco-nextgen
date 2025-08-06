package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

var MfaVerifyPayloadSchema = openapi3.NewSchema().WithProperties(properties.MfaVerifyProperties).NewRef()
