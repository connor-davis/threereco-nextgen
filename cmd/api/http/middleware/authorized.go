package middleware

import (
	"strings"

	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (m *middleware) Authorized(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)

		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "You must be logged in to access this resource.",
			})
		}

		combinedPermissions := []string{}

		for _, permission := range user.Permissions {
			combinedPermissions = append(combinedPermissions, permission)
		}

		for _, role := range user.Roles {
			for _, permission := range role.Permissions {
				combinedPermissions = append(combinedPermissions, permission)
			}
		}

		if len(combinedPermissions) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Forbidden",
				"message": "You do not have permission to access this resource.",
			})
		}

		for _, userPermission := range combinedPermissions {
			if userPermission == "*" {
				return c.Next()
			}

			userPermission = strings.TrimSuffix(strings.TrimSpace(userPermission), ".*")

			for _, requiredPermission := range permissions {
				requiredPermission = strings.TrimSpace(requiredPermission)

				if userPermission == requiredPermission {
					return c.Next()
				}

				if strings.HasPrefix(requiredPermission, userPermission) {
					return c.Next()
				}
			}
		}

		log.Warnf("⚠️ User %s does not have required permissions: %v", user.Username, permissions)

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Forbidden",
			"message": "You do not have permission to access this resource.",
		})
	}
}
