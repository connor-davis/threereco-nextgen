package roles

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/api"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

type RolesRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewRolesRouter(storage storage.Storage, middleware middleware.Middleware) Router {
	return &RolesRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *RolesRouter) LoadRoutes() []routing.Route {
	api := api.NewBaseApi[models.Role](
		r.storage,
		"/roles",
		"Role",
		"#/components/schemas/CreateRole",
		"#/components/schemas/UpdateRole",
	)

	getAllRoute := api.GetAllRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("roles.view"),
	)
	getOneRoute := api.GetOneRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("roles.view"),
	)
	createRoute := api.CreateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("roles.create"),
	)
	updateRoute := api.UpdateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("roles.update"),
	)
	deleteRoute := api.DeleteRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("roles.delete"),
	)

	return []routing.Route{
		getAllRoute,
		getOneRoute,
		createRoute,
		updateRoute,
		deleteRoute,
	}
}
