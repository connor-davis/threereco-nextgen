package businesses

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/api"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

type Router struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewBusinessesRouter(storage storage.Storage, middleware middleware.Middleware) IRouter {
	return &Router{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *Router) LoadRoutes() []routing.Route {
	businessesUserAssignmentsApi := api.NewAssignmentApi[models.Business, models.User](
		r.storage,
		"/businesses",
		"Business",
		"User",
	)

	assignUserRoute := businessesUserAssignmentsApi.AssignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.users.assign"),
	)
	unassignUserRoute := businessesUserAssignmentsApi.UnassignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.users.unassign"),
	)
	listUsersRoute := businessesUserAssignmentsApi.ListRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.users.view"),
	)

	businessesApi := api.NewBaseApi[models.Business](
		r.storage,
		"/businesses",
		"Business",
		"#/components/schemas/CreateBusiness",
		"#/components/schemas/UpdateBusiness",
	)

	getAllRoute := businessesApi.GetAllRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.view"),
	)
	getOneRoute := businessesApi.GetOneRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.view"),
	)
	createRoute := businessesApi.CreateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.create"),
	)
	updateRoute := businessesApi.UpdateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.update"),
	)
	deleteRoute := businessesApi.DeleteRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.delete"),
	)

	return []routing.Route{
		assignUserRoute,
		unassignUserRoute,
		listUsersRoute,
		getAllRoute,
		getOneRoute,
		createRoute,
		updateRoute,
		deleteRoute,
	}
}
