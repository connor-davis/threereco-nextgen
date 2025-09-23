package middleware

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Middleware interface {
	Authenticated() fiber.Handler
	Authorized(permissions ...string) fiber.Handler
	Policies(policies ...models.PolicyType) fiber.Handler
}

type middleware struct {
	storage storage.Storage
	session *session.Store
}

func New(storage storage.Storage, session *session.Store) Middleware {
	return &middleware{
		storage: storage,
		session: session,
	}
}
