package transactions

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/api"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

type TransactionMaterialsRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewTransactionMaterialsRouter(storage storage.Storage, middleware middleware.Middleware) Router {
	return &TransactionMaterialsRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *TransactionMaterialsRouter) LoadRoutes() []routing.Route {
	api := api.NewBaseApi[models.TransactionMaterial](
		r.storage,
		"/transactions/materials",
		"Transaction Material",
		"#/components/schemas/CreateTransactionMaterial",
		"#/components/schemas/UpdateTransactionMaterial",
	)

	getAllRoute := api.GetAllRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.materials.view"),
	)
	getOneRoute := api.GetOneRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.materials.view"),
	)
	createRoute := api.CreateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.materials.create"),
	)
	updateRoute := api.UpdateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.materials.update"),
	)
	deleteRoute := api.DeleteRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.materials.delete"),
	)

	return []routing.Route{
		getAllRoute,
		getOneRoute,
		createRoute,
		updateRoute,
		deleteRoute,
	}
}
