package users

import (
	"math"
	"strconv"

	"github.com/connor-davis/threereco-nextgen/internal/constants"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (r *UsersRouter) ViewRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("The list of users for the specified page and search query.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"items": []map[string]any{},
					},
					Schema: schemas.SuccessResponseSchema,
				},
			}),
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
				Name:     "page",
				In:       "query",
				Required: false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{
							"integer",
						},
					},
				},
			},
		},
		{
			Value: &openapi3.Parameter{
				Name:     "search",
				In:       "query",
				Required: false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{
							"string",
						},
					},
				},
			},
		},
	}

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     "View Users",
			Description: "Endpoint to retrieve a list of users with pagination and optional search query",
			Tags:        []string{"Users"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GetMethod,
		Path:   "/users",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			search := c.Query("search")

			page, err := strconv.Atoi(c.Query("page"))

			if err != nil {
				page = 1
			}

			searchClauses := clause.Or(
				clause.Like{
					Column: clause.Column{
						Name: "name",
					},
					Value: "%" + search + "%",
				},
				clause.Like{
					Column: clause.Column{
						Name: "email",
					},
					Value: "%" + search + "%",
				},
				clause.Like{
					Column: clause.Column{
						Name: "phone",
					},
					Value: "%" + search + "%",
				},
			)

			totalUsers, err := r.Services.Users.GetTotal(
				searchClauses,
			)

			if err != nil {
				log.Errorf("🔥 Failed to get total users: %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			pages := max(int(math.Ceil(float64(totalUsers)/10)), 1)
			nextPage := min(page+1, pages)
			previousPage := max(page-1, 1)

			pageDetails := fiber.Map{
				"count":        totalUsers,
				"nextPage":     nextPage,
				"previousPage": previousPage,
				"currentPage":  page,
				"pages":        pages,
			}

			limit := 10
			offset := (page - 1) * 10

			users, err := r.Services.Users.GetAll(
				clause.Limit{
					Limit:  &limit,
					Offset: offset,
				},
				searchClauses,
			)

			if err != nil {
				log.Errorf("🔥 Failed to get users: %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"items":       users,
				"pageDetails": pageDetails,
			})
		},
	}
}

func (r *UsersRouter) ViewByIdRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("The user was retrieved successfully.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Success,
						"details": constants.SuccessDetails,
					},
					Schema: schemas.SuccessResponseSchema,
				},
			}),
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
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("User not found.").
			WithContent(openapi3.Content{
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
							"string",
						},
					},
				},
			},
		},
	}

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     "View User",
			Description: "Endpoint to retrieve a user by their ID",
			Tags:        []string{"Users"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GetMethod,
		Path:   "/users/{id}",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			id := c.Params("id")

			if id == "" {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			idUUID, err := uuid.Parse(id)

			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			user, err := r.Services.Users.GetById(idUUID)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if user.Id == uuid.Nil {
				return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
					"error":   constants.NotFoundError,
					"details": constants.NotFoundErrorDetails,
				})
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"item": user,
			})
		},
	}
}
