package users

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

type UpdateUserPayload struct {
	Email         *string  `json:"email"`
	Name          *string  `json:"name"`
	Phone         *string  `json:"phone"`
	JobTitle      *string  `json:"jobTitle"`
	Organizations []string `json:"organizations"`
	Roles         []string `json:"roles"`
}

func (r *UsersRouter) UpdateByIdRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The user has been successfully updated."),
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
						"details": constants.BadRequestErrorDetails,
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
						"details": constants.UnauthorizedErrorDetails,
					},
					Schema: schemas.ErrorResponseSchema,
				},
			}),
	})

	responses.Set("404", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().WithJSONSchema(
			schemas.ErrorResponseSchema.Value,
		).WithDescription("User not found.").WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"error":   constants.NotFoundError,
					"details": constants.NotFoundErrorDetails,
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
						"details": constants.InternalServerErrorDetails,
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
			Summary:     "Update User by ID",
			Description: "Updates the user information for a specific user identified by their id.",
			Tags:        []string{"Users"},
			Parameters:  parameters,
			RequestBody: bodies.UpdateUserPayloadBody,
			Responses:   responses,
		},
		Method: routing.PutMethod,
		Path:   "/users/{id}",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*models.User)

			id := c.Params("id")

			if id == "" {
				log.Infof("🔥 Missing user ID in request path")

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			idUUID, err := uuid.Parse(id)

			if err != nil {
				log.Infof("🔥 Invalid user ID format: %v", err)

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			var payload UpdateUserPayload

			if err := c.BodyParser(&payload); err != nil {
				log.Infof("🔥 Failed to parse request body: %v", err)

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			existingUser, err := r.Services.Users.GetById(idUUID)

			if err != nil {
				log.Errorf("🔥 Error retrieving user by ID: %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if existingUser.Id == uuid.Nil {
				log.Infof("🔥 User with ID %s not found", id)

				return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
					"error":   constants.NotFoundError,
					"details": constants.NotFoundErrorDetails,
				})
			}

			if payload.Email != nil {
				existingUser.Email = *payload.Email
			}

			existingUser.Name = payload.Name
			existingUser.Phone = payload.Phone
			existingUser.JobTitle = payload.JobTitle

			if len(payload.Organizations) > 0 {
				for _, organizationId := range payload.Organizations {
					organizationIdUUID, err := uuid.Parse(organizationId)

					if err != nil {
						log.Infof("🔥 Invalid organization ID format: %v", err)

						continue
					}

					existingOrganization, err := r.Services.Organizations.GetById(organizationIdUUID)

					if err != nil {
						log.Errorf("🔥 Error retrieving organization by ID: %v", err)

						continue
					}

					if existingOrganization.Id == uuid.Nil {
						log.Infof("🔥 Organization with ID %s not found", organizationId)

						continue
					}

					existingUser.Organizations = append(existingUser.Organizations, *existingOrganization)
				}
			}

			if len(payload.Roles) > 0 {
				for _, roleId := range payload.Roles {
					roleIdUUID, err := uuid.Parse(roleId)

					if err != nil {
						log.Infof("🔥 Invalid role ID format: %v", err)

						continue
					}

					existingRole, err := r.Services.Roles.GetById(roleIdUUID)

					if err != nil {
						log.Errorf("🔥 Error retrieving role by ID: %v", err)

						continue
					}

					if existingRole.Id == uuid.Nil {
						log.Infof("🔥 Role with ID %s not found", roleId)

						continue
					}

					existingUser.Roles = append(existingUser.Roles, *existingRole)
				}
			}

			if err := r.Services.Users.Update(user.Id, idUUID, existingUser); err != nil {
				log.Errorf("🔥 Error updating user: %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.SendStatus(fiber.StatusOK)
		},
	}
}
