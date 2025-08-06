package properties

import "github.com/getkin/kin-openapi/openapi3"

var MfaVerifyProperties = map[string]*openapi3.Schema{
	"code": openapi3.NewStringSchema().WithMin(6).WithMax(6),
}
