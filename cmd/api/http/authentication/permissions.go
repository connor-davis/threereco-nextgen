package authentication

import (
	"slices"

	"github.com/connor-davis/threereco-nextgen/internal/constants"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

// PermissionsRoute defines the GET /authentication/permissions endpoint for AuthenticationRouter.
// The route requires authorization (via the Authorized middleware) and reads the authenticated
// user from the Fiber context (Locals("user")). If the user is missing or invalid, it responds
// with HTTP 401 and an error payload. On success, it collects all non-empty permissions from the
// user's roles, de-duplicates them, and returns the resulting list as a JSON array with HTTP 200.
// The route also provides OpenAPI metadata (summary, description, tags) and documents responses
// for 200 (success), 401 (unauthorized), and 500 (internal server error).
func (r *AuthenticationRouter) PermissionsRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Successful permissions retrieval.").
			WithJSONSchema(schemas.SuccessResponseSchema.Value).
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.SuccessResponseSchema.Value).
					WithExample("example", map[string]any{
						"item": schemas.UserSchema,
					}),
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorResponseSchema.Value).
			WithDescription(string(constants.UnauthorizedError)).
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorResponseSchema.Value).
					WithExample("example", map[string]any{
						"error":   string(constants.UnauthorizedError),
						"message": string(constants.UnauthorizedErrorDetails),
					}),
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorResponseSchema.Value).
			WithDescription("Internal server error.").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorResponseSchema.Value).
					WithExample("example", map[string]any{
						"error":   string(constants.InternalServerError),
						"message": string(constants.InternalServerErrorDetails),
					}),
			}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     "Permissions",
			Description: "Return the user's permissions.",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GetMethod,
		Path:   "/authentication/permissions",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			user, ok := c.Locals("user").(*models.User)

			if !ok || user == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   constants.UnauthorizedError,
					"message": constants.UnauthorizedErrorDetails,
				})
			}

			userPermissions := []string{}

			for _, role := range user.Roles {
				for _, permission := range role.Permissions {
					if permission == "" {
						continue
					}

					if slices.Contains(userPermissions, permission) {
						continue
					}

					userPermissions = append(userPermissions, permission)
				}
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"items": userPermissions,
			})
		},
	}
}
