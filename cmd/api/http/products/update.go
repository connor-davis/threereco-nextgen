package products

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

func (r *ProductsRouter) UpdateByIdRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The product has been successfully updated."),
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
		).WithDescription("Product not found.").WithContent(openapi3.Content{
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
			Summary:     "Update Product by ID",
			Description: "Updates the product information for a specific product identified by their id.",
			Tags:        []string{"Products"},
			Parameters:  parameters,
			RequestBody: bodies.UpdateProductPayloadBody,
			Responses:   responses,
		},
		Method: routing.PutMethod,
		Path:   "/products/{id}",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(*models.User)

			id := c.Params("id")

			if id == "" {
				log.Infof("🔥 Missing product ID in request path")

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			idUUID, err := uuid.Parse(id)

			if err != nil {
				log.Infof("🔥 Invalid product ID format: %v", err)

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			var payload models.UpdateProductPayload

			if err := c.BodyParser(&payload); err != nil {
				log.Infof("🔥 Failed to parse request body: %v", err)

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			if err := r.Services.Products.Update(currentUser.Id, idUUID, payload); err != nil {
				log.Errorf("🔥 Error updating product: %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			return c.SendStatus(fiber.StatusOK)
		},
	}
}
