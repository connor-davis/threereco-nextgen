package organizations

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/organizations/users"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/services"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type OrganizationsRouter struct {
	Storage    *storage.Storage
	Sessions   *session.Store
	Services   *services.Services
	Middleware *middleware.Middleware
	users      *users.OrganizationUsersRouter
}

func NewOrganizationsRouter(storage *storage.Storage, sessions *session.Store, services *services.Services, middleware *middleware.Middleware) *OrganizationsRouter {
	usersRouter := users.NewOrganizationUsersRouter(storage, sessions, services, middleware)

	return &OrganizationsRouter{
		Storage:    storage,
		Sessions:   sessions,
		Services:   services,
		Middleware: middleware,
		users:      usersRouter,
	}
}

func (r *OrganizationsRouter) InitializeRoutes() []routing.Route {
	viewRoute := r.ViewRoute()
	viewByIdRoute := r.ViewByIdRoute()

	createRoute := r.CreateRoute()
	updateByIdRoute := r.UpdateByIdRoute()
	removeByIdRoute := r.users.RemoveByIdRoute()
	deleteByIdRoute := r.DeleteByIdRoute()

	return []routing.Route{
		viewRoute,
		viewByIdRoute,
		createRoute,
		updateByIdRoute,
		removeByIdRoute,
		deleteByIdRoute,
	}
}
