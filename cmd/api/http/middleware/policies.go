package middleware

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm/clause"
)

func (m *middleware) Policies(policies ...models.PolicyType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)

		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "You must be logged in to access this resource.",
			})
		}

		for _, policy := range policies {
			switch policy {
			case models.CollectionsPolicy:
				switch user.Type {
				case models.CollectorUser:
					c.Locals("policies", clause.Or(
						clause.Eq{
							Column: clause.Column{
								Name: "seller_id",
							},
							Value: user.Id,
						},
					))
				case models.BusinessUser:
					c.Locals("policies", clause.Or(
						clause.Eq{
							Column: clause.Column{
								Name: "buyer_id",
							},
							Value: user.BusinessId,
						},
					))
				}

				return c.Next()
			case models.TransactionsPolicy:
				if user.Type == models.BusinessUser {
					c.Locals("policies", clause.Or(
						clause.Eq{
							Column: clause.Column{
								Name: "seller_id",
							},
							Value: user.BusinessId,
						},
						clause.Eq{
							Column: clause.Column{
								Name: "buyer_id",
							},
							Value: user.BusinessId,
						},
					))

					return c.Next()
				}
			case models.SystemPolicy:
				if user.Type == models.SystemUser {
					return c.Next()
				}
			}
		}

		log.Warnf("⚠️ User %s does not satisfy required policies: %v", user.Username, policies)

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Forbidden",
			"message": "You do not have permission to access this resource.",
		})
	}
}
