package authentication

import (
	"fmt"

	"github.com/connor-davis/threereco-nextgen/internal/constants"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/bodies"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignUpPayload struct {
	Name     string  `json:"name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Password string  `json:"password"`
}

func (r *AuthenticationRouter) SignUpRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The user has been successfully signed up."),
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
			Summary:     "Sign Up",
			Description: "Signs up a new user with email and password.",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: bodies.SignUpPayloadBody,
			Responses:   responses,
		},
		Method:      routing.PostMethod,
		Path:        "/authentication/sign-up",
		Middlewares: []fiber.Handler{},
		Handler: func(c *fiber.Ctx) error {
			var payload SignUpPayload

			if err := c.BodyParser(&payload); err != nil {
				log.Errorf("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			if payload.Email == nil && payload.Phone == nil {
				log.Warnf("⚠️ No email or phone provided")

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"message": constants.BadRequestErrorDetails,
				})
			}

			newUserId := uuid.New()
			newOrganizationId := uuid.New()
			newOrganizationAdminRoleId := uuid.New()
			newOrganizationUserRoleId := uuid.New()

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

			if err != nil {
				log.Errorf("🔥 Error hashing password: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			organizationAdminRoleDescription := "Has full access to the organization, including user management, settings, and data."
			organizationUserRoleDescription := "Has limited access to the organization, primarily for their own profile information or data related to them that has been created by other users with their information."

			organizationAdminRole := models.Role{
				Id:               newOrganizationAdminRoleId,
				Name:             "Administrator",
				Description:      &organizationAdminRoleDescription,
				Permissions:      []string{"*"},
				IsDefault:        true,
				ModifiedByUserId: newUserId,
			}

			organizationUserRole := models.Role{
				Id:               newOrganizationUserRoleId,
				Name:             "User",
				Description:      &organizationUserRoleDescription,
				Permissions:      []string{"users.view.self", "users.update.self", "users.delete.self"},
				IsDefault:        true,
				ModifiedByUserId: newUserId,
			}

			newUser := models.User{
				Id:                    newUserId,
				Name:                  &payload.Name,
				Password:              hashedPassword,
				ModifiedByUserId:      newUserId,
				PrimaryOrganizationId: &newOrganizationId,
				Roles:                 []models.Role{organizationAdminRole},
			}

			if payload.Email != nil {
				emailUser, err := r.Services.Users.GetByEmail(*payload.Email)

				if err != nil {
					log.Errorf("🔥 Error retrieving user: %s", err.Error())

					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
						"error":   constants.InternalServerError,
						"message": constants.InternalServerErrorDetails,
					})
				}

				if emailUser.Id != uuid.Nil {
					log.Warnf("⚠️ User with email %s already exists", *payload.Email)

					return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
						"error":   constants.ConflictError,
						"message": constants.ConflictErrorDetails,
					})
				}

				newUser.Email = *payload.Email
			}

			if payload.Phone != nil {
				phoneUser, err := r.Services.Users.GetByPhone(*payload.Phone)

				if err != nil {
					log.Errorf("🔥 Error retrieving user: %s", err.Error())

					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
						"error":   constants.InternalServerError,
						"message": constants.InternalServerErrorDetails,
					})
				}

				if phoneUser.Id != uuid.Nil {
					log.Warnf("⚠️ User with phone %s already exists", *payload.Phone)

					return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
						"error":   constants.ConflictError,
						"message": constants.ConflictErrorDetails,
					})
				}

				newUser.Phone = payload.Phone
			}

			newOrganization := models.Organization{
				Id:               newOrganizationId,
				Name:             fmt.Sprintf("%s's Organization", payload.Name),
				Domain:           fmt.Sprintf("%s.com", payload.Name),
				OwnerId:          newUserId,
				Owner:            newUser,
				ModifiedByUserId: newUserId,
				Users:            []models.User{newUser},
				Roles:            []models.Role{organizationAdminRole, organizationUserRole},
			}

			if err := r.Storage.Postgres.
				Set("one:ignore_audit_log", true).
				Assign(&newOrganization).
				FirstOrCreate(&newOrganization).Error; err != nil {
				log.Errorf("❌ Failed to create organization admin user: %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"message": constants.InternalServerErrorDetails,
				})
			}

			return c.SendStatus(fiber.StatusOK)
		},
	}
}
