package authentication

import (
	"fmt"
	"time"

	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (r *AuthenticationRouter) RegisterRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Registered successfully.").
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
			Summary:     "Register",
			Description: "Registers a new user with email and password.",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Ref: "#/components/requestBodies/RegisterPayload",
			},
			Responses: responses,
		},
		Method:      routing.POST,
		Path:        "/authentication/register",
		Middlewares: []fiber.Handler{},
		Handler: func(c *fiber.Ctx) error {
			user, ok := c.Locals("user").(*models.User)

			if ok && user != nil {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   "Bad Request",
					"message": "You are already logged in.",
				})
			}

			var payload models.RegisterPayload

			if err := c.BodyParser(&payload); err != nil {
				log.Errorf("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   "Bad Request",
					"message": "The request body is invalid.",
				})
			}

			if payload.Type == models.SystemUser && user.Type != models.SystemUser {
				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   "Unauthorized",
					"message": "You are not authorized to create a system user.",
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
					hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

					if err != nil {
						log.Errorf("🔥 Error hashing password: %s", err.Error())

						return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
							"error":   "Internal Server Error",
							"message": "An error occurred while processing your request.",
						})
					}

					newUser := models.User{
						Name:        payload.Name,
						Username:    payload.Username,
						Password:    hashedPassword,
						MfaVerified: false,
						MfaEnabled:  false,
						Type:        payload.Type,
					}

					if err := r.storage.Database().Create(&newUser).Error; err != nil {
						log.Errorf("🔥 Error creating user: %s", err.Error())

						return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
							"error":   "Internal Server Error",
							"message": "An error occurred while processing your request.",
						})
					}

					if payload.Type == models.BusinessUser {
						newBusinessName := fmt.Sprintf("%s Business", inflect.Pluralize(newUser.Name))

						businessOwnerRoleName := "Business Owner"
						businessOwnerRoleDescription := "Owner of the business with full access to business resources."

						businessStaffRoleName := "Business Staff"
						businessStaffRoleDescription := "Staff member of the business with limited access to business resources."

						businessUserRoleName := "Business User"
						businessUserRoleDescription := "User of the business with minimal access to business resources."

						newBusinessOwnerRole := models.Role{
							Name:        businessOwnerRoleName,
							Description: &businessOwnerRoleDescription,
							Permissions: []string{
								"materials.view",
								"collections.*",
								"transactions.*",
								"users.view.self",
								"users.update.self",
								"users.delete.self",
								"businesses.view",
								"businesses.update.self",
								"businesses.delete.self",
								"businesses.roles.assign",
								"businesses.roles.unassign",
								"businesses.roles.view",
								"businesses.users.assign",
								"businesses.users.unassign",
								"businesses.users.view",
							},
							Default: false,
						}

						newBusinessStaffRole := models.Role{
							Name:        businessStaffRoleName,
							Description: &businessStaffRoleDescription,
							Permissions: []string{
								"materials.view",
								"collections.view",
								"collections.create",
								"collections.update",
								"transactions.view",
								"transactions.create",
								"transactions.update",
								"users.view.self",
								"users.update.self",
								"users.delete.self",
								"businesses.view",
								"businesses.users.view",
							},
							Default: false,
						}

						newBusinessUserRole := models.Role{
							Name:        businessUserRoleName,
							Description: &businessUserRoleDescription,
							Permissions: []string{
								"materials.view",
								"collections.view",
								"transactions.view",
								"users.view.self",
								"users.update.self",
								"users.delete.self",
								"businesses.view",
							},
							Default: false,
						}

						if err := r.storage.Database().Where("name = ?", businessOwnerRoleName).FirstOrCreate(&newBusinessOwnerRole).Error; err != nil {
							log.Errorf("🔥 Error creating business role: %s", err.Error())

							return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
								"error":   "Internal Server Error",
								"message": "An error occurred while processing your request.",
							})
						}

						if err := r.storage.Database().Where("name = ?", businessStaffRoleName).FirstOrCreate(&newBusinessStaffRole).Error; err != nil {
							log.Errorf("🔥 Error creating business role: %s", err.Error())

							return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
								"error":   "Internal Server Error",
								"message": "An error occurred while processing your request.",
							})
						}

						if err := r.storage.Database().Where("name = ?", businessUserRoleName).FirstOrCreate(&newBusinessUserRole).Error; err != nil {
							log.Errorf("🔥 Error creating business role: %s", err.Error())

							return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
								"error":   "Internal Server Error",
								"message": "An error occurred while processing your request.",
							})
						}

						newBusiness := models.Business{
							Name: newBusinessName,
							Roles: []models.Role{
								newBusinessOwnerRole,
								newBusinessStaffRole,
								newBusinessUserRole,
							},
							Users: []models.User{
								newUser,
							},
							OwnerId: newUser.Id,
						}

						if err := r.storage.Database().Create(&newBusiness).Error; err != nil {
							log.Errorf("🔥 Error creating business: %s", err.Error())

							return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
								"error":   "Internal Server Error",
								"message": "An error occurred while processing your request.",
							})
						}

						newUser.BusinessId = &newBusiness.Id
						newUser.Roles = []models.Role{
							newBusinessOwnerRole,
						}

						if err := r.storage.Database().
							Save(&newUser).Error; err != nil {
							log.Errorf("🔥 Error updating user with business ID: %s", err.Error())

							return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
								"error":   "Internal Server Error",
								"message": "An error occurred while processing your request.",
							})
						}
					}

					currentSession, err := r.session.Get(c)

					if err != nil {
						log.Errorf("🔥 Error retrieving session: %s", err.Error())

						return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
							"error":   "Internal Server Error",
							"message": "An error occurred while processing your request.",
						})
					}

					currentSession.Set("user_id", newUser.Id.String())
					currentSession.SetExpiry(1 * time.Hour)

					if err := currentSession.Save(); err != nil {
						log.Errorf("🔥 Error saving session: %s", err.Error())

						return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
							"error":   "Internal Server Error",
							"message": "An error occurred while processing your request.",
						})
					}

					return c.SendStatus(fiber.StatusOK)
				}

				log.Errorf("🔥 Error retrieving user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   "Internal Server Error",
					"message": "An error occurred while processing your request.",
				})
			}

			return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
				"error":   "Conflict",
				"message": "An error occurred while processing your request.",
			})
		},
	}
}
