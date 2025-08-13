package users

import (
	"github.com/connor-davis/threereco-nextgen/internal/constants"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func (r *OrganizationUsersRouter) RemoveByIdRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The user has been successfully removed."),
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

	responses.Set("404", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().WithJSONSchema(
			schemas.ErrorResponseSchema.Value,
		).WithDescription("Organization not found.").WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"error":   constants.NotFoundError,
					"message": constants.NotFoundErrorDetails,
				},
				Schema: schemas.ErrorResponseSchema,
			},
		}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().WithJSONSchema(
			schemas.ErrorResponseSchema.Value,
		).WithDescription("Internal server error.").WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				},
				Schema: schemas.ErrorResponseSchema,
			},
		}),
	})

	parameters := []*openapi3.ParameterRef{
		{
			Value: &openapi3.Parameter{
				Name:     "id",
				In:       "path",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{
							"uuid",
						},
					},
				},
			},
		},
	}

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     "Remove User by ID",
			Description: "Removes a user by their id.",
			Tags:        []string{"Organizations"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.DeleteMethod,
		Path:   "/organizations/users/{id}",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(*models.User)

			if currentUser.PrimaryOrganizationId == nil {
				log.Errorf("🔥 Current user does not belong to any organization.")

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": "You must belong to and have selected an organization to create transactions.",
				})
			}

			id := c.Params("id")

			if id == "" {
				log.Infof("🔥 Organization id is required.")

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			idUUID, err := uuid.Parse(id)

			if err != nil {
				log.Errorf("🔥 Error parsing organization id: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			existingUser, err := r.Services.Users.GetById(idUUID)

			if err != nil {
				log.Errorf("🔥 Error retrieving user: %s", err.Error())

				return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			if existingUser.Id == uuid.Nil {
				log.Infof("🔥 User with ID %s not found", id)

				return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
					"error":   constants.NotFoundError,
					"message": constants.NotFoundErrorDetails,
				})
			}

			if err := r.Services.Organizations.RemoveUser(currentUser.Id, *currentUser.PrimaryOrganizationId, idUUID); err != nil {
				log.Errorf("🔥 Error removing user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			return c.SendStatus(fiber.StatusOK)
		},
	}
}
