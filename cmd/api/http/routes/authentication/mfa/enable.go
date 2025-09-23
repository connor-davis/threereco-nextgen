package mfa

import (
	"bytes"
	"encoding/base32"
	"image/png"

	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func (r *MfaRouter) EnableRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The MFA QR code has been successfully generated and returned.").
			WithContent(openapi3.Content{
				"image/png": openapi3.NewMediaType().
					WithEncoding("png", &openapi3.Encoding{
						ContentType: "image/png",
					}),
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
			Summary:     "Enable MFA",
			Description: "Enables Multi-Factor Authentication for the user.",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GET,
		Path:   "/authentication/mfa/enable",
		Middlewares: []fiber.Handler{
			r.middleware.Authenticated(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(*models.User)

			if currentUser == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   "Unauthorized",
					"message": "You must be logged in to access this resource.",
				})
			}

			if currentUser.MfaSecret == nil {
				secret, err := totp.Generate(totp.GenerateOpts{
					Issuer:      "3REco Multi-Factor Authentication",
					AccountName: currentUser.Username,
					Period:      30,
					Digits:      otp.DigitsSix,
					Algorithm:   otp.AlgorithmSHA1,
					SecretSize:  32,
				})

				if err != nil {
					return c.SendStatus(fiber.StatusInternalServerError)
				}

				currentUser.MfaSecret = []byte(secret.Secret())

				if err := r.storage.Database().Set("one:ignore_audit_log", true).
					Where("id = ?", currentUser.Id).
					Updates(&models.User{
						MfaSecret: currentUser.MfaSecret,
					}).Error; err != nil {
					log.Infof("🔥 Failed to update user: %s", err.Error())

					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": "An error occurred while processing your request.",
					})
				}
			}

			secretBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).
				DecodeString(string(currentUser.MfaSecret))

			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			secret, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "3REco Multi-Factor Authentication",
				AccountName: currentUser.Username,
				Period:      30,
				Digits:      otp.DigitsSix,
				Algorithm:   otp.AlgorithmSHA1,
				Secret:      secretBytes,
				SecretSize:  32,
			})

			if err != nil {
				log.Infof("🔥 Failed to generate TOTP secret: %s", err.Error())

				return c.SendStatus(fiber.StatusInternalServerError)
			}

			var pngBuffer bytes.Buffer

			image, err := secret.Image(256, 256)

			if err != nil {
				log.Infof("🔥 Failed to generate QR code image: %s", err.Error())

				return c.SendStatus(fiber.StatusInternalServerError)
			}

			png.Encode(&pngBuffer, image)

			c.Response().Header.Set("Content-Type", "image/png")

			return c.Status(fiber.StatusOK).Send(pngBuffer.Bytes())
		},
	}
}
