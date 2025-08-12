package auditlogs

import (
	"math"
	"strconv"

	"github.com/connor-davis/threereco-nextgen/internal/constants"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (r *AuditLogsRouter) ViewRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("The list of auditlogs for the specified page and search query.").
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
			Summary:     "View AuditLogs",
			Description: "Endpoint to retrieve a list of auditlogs with pagination and optional search query",
			Tags:        []string{"AuditLogs"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GetMethod,
		Path:   "/auditlogs",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(*models.User)

			if currentUser.PrimaryOrganizationId == nil {
				log.Errorf("🔥 Current user does not belong to any organization.")

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": "You must belong to and have selected an organization to view auditlogs.",
				})
			}

			search := c.Query("search")

			page, err := strconv.Atoi(c.Query("page"))

			if err != nil {
				page = 1
			}

			searchClauses := clause.Or(
				clause.Like{
					Column: clause.Column{
						Name: "table_name",
					},
					Value: "%" + search + "%",
				},
				clause.Like{
					Column: clause.Column{
						Name: "operation_type",
					},
					Value: "%" + search + "%",
				},
			)

			totalAuditLogs, err := r.Services.AuditLogs.GetTotal(
				*currentUser.PrimaryOrganizationId,
				searchClauses,
			)

			if err != nil {
				log.Errorf("🔥 Failed to get total auditlogs: %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			pages := max(int(math.Ceil(float64(totalAuditLogs)/10)), 1)
			nextPage := min(page+1, pages)
			previousPage := max(page-1, 1)

			pageDetails := fiber.Map{
				"count":        totalAuditLogs,
				"nextPage":     nextPage,
				"previousPage": previousPage,
				"currentPage":  page,
				"pages":        pages,
			}

			limit := 10
			offset := (page - 1) * 10

			auditlogs, err := r.Services.AuditLogs.GetAll(
				*currentUser.PrimaryOrganizationId,
				clause.Limit{
					Limit:  &limit,
					Offset: offset,
				},
				searchClauses,
			)

			if err != nil {
				log.Errorf("🔥 Failed to get auditlogs: %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"items":       auditlogs,
				"pageDetails": pageDetails,
			})
		},
	}
}

func (r *AuditLogsRouter) ViewByIdRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("The auditlog was retrieved successfully.").
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
			WithDescription("AuditLog not found.").
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
			Summary:     "View AuditLog",
			Description: "Endpoint to retrieve a auditlog by their ID",
			Tags:        []string{"AuditLogs"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GetMethod,
		Path:   "/auditlogs/{id}",
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

			auditlog, err := r.Services.AuditLogs.GetById(idUUID)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if auditlog.Id == uuid.Nil {
				return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
					"error":   constants.NotFoundError,
					"details": constants.NotFoundErrorDetails,
				})
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"item": auditlog,
			})
		},
	}
}
