package collections

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/api"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

type CollectionsRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewCollectionsRouter(storage storage.Storage, middleware middleware.Middleware) Router {
	return &CollectionsRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *CollectionsRouter) LoadRoutes() []routing.Route {
	assignMaterialsApi := api.NewAssignmentApi[models.Collection, models.CollectionMaterial](
		r.storage,
		"/collections",
		"Collection",
		"Material",
	)

	assignMaterialRoute := assignMaterialsApi.AssignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.materials.assign"),
	)
	unassignMaterialRoute := assignMaterialsApi.UnassignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.materials.unassign"),
	)
	listMaterialsRoute := assignMaterialsApi.ListRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.materials.view"),
		r.middleware.Policies(models.CollectionsPolicy, models.SystemPolicy),
	)

	api := api.NewBaseApi[models.Collection](
		r.storage,
		"/collections",
		"Collection",
		"#/components/schemas/CreateCollection",
		"#/components/schemas/UpdateCollection",
	)

	getAllRoute := api.GetAllRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.view"),
		r.middleware.Policies(models.CollectionsPolicy, models.SystemPolicy),
	)
	getOneRoute := api.GetOneRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.view"),
		r.middleware.Policies(models.CollectionsPolicy, models.SystemPolicy),
	)
	createRoute := api.CreateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.create"),
	)
	updateRoute := api.UpdateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.update"),
	)
	deleteRoute := api.DeleteRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("collections.delete"),
	)

	return []routing.Route{
		assignMaterialRoute,
		unassignMaterialRoute,
		listMaterialsRoute,
		getAllRoute,
		getOneRoute,
		createRoute,
		updateRoute,
		deleteRoute,
	}
}
