package middleware

import (
	"time"

	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (m *middleware) Authenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentSession, err := m.session.Get(c)

		if err != nil {
			log.Errorf("🔥 Failed to retrieve session: %s", err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		}

		currentUserId, ok := currentSession.Get("user_id").(string)

		if !ok || currentUserId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "You must be logged in to access this resource.",
			})
		}

		currentUserIdUUID, err := uuid.Parse(currentUserId)

		if err != nil {
			log.Errorf("🔥 Invalid user ID in session: %s", err.Error())

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "You must be logged in to access this resource.",
			})
		}

		var currentUser *models.User

		if err := m.storage.Database().Where("id = ?", currentUserIdUUID).Preload("Roles").Preload("Businesses").First(&currentUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   "Unauthorized",
					"message": "You must be logged in to access this resource.",
				})
			}

			log.Errorf("🔥 Failed to retrieve user from database: %s", err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		}

		c.Locals("user_id", currentUser.Id.String())
		c.Locals("user", currentUser)

		currentSession.Set("user_id", currentUser.Id.String())
		currentSession.SetExpiry(1 * time.Hour)

		if err := currentSession.Save(); err != nil {
			log.Errorf("🔥 Failed to save session: %s", err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		}

		return c.Next()
	}
}
