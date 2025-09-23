package transactions

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/api"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

type TransactionsRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewTransactionsRouter(storage storage.Storage, middleware middleware.Middleware) Router {
	return &TransactionsRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *TransactionsRouter) LoadRoutes() []routing.Route {
	assignMaterialsApi := api.NewAssignmentApi[models.Transaction, models.TransactionMaterial](
		r.storage,
		"/transactions",
		"Transaction",
		"Material",
	)

	assignMaterialRoute := assignMaterialsApi.AssignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.materials.assign"),
	)
	unassignMaterialRoute := assignMaterialsApi.UnassignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.materials.unassign"),
	)
	listMaterialsRoute := assignMaterialsApi.ListRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.materials.view"),
		r.middleware.Policies(models.TransactionsPolicy, models.SystemPolicy),
	)

	api := api.NewBaseApi[models.Transaction](
		r.storage,
		"/transactions",
		"Transaction",
		"#/components/schemas/CreateTransaction",
		"#/components/schemas/UpdateTransaction",
	)

	getAllRoute := api.GetAllRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.view"),
		r.middleware.Policies(models.TransactionsPolicy, models.SystemPolicy),
	)
	getOneRoute := api.GetOneRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.view"),
		r.middleware.Policies(models.TransactionsPolicy, models.SystemPolicy),
	)
	createRoute := api.CreateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.create"),
	)
	updateRoute := api.UpdateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.update"),
	)
	deleteRoute := api.DeleteRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("transactions.delete"),
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
