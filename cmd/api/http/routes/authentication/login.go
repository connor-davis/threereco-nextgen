package authentication

import (
	"time"

	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (r *AuthenticationRouter) LoginRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Logged in successfully.").
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
			Summary:     "Login",
			Description: "Logs in a user with email and password.",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Ref: "#/components/requestBodies/LoginPayload",
			},
			Responses: responses,
		},
		Method:      routing.POST,
		Path:        "/authentication/login",
		Middlewares: []fiber.Handler{},
		Handler: func(c *fiber.Ctx) error {
			var payload models.LoginPayload

			if err := c.BodyParser(&payload); err != nil {
				log.Errorf("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   "Bad Request",
					"message": "The request body is invalid.",
				})
			}

			var existingUser models.User

			if err := r.storage.Database().
				Where(
					"username = ?",
					payload.Username,
				).
				First(&existingUser).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					log.Warnf("⚠️ User with username %s not found", payload.Username)

					return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
						"error":   "Unauthorized",
						"message": "You are not authorized to access this resource.",
					})
				}

				log.Errorf("🔥 Error retrieving user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   "Internal Server Error",
					"message": "An error occurred while processing your request.",
				})
			}

			if err := bcrypt.CompareHashAndPassword(existingUser.Password, []byte(payload.Password)); err != nil {
				log.Warnf("⚠️ Invalid password for user %s: %s", payload.Username, err.Error())

				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   "Unauthorized",
					"message": "You are not authorized to access this resource.",
				})
			}

			currentSession, err := r.session.Get(c)

			if err != nil {
				log.Errorf("🔥 Error retrieving session: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   "Internal Server Error",
					"message": "An error occurred while processing your request.",
				})
			}

			currentSession.Set("user_id", existingUser.Id.String())
			currentSession.SetExpiry(1 * time.Hour)

			if err := currentSession.Save(); err != nil {
				log.Errorf("🔥 Error saving session: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   "Internal Server Error",
					"message": "An error occurred while processing your request.",
				})
			}

			if err := r.storage.Database().
				Model(&existingUser).
				Save(map[string]any{
					"mfa_verified": false,
				}).Error; err != nil {
				log.Errorf("🔥 Error updating MFA status for user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   "Internal Server Error",
					"message": "An error occurred while processing your request.",
				})
			}

			return c.SendStatus(fiber.StatusOK)
		},
	}
}
