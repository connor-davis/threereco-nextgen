package http

import (
	"fmt"
	"regexp"

	"github.com/connor-davis/threereco-nextgen/cmd/api/http/authentication"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/materials"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/notifications"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/organizations"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/organizations/invites"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/products"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/roles"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/transactions"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/users"
	"github.com/connor-davis/threereco-nextgen/env"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/connor-davis/threereco-nextgen/internal/services"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// HttpRouter encapsulates the dependencies and configuration required to set up HTTP routing.
// It includes references to storage, session management, service layer, middleware, and route definitions.
type HttpRouter struct {
	Storage    *storage.Storage
	Sessions   *session.Store
	Services   *services.Services
	Middleware *middleware.Middleware
	Routes     []routing.Route
}

// NewHttpRouter creates and initializes a new HttpRouter instance.
// It sets up middleware and authentication routes using the provided storage,
// session store, and services. The returned HttpRouter contains all configured
// routes and dependencies required for handling HTTP requests.
//
// Parameters:
//   - storage: Pointer to the application's storage layer.
//   - sessions: Pointer to the session store for managing user sessions.
//   - services: Pointer to the application's service layer.
//
// Returns:
//   - *HttpRouter: A pointer to the initialized HttpRouter.
func NewHttpRouter(storage *storage.Storage, sessions *session.Store, services *services.Services, middleware *middleware.Middleware) *HttpRouter {
	authenticationRouter := authentication.NewAuthenticationRouter(storage, sessions, services, middleware)
	authenticationRoutes := authenticationRouter.InitializeRoutes()

	usersRouter := users.NewUsersRouter(storage, sessions, services, middleware)
	usersRoutes := usersRouter.InitializeRoutes()

	rolesRouter := roles.NewRolesRouter(storage, sessions, services, middleware)
	rolesRoutes := rolesRouter.InitializeRoutes()

	organizationsRouter := organizations.NewOrganizationsRouter(storage, sessions, services, middleware)
	organizationsRoutes := organizationsRouter.InitializeRoutes()

	invitesRouter := invites.NewOrganizationsInvitesRouter(storage, sessions, services, middleware)
	invitesRoutes := invitesRouter.InitializeRoutes()

	materialsRouter := materials.NewMaterialsRouter(storage, sessions, services, middleware)
	materialsRoutes := materialsRouter.InitializeRoutes()

	productsRouter := products.NewProductsRouter(storage, sessions, services, middleware)
	productsRoutes := productsRouter.InitializeRoutes()

	transactionsRouter := transactions.NewTransactionsRouter(storage, sessions, services, middleware)
	transactionsRoutes := transactionsRouter.InitializeRoutes()

	notificationsRouter := notifications.NewNotificationsRouter(storage, sessions, services, middleware)
	notificationsRoutes := notificationsRouter.InitializeRoutes()

	routes := []routing.Route{}

	routes = append(routes, authenticationRoutes...)
	routes = append(routes, usersRoutes...)
	routes = append(routes, rolesRoutes...)
	routes = append(routes, organizationsRoutes...)
	routes = append(routes, invitesRoutes...)
	routes = append(routes, materialsRoutes...)
	routes = append(routes, productsRoutes...)
	routes = append(routes, transactionsRoutes...)
	routes = append(routes, notificationsRoutes...)

	return &HttpRouter{
		Storage:    storage,
		Sessions:   sessions,
		Services:   services,
		Middleware: middleware,
		Routes:     routes,
	}
}

// InitializeRoutes registers all HTTP routes defined in the HttpRouter with the provided Fiber router.
// It converts route paths from the format "{param}" to Fiber's ":param" syntax using regular expressions.
// For each route, it attaches the specified middlewares and handler to the corresponding HTTP method.
// Supported methods include GET, POST, PUT, PATCH, OPTIONS, and DELETE.
func (r *HttpRouter) InitializeRoutes(router fiber.Router) {
	for _, route := range r.Routes {
		path := regexp.MustCompile(`\{([^}]+)\}`).ReplaceAllString(route.Path, ":$1")

		switch route.Method {
		case routing.GetMethod:
			router.Get(path, append(route.Middlewares, route.Handler)...)
		case routing.PostMethod:
			router.Post(path, append(route.Middlewares, route.Handler)...)
		case routing.PutMethod:
			router.Put(path, append(route.Middlewares, route.Handler)...)
		case routing.PatchMethod:
			router.Patch(path, append(route.Middlewares, route.Handler)...)
		case routing.OptionsMethod:
			router.Options(path, append(route.Middlewares, route.Handler)...)
		case routing.DeleteMethod:
			router.Delete(path, append(route.Middlewares, route.Handler)...)
		}
	}
}

// InitializeOpenAPI generates and returns an OpenAPI 3 specification for the HTTP router.
// It iterates through the defined routes, constructs OpenAPI PathItem objects for each HTTP method,
// and sets up the API paths, operations, and components (schemas, servers, etc.).
// The resulting OpenAPI specification includes metadata such as title, version, server URLs, and
// reusable schemas for success and error responses.
//
// Returns:
//
//	*openapi3.T: The constructed OpenAPI 3 specification for the API.
func (h *HttpRouter) InitializeOpenAPI() *openapi3.T {
	paths := openapi3.NewPaths()

	for _, route := range h.Routes {
		pathItem := &openapi3.PathItem{}

		switch route.Method {
		case routing.GetMethod:
			pathItem.Get = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				Responses:   route.Responses,
			}
		case routing.PostMethod:
			pathItem.Post = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.PutMethod:
			pathItem.Put = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.PatchMethod:
			pathItem.Patch = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.OptionsMethod:
			pathItem.Options = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.DeleteMethod:
			pathItem.Delete = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				Responses:   route.Responses,
			}
		}

		path := fmt.Sprintf("/api%s", route.Path)

		existingPathItem := paths.Find(path)

		if existingPathItem != nil {
			switch route.Method {
			case routing.GetMethod:
				existingPathItem.Get = pathItem.Get
			case routing.PostMethod:
				existingPathItem.Post = pathItem.Post
			case routing.PutMethod:
				existingPathItem.Put = pathItem.Put
			case routing.DeleteMethod:
				existingPathItem.Delete = pathItem.Delete
			}
		} else {
			paths.Set(path, pathItem)
		}
	}

	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "Thusa One API",
			Version: "1.0.0",
		},
		Servers: openapi3.Servers{
			{
				URL:         fmt.Sprintf("http://localhost:%s", string(env.PORT)),
				Description: "Development",
			},
			{
				URL:         "https://one.thusa.co.za",
				Description: "Production",
			},
		},
		Tags:  openapi3.Tags{},
		Paths: paths,
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"SuccessResponse":           schemas.SuccessResponseSchema,
				"ErrorResponse":             schemas.ErrorResponseSchema,
				"User":                      schemas.UserSchema,
				"CreateUserPayload":         schemas.CreateUserPayloadSchema,
				"UpdateUserPayload":         schemas.UpdateUserPayloadSchema,
				"Role":                      schemas.RoleSchema,
				"CreateRolePayload":         schemas.CreateRolePayloadSchema,
				"UpdateRolePayload":         schemas.UpdateRolePayloadSchema,
				"Organization":              schemas.OrganizationSchema,
				"CreateOrganizationPayload": schemas.CreateOrganizationPayloadSchema,
				"UpdateOrganizationPayload": schemas.UpdateOrganizationPayloadSchema,
				"AuditLog":                  schemas.AuditLogSchema,
				"MfaVerifyPayload":          schemas.MfaVerifyPayloadSchema,
				"LoginPayload":              schemas.LoginPayloadSchema,
				"CreateMaterialPayload":     schemas.CreateMaterialPayloadSchema,
				"UpdateMaterialPayload":     schemas.UpdateMaterialPayloadSchema,
				"Material":                  schemas.MaterialSchema,
				"CreateProductPayload":      schemas.CreateProductPayloadSchema,
				"UpdateProductPayload":      schemas.UpdateProductPayloadSchema,
				"Product":                   schemas.ProductSchema,
				"CreateTransactionPayload":  schemas.CreateTransactionPayloadSchema,
				"UpdateTransactionPayload":  schemas.UpdateTransactionPayloadSchema,
				"Transaction":               schemas.TransactionSchema,
			},
		},
	}
}
