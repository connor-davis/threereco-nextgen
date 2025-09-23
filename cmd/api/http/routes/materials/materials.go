package materials

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/api"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

type MaterialsRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewMaterialsRouter(storage storage.Storage, middleware middleware.Middleware) Router {
	return &MaterialsRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *MaterialsRouter) LoadRoutes() []routing.Route {
	api := api.NewBaseApi[models.Material](
		r.storage,
		"/materials",
		"Material",
		"#/components/schemas/CreateMaterial",
		"#/components/schemas/UpdateMaterial",
	)

	getAllRoute := api.GetAllRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("materials.view"),
	)
	getOneRoute := api.GetOneRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("materials.view"),
	)
	createRoute := api.CreateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("materials.create"),
	)
	updateRoute := api.UpdateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("materials.update"),
	)
	deleteRoute := api.DeleteRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("materials.delete"),
	)

	return []routing.Route{
		getAllRoute,
		getOneRoute,
		createRoute,
		updateRoute,
		deleteRoute,
	}
}
