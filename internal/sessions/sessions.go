package sessions

import (
	"time"

	"github.com/connor-davis/threereco-nextgen/common"
	"github.com/gofiber/fiber/v2/middleware/session"
	fiberPg "github.com/gofiber/storage/postgres/v2"
)

func New() *session.Store {
	return session.New(session.Config{
		Storage: fiberPg.New(fiberPg.Config{
			Table:         "sessions",
			ConnectionURI: common.EnvString("APP_DSN", "host=localhost user=<user> password=<password> dbname=<database> port=5432 sslmode=disable TimeZone=Africa/Johannesburg"),
		}),
		KeyLookup:         common.EnvString("APP_SESSION_KEY", "cookie:threereco_session"),
		CookieDomain:      common.EnvString("APP_DOMAIN", "localhost"),
		CookiePath:        "/",
		CookieSecure:      true,
		CookieSameSite:    "Strict",
		CookieSessionOnly: false,
		CookieHTTPOnly:    false,
		Expiration:        1 * time.Hour,
	})
}
