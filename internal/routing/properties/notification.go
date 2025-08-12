package properties

import "github.com/getkin/kin-openapi/openapi3"

var NotificationActionProperties = map[string]*openapi3.Schema{
	"link":     openapi3.NewStringSchema().WithFormat("uri"),
	"linkText": openapi3.NewStringSchema(),
}

var NotificationProperties = map[string]*openapi3.Schema{
	"id":           openapi3.NewStringSchema().WithFormat("uuid"),
	"title":        openapi3.NewStringSchema(),
	"message":      openapi3.NewStringSchema(),
	"action":       openapi3.NewSchema().WithProperties(NotificationActionProperties),
	"userId":       openapi3.NewStringSchema().WithFormat("uuid"),
	"closed":       openapi3.NewBoolSchema(),
	"modifiedById": openapi3.NewStringSchema().WithFormat("uuid"),
	"createdAt":    openapi3.NewDateTimeSchema(),
	"updatedAt":    openapi3.NewDateTimeSchema(),
}

var CreateNotificationPayloadProperties = map[string]*openapi3.Schema{
	"title":   openapi3.NewStringSchema(),
	"message": openapi3.NewStringSchema(),
	"action":  openapi3.NewSchema().WithProperties(NotificationActionProperties),
	"userId":  openapi3.NewStringSchema().WithFormat("uuid"),
	"closed":  openapi3.NewBoolSchema(),
}

var UpdateNotificationPayloadProperties = map[string]*openapi3.Schema{
	"title":   openapi3.NewStringSchema().WithNullable(),
	"message": openapi3.NewStringSchema().WithNullable(),
	"action":  openapi3.NewSchema().WithProperties(NotificationActionProperties).WithNullable(),
	"userId":  openapi3.NewStringSchema().WithFormat("uuid").WithNullable(),
	"closed":  openapi3.NewBoolSchema().WithNullable(),
}
