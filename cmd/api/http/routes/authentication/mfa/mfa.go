package mfa

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type MfaRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
	session    *session.Store
}

func NewMfaRouter(storage storage.Storage, middleware middleware.Middleware, session *session.Store) Router {
	return &MfaRouter{
		storage:    storage,
		middleware: middleware,
		session:    session,
	}
}

func (r *MfaRouter) LoadRoutes() []routing.Route {
	enableRoute := r.EnableRoute()
	verifyRoute := r.VerifyRoute()

	return []routing.Route{
		enableRoute,
		verifyRoute,
	}
}
