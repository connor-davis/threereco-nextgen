package mfa

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pquerna/otp/totp"
)

func (r *MfaRouter) VerifyRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Successful authentication check.").
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
			Summary:     "Verify MFA",
			Description: "Verifies the Multi-Factor Authentication code for the user.",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Ref: "#/components/requestBodies/VerifyMfaPayload",
			},
			Responses: responses,
		},
		Method: routing.POST,
		Path:   "/authentication/mfa/verify",
		Middlewares: []fiber.Handler{
			r.middleware.Authenticated(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(*models.User)

			var payload models.VerifyMfaPayload

			if err := c.BodyParser(&payload); err != nil {
				log.Infof("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": "The request body is invalid.",
				})
			}

			if payload.Code == "" || len(payload.Code) < 6 || len(payload.Code) > 6 {
				log.Warn("🚫 Unauthorized access attempt: No MFA code provided")

				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": "Unable to verify Multi-Factor Authentication (MFA) status. Please provide a valid MFA code.",
				})
			}

			if currentUser == nil || currentUser.MfaSecret == nil {
				log.Warn("🚫 Unauthorized access attempt: User not found or MFA not enabled")

				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": "Unable to verify Multi-Factor Authentication (MFA) status. Please ensure MFA is enabled for your account.",
				})
			}

			if !totp.Validate(payload.Code, string(currentUser.MfaSecret)) {
				return c.Status(fiber.StatusUnauthorized).
					JSON(fiber.Map{
						"error":   "Unauthorized",
						"message": "Invalid Multi-Factor Authentication code. Please try again.",
					})
			}

			currentUser.MfaEnabled = true
			currentUser.MfaVerified = true

			if err := r.storage.Database().Set("one:ignore_audit_log", true).
				Where("id = ?", currentUser.Id).
				Updates(&currentUser).Error; err != nil {
				log.Errorf("🔥 Error updating user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": "An error occurred while processing your request.",
				})
			}

			return c.SendStatus(fiber.StatusOK)
		},
	}
}
