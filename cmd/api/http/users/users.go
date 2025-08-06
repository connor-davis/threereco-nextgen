package users

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/services"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type UsersRouter struct {
	Storage    *storage.Storage
	Sessions   *session.Store
	Services   *services.Services
	Middleware *middleware.Middleware
}

func NewUsersRouter(storage *storage.Storage, sessions *session.Store, services *services.Services, middleware *middleware.Middleware) *UsersRouter {
	return &UsersRouter{
		Storage:    storage,
		Sessions:   sessions,
		Services:   services,
		Middleware: middleware,
	}
}

func (r *UsersRouter) InitializeRoutes() []routing.Route {
	createRoute := r.CreateRoute()

	return []routing.Route{
		createRoute,
	}
}
