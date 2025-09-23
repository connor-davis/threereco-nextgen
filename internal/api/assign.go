package api

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AssignmentApi[Parent any, Child any] interface {
	AssignRoute(middleware ...fiber.Handler) routing.Route
	UnassignRoute(middleware ...fiber.Handler) routing.Route
	ListRoute(middleware ...fiber.Handler) routing.Route
}

type assignmentApi[Parent any, Child any] struct {
	storage    storage.Storage
	baseUrl    string
	parentName string
	childName  string
}

type AssignParams struct {
	ParentId uuid.UUID `param:"parentId"`
	ChildId  uuid.UUID `param:"childId"`
}

type UnassignParams struct {
	ParentId uuid.UUID `param:"parentId"`
	ChildId  uuid.UUID `param:"childId"`
}

type ListParams struct {
	ParentId uuid.UUID `param:"parentId"`
	ChildId  uuid.UUID `param:"childId"`
}

type ListQueryParams struct {
	Page          int            `query:"page"`
	PageSize      int            `query:"pageSize"`
	Preloads      pq.StringArray `query:"preload"`
	SearchTerm    string         `query:"searchTerm"`
	SearchColumns pq.StringArray `query:"searchColumn"`
}

func NewAssignmentApi[Parent any, Child any](storage storage.Storage, baseUrl string, parentName string, childName string) AssignmentApi[Parent, Child] {
	return &assignmentApi[Parent, Child]{
		storage:    storage,
		baseUrl:    baseUrl,
		parentName: parentName,
		childName:  childName,
	}
}

func (c *assignmentApi[Parent, Child]) AssignRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s assigned to %s successfully.", c.childName, c.parentName)).
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

	responses.Set("404", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Not Found").
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
			Summary:     fmt.Sprintf("Assign %s", c.childName),
			Description: fmt.Sprintf("This endpoint assigns a %s to a %s", strings.ToLower(c.childName), strings.ToLower(c.parentName)),
			Tags:        []string{fmt.Sprintf("%s Assignments", inflect.Pluralize(c.parentName))},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewPathParameter(fmt.Sprintf("%sId", inflect.Parameterize(c.parentName))).
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
				{
					Value: openapi3.NewPathParameter(fmt.Sprintf("%sId", inflect.Parameterize(c.childName))).
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Entity:    c.childName,
		CreateRef: nil,
		UpdateRef: nil,
		Method:    routing.POST,
		Path: fmt.Sprintf(
			"%s/assign-%s/{%sId}/{%sId}",
			c.baseUrl,
			strings.ToLower(inflect.Dasherize(c.childName)),
			inflect.Parameterize(c.parentName),
			inflect.Parameterize(c.childName),
		),
		Middlewares: middleware,
		Handler: func(ctx *fiber.Ctx) error {
			parentId := ctx.Params(fmt.Sprintf("%sId", inflect.Parameterize(c.parentName)))
			childId := ctx.Params(fmt.Sprintf("%sId", inflect.Parameterize(c.childName)))

			var parentEntity Parent
			var childEntity Child

			if err := c.storage.Database().Model(&parentEntity).Where("id = ?", parentId).First(&parentEntity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.parentName)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			if err := c.storage.Database().Model(&childEntity).Where("id = ?", childId).First(&childEntity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.childName)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			var existingAssociation Child

			if err := c.storage.Database().Model(&parentEntity).Association(fmt.Sprintf("%ss", c.childName)).Find(&existingAssociation); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			if !reflect.ValueOf(existingAssociation).IsZero() {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": fmt.Sprintf("The %s is already assigned to the %s.", strings.ToLower(c.childName), strings.ToLower(c.parentName)),
				})
			}

			if err := c.storage.Database().Model(&parentEntity).Association(fmt.Sprintf("%ss", c.childName)).Append(&childEntity); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).SendString("OK")
		},
	}
}

func (c *assignmentApi[Parent, Child]) UnassignRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s unassigned from %s successfully.", c.childName, c.parentName)).
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

	responses.Set("404", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Not Found").
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
			Summary:     fmt.Sprintf("Unassign %s", c.childName),
			Description: fmt.Sprintf("This endpoint unassigns a %s from a %s", strings.ToLower(c.childName), strings.ToLower(c.parentName)),
			Tags:        []string{fmt.Sprintf("%s Assignments", inflect.Pluralize(c.parentName))},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewPathParameter(fmt.Sprintf("%sId", inflect.Parameterize(c.parentName))).
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
				{
					Value: openapi3.NewPathParameter(fmt.Sprintf("%sId", inflect.Parameterize(c.childName))).
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Entity:    c.childName,
		CreateRef: nil,
		UpdateRef: nil,
		Method:    routing.POST,
		Path: fmt.Sprintf(
			"%s/unassign-%s/{%sId}/{%sId}",
			c.baseUrl,
			strings.ToLower(inflect.Dasherize(c.childName)),
			inflect.Parameterize(c.parentName),
			inflect.Parameterize(c.childName),
		),
		Middlewares: middleware,
		Handler: func(ctx *fiber.Ctx) error {
			parentId := ctx.Params(fmt.Sprintf("%sId", inflect.Parameterize(c.parentName)))
			childId := ctx.Params(fmt.Sprintf("%sId", inflect.Parameterize(c.childName)))

			var parentEntity Parent
			var childEntity Child

			if err := c.storage.Database().Model(&parentEntity).Where("id = ?", parentId).First(&parentEntity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.parentName)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			if err := c.storage.Database().Model(&childEntity).Where("id = ?", childId).First(&childEntity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.childName)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			var existingAssociation Child

			if err := c.storage.Database().Model(&parentEntity).Association(fmt.Sprintf("%ss", c.childName)).Find(&existingAssociation); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			if reflect.ValueOf(existingAssociation).IsZero() {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": fmt.Sprintf("The %s is not assigned to the %s.", strings.ToLower(c.childName), strings.ToLower(c.parentName)),
				})
			}

			if err := c.storage.Database().Model(&parentEntity).Association(fmt.Sprintf("%ss", c.childName)).Delete(&childEntity); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).SendString("OK")
		},
	}
}

