package users

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/services"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type OrganizationUsersRouter struct {
	Storage    *storage.Storage
	Sessions   *session.Store
	Services   *services.Services
	Middleware *middleware.Middleware
}

func NewOrganizationUsersRouter(storage *storage.Storage, sessions *session.Store, services *services.Services, middleware *middleware.Middleware) *OrganizationUsersRouter {
	return &OrganizationUsersRouter{
		Storage:    storage,
		Sessions:   sessions,
		Services:   services,
		Middleware: middleware,
	}
}

func (r *OrganizationUsersRouter) InitializeRoutes() []routing.Route {
	removeByIdRoute := r.RemoveByIdRoute()

	return []routing.Route{
		removeByIdRoute,
	}
}
