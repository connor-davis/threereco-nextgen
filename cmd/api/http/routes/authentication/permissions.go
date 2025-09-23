package authentication

import (
	"github.com/connor-davis/threereco-nextgen/common"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func (r *AuthenticationRouter) PermissionsRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Permissions retrieved successfully.").
			WithContent(openapi3.Content{
				"text/plain": openapi3.NewMediaType().
					WithSchemaRef(schemas.SuccessSchema),
			}),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Bad Request").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(schemas.ErrorSchema),
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Unauthorized").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(schemas.ErrorSchema),
			}),
	})

	responses.Set("403", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Forbidden").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(schemas.ErrorSchema),
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Internal Server Error").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(schemas.ErrorSchema),
			}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     "Permissions",
			Description: "Retrieves the permissions available in the system.",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GET,
		Path:   "/authentication/permissions",
		Middlewares: []fiber.Handler{
			r.middleware.Authenticated(),
			r.middleware.Authorized("permissions.view"),
		},
		Handler: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"items": common.Permissions,
			})
		},
	}
}
