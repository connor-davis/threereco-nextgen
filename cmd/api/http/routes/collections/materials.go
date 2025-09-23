package collections

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/api"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

type CollectionMaterialsRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewCollectionMaterialsRouter(storage storage.Storage, middleware middleware.Middleware) Router {
	return &CollectionMaterialsRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *CollectionMaterialsRouter) LoadRoutes() []routing.Route {
	api := api.NewBaseApi[models.CollectionMaterial](
		r.storage,
		"/collections/materials",
		"Collection Material",
		"#/components/schemas/CreateCollectionMaterial",
		"#/components/schemas/UpdateCollectionMaterial",
	)

	getAllRoute := api.GetAllRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.materials.view"),
	)
	getOneRoute := api.GetOneRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.materials.view"),
	)
	createRoute := api.CreateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.materials.create"),
	)
	updateRoute := api.UpdateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.materials.update"),
	)
	deleteRoute := api.DeleteRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.materials.delete"),
	)

	return []routing.Route{
		getAllRoute,
		getOneRoute,
		createRoute,
		updateRoute,
		deleteRoute,
	}
}
