package authentication

import (
	"github.com/connor-davis/threereco-nextgen/internal/constants"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/bodies"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pquerna/otp/totp"
)

type MfaVerifyPayload struct {
	Code string `json:"code"`
}

func (r *AuthenticationRouter) MfaVerifyRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("MFA verification successful."),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorResponseSchema.Value).
			WithDescription(constants.BadRequestError).
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorResponseSchema.Value).
					WithExample("example", map[string]any{
						"error":   constants.BadRequestError,
						"message": "Invalid request payload",
					}),
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorResponseSchema.Value).
			WithDescription(constants.UnauthorizedError).
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorResponseSchema.Value).
					WithExample("example", map[string]any{
						"error":   constants.UnauthorizedError,
						"message": constants.UnauthorizedErrorDetails,
					}),
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorResponseSchema.Value).
			WithDescription(constants.InternalServerError).
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorResponseSchema.Value).
					WithExample("example", map[string]any{
						"error":   constants.InternalServerError,
						"message": constants.InternalServerErrorDetails,
					}),
			}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     "Verify Multi-Factor Authentication (MFA)",
			Description: "Verifies the user's MFA status.",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: bodies.MfaVerifyPayloadBody,
			Responses:   responses,
		},
		Method: routing.PostMethod,
		Path:   "/authentication/mfa/verify",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(*models.User)

			var payload MfaVerifyPayload

			if err := c.BodyParser(&payload); err != nil {
				log.Infof("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   constants.BadRequestError,
					"message": "Invalid request payload",
				})
			}

			if payload.Code == "" || len(payload.Code) < 6 || len(payload.Code) > 6 {
				log.Warn("🚫 Unauthorized access attempt: No MFA code provided")

				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   constants.BadRequestError,
					"message": "Unable to verify Multi-Factor Authentication (MFA) status. Please provide a valid MFA code.",
				})
			}

			if currentUser == nil || currentUser.MfaSecret == nil {
				log.Warn("🚫 Unauthorized access attempt: User not found or MFA not enabled")

				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   constants.BadRequestError,
					"message": "Unable to verify Multi-Factor Authentication (MFA) status. Please ensure MFA is enabled for your account.",
				})
			}

			if !totp.Validate(payload.Code, string(currentUser.MfaSecret)) {
				return c.Status(fiber.StatusUnauthorized).
					JSON(fiber.Map{
						"error":   constants.UnauthorizedError,
						"message": "Invalid Multi-Factor Authentication code. Please try again.",
					})
			}

			currentUser.MfaEnabled = true
			currentUser.MfaVerified = true

			if err := r.Storage.Postgres.Set("one:ignore_audit_log", true).
				Where(&models.User{
					Id: currentUser.Id,
				}).
				Updates(&currentUser).Error; err != nil {
				log.Errorf("🔥 Error updating user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			return c.SendStatus(fiber.StatusOK)
		},
	}
}
