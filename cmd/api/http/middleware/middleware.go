package middleware

import (
	"github.com/connor-davis/threereco-nextgen/internal/services"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Middleware encapsulates dependencies required for HTTP middleware operations,
// including storage access, session management, and service interactions.
// It provides a convenient way to inject these dependencies into middleware handlers.
type Middleware struct {
	Storage  *storage.Storage
	Sessions *session.Store
	Services *services.Services
}

// NewMiddleware creates and returns a new Middleware instance, initializing it with the provided
// storage, session store, and services. This function is typically used to set up middleware
// dependencies for handling HTTP requests in the application.
//
// Parameters:
//
//	storage  - Pointer to the application's storage layer.
//	sessions - Pointer to the session store for managing user sessions.
//	services - Pointer to the application's service layer.
//
// Returns:
//
//	A pointer to a newly constructed Middleware instance.
func NewMiddleware(storage *storage.Storage, sessions *session.Store, services *services.Services) *Middleware {
	return &Middleware{
		Storage:  storage,
		Sessions: sessions,
		Services: services,
	}
}
