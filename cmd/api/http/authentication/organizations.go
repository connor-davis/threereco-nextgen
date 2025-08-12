package authentication

import (
	"github.com/connor-davis/threereco-nextgen/internal/constants"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

// OrganizationsRoute defines the GET /authentication/organizations endpoint.
//
// Behavior:
//   - Requires authentication via the Authorized middleware.
//   - Retrieves the authenticated user from c.Locals("user").
//   - On success, returns the user's organizations as a JSON array ([]models.Organization).
//   - Returns 401 Unauthorized if no authenticated user is present.
//
// OpenAPI metadata:
//   - Summary: "Organizations"
//   - Description: "Return the user's organizations."
//   - Tags: ["Authentication"]
//
// Possible responses:
//   - 200 OK: JSON array of organizations.
//   - 401 Unauthorized: JSON error with "error" and "message" fields.
//   - 500 Internal Server Error: JSON error with "error" and "message" fields.
func (r *AuthenticationRouter) OrganizationsRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Successful organizations retrieval.").
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
			Summary:     "Organizations",
			Description: "Return the user's organizations.",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GetMethod,
		Path:   "/authentication/organizations",
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

			userOrganizations := []models.Organization{}

			for _, organization := range user.Organizations {
				userOrganizations = append(userOrganizations, organization)
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"items": userOrganizations,
			})
		},
	}
}
