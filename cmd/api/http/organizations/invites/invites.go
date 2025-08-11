package invites

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/services"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type OrganizationsInvitesRouter struct {
	Storage    *storage.Storage
	Sessions   *session.Store
	Services   *services.Services
	Middleware *middleware.Middleware
}

func NewOrganizationsInvitesRouter(storage *storage.Storage, sessions *session.Store, services *services.Services, middleware *middleware.Middleware) *OrganizationsInvitesRouter {
	return &OrganizationsInvitesRouter{
		Storage:    storage,
		Sessions:   sessions,
		Services:   services,
		Middleware: middleware,
	}
}

func (r *OrganizationsInvitesRouter) InitializeRoutes() []routing.Route {
	acceptInvite := r.AcceptInvite()
	sendInvite := r.SendInvite()

	return []routing.Route{
		acceptInvite,
		sendInvite,
	}
}
