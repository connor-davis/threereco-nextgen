package schemas

import (
	"github.com/connor-davis/threereco-nextgen/internal/routing/properties"
	"github.com/getkin/kin-openapi/openapi3"
)

var NotificationSchema = openapi3.NewSchema().WithProperties(properties.NotificationProperties).NewRef()

var NotificationArraySchema = openapi3.NewArraySchema().
	WithItems(NotificationSchema.Value).NewRef()

var CreateNotificationPayloadSchema = openapi3.NewSchema().WithProperties(properties.CreateNotificationPayloadProperties).NewRef()

var UpdateNotificationPayloadSchema = openapi3.NewSchema().WithProperties(properties.UpdateNotificationPayloadProperties).NewRef()
