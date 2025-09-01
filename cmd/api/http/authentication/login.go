package authentication

import (
	"time"

	"github.com/connor-davis/threereco-nextgen/internal/constants"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/bodies"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginPayload struct {
	EmailOrPhone string `json:"emailOrPhone"`
	Password     string `json:"password"`
}

func (r *AuthenticationRouter) LoginRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The user has been successfully logged in."),
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

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().WithJSONSchema(
			schemas.ErrorResponseSchema.Value,
		).WithDescription("Internal Server Error.").WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				},
				Schema: schemas.ErrorResponseSchema,
			},
		}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     "Login",
			Description: "Logs in a user with email and password.",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: bodies.LoginPayloadBody,
			Responses:   responses,
		},
		Method:      routing.PostMethod,
		Path:        "/authentication/login",
		Middlewares: []fiber.Handler{},
		Handler: func(c *fiber.Ctx) error {
			var payload LoginPayload

			if err := c.BodyParser(&payload); err != nil {
				log.Errorf("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			user, err := r.Services.Users.GetByEmailOrPhone(payload.EmailOrPhone)

			if err != nil {
				log.Errorf("🔥 Error retrieving user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			if user.Id == uuid.Nil {
				log.Warnf("⚠️ User with email or phone %s not found", payload.EmailOrPhone)

				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   constants.UnauthorizedError,
					"message": constants.UnauthorizedErrorDetails,
				})
			}

			err = bcrypt.CompareHashAndPassword(user.Password, []byte(payload.Password))

			if err != nil {
				log.Warnf("⚠️ Invalid password for user %s: %s", payload.EmailOrPhone, err.Error())

				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   constants.UnauthorizedError,
					"message": constants.UnauthorizedErrorDetails,
				})
			}

			currentSession, err := r.Sessions.Get(c)

			if err != nil {
				log.Errorf("🔥 Error retrieving session: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			currentSession.Set("user_id", user.Id.String())
			currentSession.SetExpiry(1 * time.Hour)

			if err := currentSession.Save(); err != nil {
				log.Errorf("🔥 Error saving session: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			if err := r.Storage.Postgres.
				Set("one:ignore_audit_log", true).
				Model(&user).
				Updates(map[string]any{
					"mfa_verified": false,
				}).Error; err != nil {
				log.Errorf("🔥 Error updating MFA status for user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			return c.SendStatus(fiber.StatusOK)
		},
	}
}
