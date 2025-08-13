package invites

import (
	"github.com/connor-davis/threereco-nextgen/internal/constants"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (r *OrganizationsInvitesRouter) AcceptInvite() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The invitation has been successfully accepted."),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Invalid request.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"error":   constants.BadRequestError,
						"message": constants.BadRequestErrorDetails,
					},
					Schema: schemas.ErrorResponseSchema,
				},
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().WithJSONSchema(
			schemas.ErrorResponseSchema.Value,
		).WithDescription("Unauthorized.").WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"error":   constants.UnauthorizedError,
					"message": constants.UnauthorizedErrorDetails,
				},
				Schema: schemas.ErrorResponseSchema,
			},
		}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().WithJSONSchema(
			schemas.ErrorResponseSchema.Value,
		).WithDescription("Internal Server Error.").WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				},
				Schema: schemas.ErrorResponseSchema,
			},
		}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     "Accept Organization Invite",
			Description: "Accepts an invitation to join the organization.",
			Tags:        []string{"Organizations"},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.PostMethod,
		Path:   "/organizations/{id}/invites/accept",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(*models.User)

			organizationId := c.Params("id")

			existingOrganization, err := r.Services.Organizations.GetById(uuid.MustParse(organizationId))

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			if existingOrganization == nil || existingOrganization.Id == uuid.Nil {
				return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
					"error":   constants.NotFoundError,
					"message": constants.NotFoundErrorDetails,
				})
			}

			if currentUser.PrimaryOrganizationId != nil && currentUser.PrimaryOrganizationId.String() == organizationId {
				return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
					"error":   constants.ConflictError,
					"message": "You can not accept an invite to an organization you are already a member of.",
				})
			}

			if err := r.Services.Organizations.AcceptInvite(currentUser.Id, existingOrganization.Id, currentUser.Id); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			return c.SendStatus(fiber.StatusOK)
		},
	}
}
