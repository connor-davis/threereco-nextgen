package authentication

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/services"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// AuthenticationRouter encapsulates dependencies required for handling authentication-related HTTP routes.
// It provides access to storage, session management, service layer, and middleware components.
type AuthenticationRouter struct {
	Storage    *storage.Storage
	Sessions   *session.Store
	Services   *services.Services
	Middleware *middleware.Middleware
}

// NewAuthenticationRouter creates and returns a new instance of AuthenticationRouter,
// initializing it with the provided storage, session store, services, and middleware.
// This function is typically used to set up authentication-related HTTP routing
// with the necessary dependencies injected.
//
// Parameters:
//   - storage: Pointer to the application's storage layer.
//   - sessions: Pointer to the session store for managing user sessions.
//   - services: Pointer to the services container providing business logic.
//   - middleware: Pointer to the middleware manager for HTTP request handling.
//
// Returns:
//   - *AuthenticationRouter: A pointer to the newly constructed AuthenticationRouter.
func NewAuthenticationRouter(storage *storage.Storage, sessions *session.Store, services *services.Services, middleware *middleware.Middleware) *AuthenticationRouter {
	return &AuthenticationRouter{
		Storage:    storage,
		Sessions:   sessions,
		Services:   services,
		Middleware: middleware,
	}
}

// InitializeRoutes sets up and returns the authentication-related routes for the router.
// Specifically, it includes routes for handling Microsoft authentication redirects and callbacks.
// Returns a slice of routing.Route containing the configured routes.
func (r *AuthenticationRouter) InitializeRoutes() []routing.Route {
	// General routes
	checkRoute := r.CheckRoute()
	logoutRoute := r.LogoutRoute()
	loginRoute := r.LoginRoute()
	permissionsRoute := r.PermissionsRoute()
	organizationsRoute := r.OrganizationsRoute()

	// MFA routes
	mfaEnableRoute := r.MfaEnableRoute()
	mfaVerifyRoute := r.MfaVerifyRoute()

	return []routing.Route{
		checkRoute,
		logoutRoute,
		loginRoute,
		permissionsRoute,
		organizationsRoute,
		mfaEnableRoute,
		mfaVerifyRoute,
	}
}
