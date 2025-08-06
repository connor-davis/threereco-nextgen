package sessions

import (
	"time"

	"github.com/connor-davis/threereco-nextgen/env"
	"github.com/gofiber/fiber/v2/middleware/session"
	fiberPg "github.com/gofiber/storage/postgres/v2"
)

// NewSessions initializes and returns a new session.Store configured to use PostgreSQL as the storage backend.
// It sets up session management with specific cookie settings, including domain, path, security, and expiration.
// The session data is stored in the "sessions" table using the provided PostgreSQL connection URI.
// Returns a pointer to the configured session.Store.
func NewSessions() *session.Store {
	return session.New(session.Config{
		Storage: fiberPg.New(fiberPg.Config{
			Table:         "sessions",
			ConnectionURI: string(env.POSTGRES_DSN),
		}),
		KeyLookup:         "cookie:threereco_session",
		CookieDomain:      string(env.COOKIE_DOMAIN),
		CookiePath:        "/",
		CookieSecure:      true,
		CookieSameSite:    "Strict",
		CookieSessionOnly: false,
		CookieHTTPOnly:    false,
		Expiration:        1 * time.Hour,
	})
}
