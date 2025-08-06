package routing

import "github.com/gofiber/fiber/v2"

type RouteMethod string

const (
	GetMethod     RouteMethod = "GET"
	PostMethod    RouteMethod = "POST"
	PutMethod     RouteMethod = "PUT"
	PatchMethod   RouteMethod = "PATCH"
	OptionsMethod RouteMethod = "OPTIONS"
	DeleteMethod  RouteMethod = "DELETE"
)

// Route represents an HTTP route configuration, including its method, path, associated middlewares, and handler function.
// It also embeds OpenAPIMetadata for OpenAPI specification support.
//
// Fields:
//   - Method: The HTTP method (e.g., GET, POST) for the route.
//   - Path: The URL path pattern for the route.
//   - Middlewares: A slice of Fiber middleware handlers to be executed before the main handler.
//   - Handler: The main handler function for processing requests to this route.
type Route struct {
	OpenAPIMetadata

	Method      RouteMethod
	Path        string
	Middlewares []fiber.Handler
	Handler     func(*fiber.Ctx) error
}
