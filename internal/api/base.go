package api

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/connor-davis/threereco-nextgen/internal/models"
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

type BaseApi[Entity any] interface {
	CreateRoute(middleware ...fiber.Handler) routing.Route
	UpdateRoute(middleware ...fiber.Handler) routing.Route
	DeleteRoute(middleware ...fiber.Handler) routing.Route
	GetOneRoute(middleware ...fiber.Handler) routing.Route
	GetAllRoute(middleware ...fiber.Handler) routing.Route
}

type baseApi[Entity any] struct {
	storage   storage.Storage
	baseUrl   string
	name      string
	createRef string
	updateRef string
}

type UpdateParams struct {
	Id uuid.UUID `param:"id"`
}

type DeleteParams struct {
	Id uuid.UUID `param:"id"`
}

type GetOneParams struct {
	Id uuid.UUID `param:"id"`
}

type GetOneQueryParams struct {
	Preloads pq.StringArray `query:"preload"`
}

type GetAllQueryParams struct {
	Page          int            `query:"page"`
	PageSize      int            `query:"pageSize"`
	Preloads      pq.StringArray `query:"preload"`
	SearchTerm    string         `query:"searchTerm"`
	SearchColumns pq.StringArray `query:"searchColumn"`
}

func NewBaseApi[Entity any](storage storage.Storage, baseUrl string, name string, createRef string, updateRef string) BaseApi[Entity] {
	return &baseApi[Entity]{
		storage:   storage,
		baseUrl:   baseUrl,
		name:      name,
		createRef: createRef,
		updateRef: updateRef,
	}
}

func (c *baseApi[Entity]) CreateRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s created successfully.", c.name)).
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
			Summary:     fmt.Sprintf("Create %s", c.name),
			Description: fmt.Sprintf("This endpoint creates a new %s.", strings.ToLower(c.name)),
			Tags:        []string{fmt.Sprintf("%s", inflect.Pluralize(c.name))},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: openapi3.NewRequestBody().
					WithRequired(true).
					WithJSONSchemaRef(&openapi3.SchemaRef{
						Ref: c.createRef,
					}).
					WithDescription(fmt.Sprintf("Payload to create a new %s.", strings.ToLower(c.name))),
			},
			Responses: responses,
		},
		Entity:      c.name,
		CreateRef:   &c.createRef,
		UpdateRef:   nil,
		Method:      routing.POST,
		Path:        fmt.Sprintf("%s", c.baseUrl),
		Middlewares: middleware,
		Handler: func(ctx *fiber.Ctx) error {
			var entity Entity

			if err := ctx.BodyParser(&entity); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			if err := c.storage.Database().Save(&entity).Error; err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			entityValue := reflect.ValueOf(entity)
			entityId := entityValue.FieldByName("Id").Interface()

			return ctx.Status(fiber.StatusOK).SendString(entityId.(uuid.UUID).String())
		},
	}
}

func (c *baseApi[Entity]) UpdateRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s updated successfully", c.name)).
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
			Summary:     fmt.Sprintf("Update %s", c.name),
			Description: fmt.Sprintf("This endpoint updates an existing %s.", strings.ToLower(c.name)),
			Tags:        []string{fmt.Sprintf("%s", inflect.Pluralize(c.name))},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewPathParameter("id").
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
			},
			RequestBody: &openapi3.RequestBodyRef{
				Value: openapi3.NewRequestBody().
					WithRequired(true).
					WithJSONSchemaRef(&openapi3.SchemaRef{
						Ref: c.updateRef,
					}).
					WithDescription(fmt.Sprintf("Payload to update an existing %s.", strings.ToLower(c.name))),
			},
			Responses: responses,
		},
		Entity:      c.name,
		CreateRef:   nil,
		UpdateRef:   &c.updateRef,
		Method:      routing.PUT,
		Path:        fmt.Sprintf("%s/{id}", c.baseUrl),
		Middlewares: middleware,
		Handler: func(ctx *fiber.Ctx) error {
			var params UpdateParams

			if err := ctx.ParamsParser(&params); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			var entity Entity

			if err := ctx.BodyParser(&entity); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			var existingEntity Entity

			if err := c.storage.Database().Where("id = ?", params.Id).First(&existingEntity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.name)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			if err := c.storage.Database().Model(&existingEntity).Updates(&entity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.name)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			if err := replaceAssociations(c.storage.Database(), &existingEntity, &entity); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).SendString("OK")
		},
	}
}

