package analytics

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func (r *AnalyticsRouter) MaterialTypesRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Successful material types list with percentages towards the total waste collected.").
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
			Summary:     "Material Types",
			Description: "Returns the material types in the system.",
			Tags:        []string{"Analytics"},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GET,
		Path:   "/analytics/material-types",
		Middlewares: []fiber.Handler{
			r.middleware.Authenticated(),
		},
		Handler: func(c *fiber.Ctx) error {
			user, ok := c.Locals("user").(*models.User)

			if !ok || user == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   "Unauthorized",
					"message": "You must be logged in to access this resource.",
				})
			}

			baseQuery := r.storage.Database().Model(&models.Collection{})

			if user.Type == models.BusinessUser {
				baseQuery = baseQuery.Where("buyer_id = ?", user.BusinessId)
			}

			var totalWeight float64

			if err := baseQuery.
				Table("collections").
				Joins("JOIN collections_materials ON collections_materials.collection_id = collections.id").
				Joins("JOIN collection_materials ON collection_materials.id = collections_materials.collection_material_id").
				Select("COALESCE(SUM(collection_materials.weight), 0)"). // returns 0 instead of NULL
				Scan(&totalWeight).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": "Failed to retrieve total weight.",
				})
			}

			var materialTypes []models.MaterialType

			if err := r.storage.Database().
				Table("collections").
				Joins("JOIN collections_materials ON collections_materials.collection_id = collections.id").
				Joins("JOIN collection_materials ON collection_materials.id = collections_materials.collection_material_id").
				Select("collection_materials.name AS type, COALESCE(SUM(collection_materials.weight), 0) AS weight").
				Group("collection_materials.name").
				Order("SUM(collection_materials.weight) DESC").
				Scan(&materialTypes).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": "Failed to retrieve material types.",
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"items": materialTypes,
			})
		},
	}
}
