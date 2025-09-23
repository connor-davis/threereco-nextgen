package users

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/api"
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

type UsersRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewUsersRouter(storage storage.Storage, middleware middleware.Middleware) Router {
	return &UsersRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *UsersRouter) LoadRoutes() []routing.Route {
	assignRoleApi := api.NewAssignmentApi[models.User, models.Role](
		r.storage,
		"/users",
		"User",
		"Role",
	)

	assignRoleRoute := assignRoleApi.AssignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("users.roles.assign"),
	)
	unassignRoleRoute := assignRoleApi.UnassignRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("users.roles.unassign"),
	)
	listRolesRoute := assignRoleApi.ListRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("users.roles.view"),
	)

	api := api.NewBaseApi[models.User](
		r.storage,
		"/users",
		"User",
		"#/components/schemas/CreateUser",
		"#/components/schemas/UpdateUser",
	)

	getAllRoute := api.GetAllRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("users.view"),
	)
	getOneRoute := api.GetOneRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("users.view"),
	)
	createRoute := api.CreateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("users.create"),
	)
	updateRoute := api.UpdateRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("users.update.any", "users.update.self"),
	)
	deleteRoute := api.DeleteRoute(
		r.middleware.Authenticated(),
		r.middleware.Authorized("users.delete.any", "users.delete.self"),
	)

	return []routing.Route{
		assignRoleRoute,
		unassignRoleRoute,
		listRolesRoute,
		getAllRoute,
		getOneRoute,
		createRoute,
		updateRoute,
		deleteRoute,
	}
}