func replaceAssociations(db *gorm.DB, existingEntity any, newEntity any) error {
	v := reflect.ValueOf(newEntity)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// Only handle exported fields
		if !fieldVal.CanInterface() {
			continue
		}

		// Detect slice or struct (associations)
		if fieldVal.Kind() == reflect.Slice && (fieldVal.Type() == reflect.TypeOf([]models.Business{}) || fieldVal.Type() == reflect.TypeOf([]models.Role{})) {
			assocName := field.Name
			val := fieldVal.Interface()

			// Replace association
			if err := db.Model(existingEntity).Association(assocName).Replace(val); err != nil {
				return fmt.Errorf("failed replacing %s: %w", assocName, err)
			}
		}
	}
	return nil
}

func (c *baseApi[Entity]) DeleteRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s deleted successfully.", c.name)).
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
			Summary:     fmt.Sprintf("Delete %s", c.name),
			Description: fmt.Sprintf("This endpoint deletes an existing %s.", strings.ToLower(c.name)),
			Tags:        []string{fmt.Sprintf("%s", inflect.Pluralize(c.name))},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewPathParameter("id").
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Entity:      c.name,
		CreateRef:   nil,
		UpdateRef:   nil,
		Method:      routing.DELETE,
		Path:        fmt.Sprintf("%s/{id}", c.baseUrl),
		Middlewares: middleware,
		Handler: func(ctx *fiber.Ctx) error {
			var params DeleteParams

			if err := ctx.ParamsParser(&params); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			if err := c.storage.Database().Model(new(Entity)).Where("id = ?", params.Id).Delete(new(Entity)).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.name)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).SendString("OK")
		},
	}
}

func (c *baseApi[Entity]) GetOneRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s retrieved successfully.", c.name)).
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
			Summary:     fmt.Sprintf("Get %s", c.name),
			Description: fmt.Sprintf("This endpoint retrieves an existing %s.", strings.ToLower(c.name)),
			Tags:        []string{fmt.Sprintf("%s", inflect.Pluralize(c.name))},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewPathParameter("id").
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
				{
					Value: openapi3.NewQueryParameter("preload").
						WithRequired(false).
						WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Entity:      c.name,
		CreateRef:   nil,
		UpdateRef:   nil,
		Method:      routing.GET,
		Path:        fmt.Sprintf("%s/{id}", c.baseUrl),
		Middlewares: middleware,
		Handler: func(ctx *fiber.Ctx) error {
			var queryParams GetOneQueryParams

			if err := ctx.QueryParser(&queryParams); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			var params GetOneParams

			if err := ctx.ParamsParser(&params); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			var entity Entity

			query := c.storage.Database().Model(&entity)

			for _, preload := range queryParams.Preloads {
				parts := strings.Split(preload, ".")

				for i, part := range parts {
					parts[i] = inflect.Camelize(strings.ToLower(inflect.Dasherize(part)))
				}

				query = query.Preload(strings.Join(parts, "."))
			}

			policies, ok := ctx.Locals("policies").(clause.Expression)

			if ok && policies != nil {
				query = query.Clauses(policies)
			}

			if err := query.Where("id = ?", params.Id).First(&entity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.name)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
				"item": entity,
			})
		},
	}
}

func (c *baseApi[Entity]) GetAllRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s's retrieved successfully.", c.name)).
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
			Summary:     fmt.Sprintf("Get %s", inflect.Pluralize(c.name)),
			Description: fmt.Sprintf("This endpoint retrieves a list of %s.", strings.ToLower(c.name)),
			Tags:        []string{fmt.Sprintf("%s", inflect.Pluralize(c.name))},
			Parameters: []*openapi3.ParameterRef{
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
		Entity:      c.name,
		CreateRef:   nil,
		UpdateRef:   nil,
		Method:      routing.GET,
		Path:        fmt.Sprintf("%s", c.baseUrl),
		Middlewares: middleware,
		Handler: func(ctx *fiber.Ctx) error {
			var queryParams GetAllQueryParams

			if err := ctx.QueryParser(&queryParams); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
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

			var entities []Entity
			totalEntities := int64(0)
			countQuery := c.storage.Database().Model(new(Entity))

			if len(clauses) > 0 {
				countQuery = countQuery.Clauses(clauses...)
			}

			if err := countQuery.Count(&totalEntities).Error; err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

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

			query := c.storage.Database().Model(&entities)

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

			if err := query.Limit(limit).Offset(offset).Find(&entities).Error; err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"items": entities,
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
