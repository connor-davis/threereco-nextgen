package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/common"
	"github.com/connor-davis/threereco-nextgen/internal/sessions"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	storage := storage.New()

	err := storage.Migrate()

	if err != nil {
		log.Fatalf("🔥 Failed to migrate database: %v", err)
	}

	err = storage.SeedAdmin()

	if err != nil {
		log.Fatalf("🔥 Failed to seed admin user: %v", err)
	}

	err = storage.SeedDefaultBusiness()

	if err != nil {
		log.Fatalf("🔥 Failed to seed default business: %v", err)
	}

	session := sessions.New()
	middleware := middleware.New(storage, session)

	app := fiber.New(fiber.Config{
		AppName:       common.EnvString("APP_NAME", "Dynamic CRUD API"),
		ServerHeader:  common.EnvString("APP_HEADER", "Dynamic-CRUD"),
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		StrictRouting: true,
		CaseSensitive: true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: fmt.Sprintf(
			"%s,%s",
			common.EnvString("APP_BASE_URL", "http://localhost:3000"),
			"http://localhost:3000",
		),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Africa/Johannesburg",
	}))

	api := app.Group("/api")

	httpRouter := http.NewHttpRouter(storage, middleware, session)
	httpRouter.InitializeRoutes(api)

	openapi := httpRouter.InitializeOpenAPI()

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
		return c.Status(fiber.StatusOK).JSON(openapi)
	})

	api.Get("/api-doc", func(c *fiber.Ctx) error {
		html, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: func() string {
				if common.EnvString("APP_ENV", "development") == "production" {
					return fmt.Sprintf("%s/api/api-spec", common.EnvString("APP_BASE_URL", "https://example.com"))
				}

				return fmt.Sprintf("http://localhost:%s/api/api-spec", common.EnvString("APP_PORT", "6173"))
			}(),
			Theme:  scalar.ThemeDefault,
			Layout: scalar.LayoutModern,
			BaseServerURL: func() string {
				if common.EnvString("APP_ENV", "development") == "production" {
					return common.EnvString("APP_BASE_URL", "https://example.com")
				}

				return fmt.Sprintf("http://localhost:%s", common.EnvString("APP_PORT", "6173"))
			}(),
			DarkMode: true,
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Type("html").SendString(html)
	})

	log.Printf("✅ Starting API on port %s...", common.EnvString("APP_PORT", "6173"))

	if err := app.Listen(fmt.Sprintf(":%s", common.EnvString("APP_PORT", "6173"))); err != nil {
		log.Printf("🔥 Failed to start server: %v", err)
	}
}
