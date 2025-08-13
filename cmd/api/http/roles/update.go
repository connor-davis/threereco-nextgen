package roles

import (
	"github.com/connor-davis/threereco-nextgen/internal/constants"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/bodies"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type UpdateRolePayload struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Permissions []string `json:"permissions"`
}

func (r *RolesRouter) UpdateByIdRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The role has been successfully updated."),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Bad Request.").
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
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Unauthorized.").
			WithContent(openapi3.Content{
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
		).WithDescription("Role not found.").WithContent(openapi3.Content{
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
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Internal Server Error.").
			WithContent(openapi3.Content{
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
			Summary:     "Update Role by ID",
			Description: "Updates the role information for a specific role identified by their id.",
			Tags:        []string{"Roles"},
			Parameters:  parameters,
			RequestBody: bodies.UpdateRolePayloadBody,
			Responses:   responses,
		},
		Method: routing.PutMethod,
		Path:   "/roles/{id}",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(*models.User)

			id := c.Params("id")

			if id == "" {
				log.Infof("🔥 Missing role ID in request path")

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			idUUID, err := uuid.Parse(id)

			if err != nil {
				log.Infof("🔥 Invalid role ID format: %v", err)

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			var payload models.UpdateRolePayload

			if err := c.BodyParser(&payload); err != nil {
				log.Infof("🔥 Failed to parse request body: %v", err)

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			if err := r.Services.Roles.Update(currentUser.Id, idUUID, payload); err != nil {
				log.Errorf("🔥 Error updating role: %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			return c.SendStatus(fiber.StatusOK)
		},
	}
}
