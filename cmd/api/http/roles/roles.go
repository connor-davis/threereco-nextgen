package roles

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/services"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type RolesRouter struct {
	Storage    *storage.Storage
	Sessions   *session.Store
	Services   *services.Services
	Middleware *middleware.Middleware
}

func NewRolesRouter(storage *storage.Storage, sessions *session.Store, services *services.Services, middleware *middleware.Middleware) *RolesRouter {
	return &RolesRouter{
		Storage:    storage,
		Sessions:   sessions,
		Services:   services,
		Middleware: middleware,
	}
}

func (r *RolesRouter) InitializeRoutes() []routing.Route {
	viewRoute := r.ViewRoute()
	viewByIdRoute := r.ViewByIdRoute()

	createRoute := r.CreateRoute()
	updateByIdRoute := r.UpdateByIdRoute()
	deleteByIdRoute := r.DeleteByIdRoute()

	return []routing.Route{
		viewRoute,
		viewByIdRoute,
		createRoute,
		updateByIdRoute,
		deleteByIdRoute,
	}
}
