package authentication

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AuthenticationRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
	session    *session.Store
}

func NewAuthenticationRouter(storage storage.Storage, middleware middleware.Middleware, session *session.Store) Router {
	return &AuthenticationRouter{
		storage:    storage,
		middleware: middleware,
		session:    session,
	}
}

func (r *AuthenticationRouter) LoadRoutes() []routing.Route {
	checkRoute := r.CheckRoute()
	loginRoute := r.LoginRoute()
	registerRoute := r.RegisterRoute()
	permissionsRoute := r.PermissionsRoute()
	logoutRoute := r.LogoutRoute()

	return []routing.Route{
		checkRoute,
		loginRoute,
		registerRoute,
		permissionsRoute,
		logoutRoute,
	}
}
