package auditlogs

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/services"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AuditLogsRouter struct {
	Storage    *storage.Storage
	Sessions   *session.Store
	Services   *services.Services
	Middleware *middleware.Middleware
}

func NewAuditLogsRouter(storage *storage.Storage, sessions *session.Store, services *services.Services, middleware *middleware.Middleware) *AuditLogsRouter {
	return &AuditLogsRouter{
		Storage:    storage,
		Sessions:   sessions,
		Services:   services,
		Middleware: middleware,
	}
}

func (r *AuditLogsRouter) InitializeRoutes() []routing.Route {
	viewRoute := r.ViewRoute()
	viewByIdRoute := r.ViewByIdRoute()

	return []routing.Route{
		viewRoute,
		viewByIdRoute,
	}
}
