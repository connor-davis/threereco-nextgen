package businesses

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/api"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

type BusinessesRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewBusinessesRouter(storage storage.Storage, middleware middleware.Middleware) Router {
	return &BusinessesRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *BusinessesRouter) LoadRoutes() []routing.Route {
	assignRoleApi := api.NewAssignmentApi[models.Business, models.Role](
		r.storage,
		"/businesses",
		"Business",
		"Role",
	)

	assignRoleRoute := assignRoleApi.AssignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.roles.assign"),
	)
	unassignRoleRoute := assignRoleApi.UnassignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.roles.unassign"),
	)
	listRolesRoute := assignRoleApi.ListRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.roles.view"),
	)

	assignUserApi := api.NewAssignmentApi[models.Business, models.User](
		r.storage,
		"/businesses",
		"Business",
		"User",
	)

	assignUserRoute := assignUserApi.AssignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.users.assign"),
	)
	unassignUserRoute := assignUserApi.UnassignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.users.unassign"),
	)
	listUsersRoute := assignUserApi.ListRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.users.view"),
	)

	api := api.NewBaseApi[models.Business](
		r.storage,
		"/businesses",
		"Business",
		"#/components/schemas/CreateBusiness",
		"#/components/schemas/UpdateBusiness",
	)

	getAllRoute := api.GetAllRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.view"),
	)
	getOneRoute := api.GetOneRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.view"),
	)
	createRoute := api.CreateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.create"),
	)
	updateRoute := api.UpdateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.update"),
	)
	deleteRoute := api.DeleteRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("businesses.delete"),
	)

	return []routing.Route{
		assignRoleRoute,
		unassignRoleRoute,
		listRolesRoute,
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