func (c *assignmentApi[Parent, Child]) ListRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s retrieved from %s successfully.", c.childName, c.parentName)).
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
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

	responses.Set("404", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(schemas.ErrorSchema).
			WithDescription("Not Found").
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
			Summary:     fmt.Sprintf("List %s", c.childName),
			Description: fmt.Sprintf("This endpoint retrieves a list of %s assigned to a %s", strings.ToLower(c.childName), strings.ToLower(c.parentName)),
			Tags:        []string{fmt.Sprintf("%s Assignments", inflect.Pluralize(c.parentName))},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewPathParameter(fmt.Sprintf("%sId", inflect.Parameterize(c.parentName))).
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
				{
					Value: openapi3.NewQueryParameter("page").
						WithRequired(true).
						WithSchema(openapi3.NewInt64Schema().WithDefault(1)),
				},
				{
					Value: openapi3.NewQueryParameter("pageSize").
						WithRequired(true).
						WithSchema(openapi3.NewInt64Schema().WithDefault(10)),
				},
				{
					Value: openapi3.NewQueryParameter("preload").
						WithRequired(false).
						WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
				},
				{
					Value: openapi3.NewQueryParameter("searchTerm").
						WithRequired(false).
						WithSchema(openapi3.NewStringSchema()),
				},
				{
					Value: openapi3.NewQueryParameter("searchColumn").
						WithRequired(false).
						WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Entity:    c.childName,
		CreateRef: nil,
		UpdateRef: nil,
		Method:    routing.GET,
		Path: fmt.Sprintf(
			"%s/list-%s/{%sId}",
			c.baseUrl,
			strings.ToLower(inflect.Dasherize(inflect.Pluralize(c.childName))),
			inflect.Parameterize(c.parentName),
		),
		Middlewares: middleware,
		Handler: func(ctx *fiber.Ctx) error {
			parentId := ctx.Params(fmt.Sprintf("%sId", inflect.Parameterize(c.parentName)))

			var queryParams ListQueryParams

			if err := ctx.QueryParser(&queryParams); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			var parentEntity Parent

			if err := c.storage.Database().Model(&parentEntity).Where("id = ?", parentId).First(&parentEntity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.parentName)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			clauses := []clause.Expression{}

			for _, column := range queryParams.SearchColumns {
				clauses = append(clauses, clause.Like{
					Column: clause.Column{Name: column},
					Value:  fmt.Sprintf("%%%s%%", queryParams.SearchTerm),
				})
			}

			policies, ok := ctx.Locals("policies").(clause.Expression)

			if ok && policies != nil {
				clauses = append(clauses, policies)
			}

			countQuery := c.storage.Database().Model(&parentEntity)

			if len(clauses) > 0 {
				countQuery = countQuery.Clauses(clause.Or(clauses...))
			}

			totalEntities := countQuery.Association(fmt.Sprintf("%ss", c.childName)).Count()

			if queryParams.Page < 1 {
				queryParams.Page = 1
			}

			limit := queryParams.PageSize
			offset := (queryParams.Page - 1) * queryParams.PageSize
			totalPages := int64(math.Ceil(float64(totalEntities) / float64(limit)))

			if totalPages == 0 {
				totalPages = 1
			}

			nextPage := queryParams.Page + 1
			previousPage := queryParams.Page - 1

			if nextPage > int(totalPages) {
				nextPage = int(totalPages)
			}

			if previousPage < 1 {
				previousPage = 1
			}

			var existingAssociations []Child

			query := c.storage.Database().Model(&parentEntity)

			for _, preload := range queryParams.Preloads {
				parts := strings.Split(preload, ".")

				for i, part := range parts {
					parts[i] = inflect.Camelize(strings.ToLower(part))
				}

				query = query.Preload(strings.Join(parts, "."))
			}

			if len(clauses) > 0 {
				query = query.Clauses(clause.Or(clauses...))
			}

			if err := query.Limit(limit).Offset(offset).Association(fmt.Sprintf("%ss", c.childName)).Find(&existingAssociations); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"items": existingAssociations,
				"pagination": fiber.Map{
					"count":        totalEntities,
					"pages":        totalPages,
					"pageSize":     queryParams.PageSize,
					"currentPage":  queryParams.Page,
					"nextPage":     nextPage,
					"previousPage": previousPage,
				},
			})
		},
	}
}
