package main

import (
	"encoding/json"
	"fmt"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/env"
	"github.com/connor-davis/threereco-nextgen/internal/services"
	"github.com/connor-davis/threereco-nextgen/internal/sessions"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// main initializes and starts the Zingfibre Reporting API server.
// It sets up storage, database connections, migrations, and seeds initial data.
// The function configures session management, service dependencies, and the Fiber web application,
// including CORS and logging middleware. It defines API routes for health checks, OpenAPI specification,
// and API documentation, and registers additional routes via the HTTP router.
// Finally, it starts the server on the configured port and logs startup or error messages.
func main() {
	storage := storage.New()

	storage.ConnectPostgres()
	storage.MigratePostgres()
	storage.SeedPostgres()

	sessions := sessions.NewSessions()
	services := services.NewServices(storage)

	middleware := middleware.NewMiddleware(storage, sessions, services)

	app := fiber.New(fiber.Config{
		AppName:      "3rEco API",
		ServerHeader: "3rEco-NextGen-API",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,https://3reco.co.za",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Africa/Johannesburg",
	}))

	httpRouter := http.NewHttpRouter(storage, sessions, services, middleware)

	openapiSpecification := httpRouter.InitializeOpenAPI()

	api := app.Group("/api")

	api.Get(
		"/health",
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "ok",
				"message": "API is running",
			})
		},
	)

	api.Get("/api-spec", func(c *fiber.Ctx) error {
		if string(env.MODE) == "production" {
			return c.Status(fiber.StatusOK).JSON(openapiSpecification)
		}

		return c.Status(fiber.StatusOK).JSON(openapiSpecification)
	})

	api.Get("/api-doc", middleware.Authorized(), func(c *fiber.Ctx) error {
		html, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: func() string {
				if string(env.MODE) == "production" {
					return "https://3reco.co.za/api/api-spec"
				}

				return fmt.Sprintf("http://localhost:%s/api/api-spec", env.PORT)
			}(),
			Theme:  scalar.ThemeDefault,
			Layout: scalar.LayoutModern,
			BaseServerURL: func() string {
				if string(env.MODE) == "production" {
					return "https://3reco.co.za"
				}

				return fmt.Sprintf("http://localhost:%s", env.PORT)
			}(),
			DarkMode: true,
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Type("html").SendString(html)
	})

	httpRouter.InitializeRoutes(api)

	log.Infof("✅ Starting 3rEco API on port %s...", string(env.PORT))

	if err := app.Listen(fmt.Sprintf(":%s", env.PORT)); err != nil {
		log.Errorf("🔥 Failed to start server: %v", err)
	}
}
