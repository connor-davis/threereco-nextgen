package http

import (
	"fmt"
	"regexp"

	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/routes/authentication"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/routes/authentication/mfa"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/routes/businesses"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/routes/collections"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/routes/materials"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/routes/roles"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/routes/transactions"
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/routes/users"
	"github.com/connor-davis/threereco-nextgen/common"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/routing/schemas"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type HttpRouter interface {
	InitializeRoutes(router fiber.Router)
	InitializeOpenAPI() *openapi3.T
}

type httpRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
	session    *session.Store
	routes     []routing.Route
}

func NewHttpRouter(storage storage.Storage, middleware middleware.Middleware, session *session.Store) HttpRouter {
	mfaRouter := mfa.NewMfaRouter(storage, middleware, session)
	mfaRoutes := mfaRouter.LoadRoutes()

	authenticationRouter := authentication.NewAuthenticationRouter(storage, middleware, session)
	authenticationRoutes := authenticationRouter.LoadRoutes()

	usersRouter := users.NewUsersRouter(storage, middleware)
	usersRoutes := usersRouter.LoadRoutes()

	rolesRouter := roles.NewRolesRouter(storage, middleware)
	rolesRoutes := rolesRouter.LoadRoutes()

	materialsRouter := materials.NewMaterialsRouter(storage, middleware)
	materialsRoutes := materialsRouter.LoadRoutes()

	collectionMaterialsRouter := collections.NewCollectionMaterialsRouter(storage, middleware)
	collectionMaterialsRoutes := collectionMaterialsRouter.LoadRoutes()

	collectionsRouter := collections.NewCollectionsRouter(storage, middleware)
	collectionsRoutes := collectionsRouter.LoadRoutes()

	transactionMaterialsRouter := transactions.NewTransactionMaterialsRouter(storage, middleware)
	transactionMaterialsRoutes := transactionMaterialsRouter.LoadRoutes()

	transactionsRouter := transactions.NewTransactionsRouter(storage, middleware)
	transactionsRoutes := transactionsRouter.LoadRoutes()

	businessesRouter := businesses.NewBusinessesRouter(storage, middleware)
	businessesRoutes := businessesRouter.LoadRoutes()

	routes := []routing.Route{}

	routes = append(routes, mfaRoutes...)
	routes = append(routes, authenticationRoutes...)
	routes = append(routes, usersRoutes...)
	routes = append(routes, rolesRoutes...)
	routes = append(routes, materialsRoutes...)
	routes = append(routes, collectionMaterialsRoutes...)
	routes = append(routes, collectionsRoutes...)
	routes = append(routes, transactionMaterialsRoutes...)
	routes = append(routes, transactionsRoutes...)
	routes = append(routes, businessesRoutes...)

	return &httpRouter{
		storage:    storage,
		middleware: middleware,
		session:    session,
		routes:     routes,
	}
}

func (h *httpRouter) InitializeRoutes(router fiber.Router) {
	for _, route := range h.routes {
		path := regexp.MustCompile(`\{([^}]+)\}`).ReplaceAllString(route.Path, ":$1")

		switch route.Method {
		case routing.GET:
			router.Get(path, append(route.Middlewares, route.Handler)...)
		case routing.POST:
			router.Post(path, append(route.Middlewares, route.Handler)...)
		case routing.PUT:
			router.Put(path, append(route.Middlewares, route.Handler)...)
		case routing.DELETE:
			router.Delete(path, append(route.Middlewares, route.Handler)...)
		}
	}
}

func (h *httpRouter) InitializeOpenAPI() *openapi3.T {
	paths := openapi3.NewPaths()

	bodies := openapi3.RequestBodies{
		"LoginPayload":     schemas.LoginPayloadSchema,
		"RegisterPayload":  schemas.RegisterPayloadSchema,
		"VerifyMfaPayload": schemas.VerifyMfaPayloadSchema,
	}

	schemas := openapi3.Schemas{
		"SuccessResponse":            schemas.SuccessSchema,
		"ErrorResponse":              schemas.ErrorSchema,
		"Query":                      schemas.QuerySchema,
		"Pagination":                 schemas.PaginationSchema,
		"Address":                    schemas.AddressSchema,
		"BankDetails":                schemas.BankDetailsSchema,
		"User":                       schemas.UserSchema,
		"Users":                      schemas.UsersSchema,
		"AssignUser":                 schemas.AssignUserSchema,
		"AssignUsers":                schemas.AssignUsersSchema,
		"CreateUser":                 schemas.CreateUserSchema,
		"UpdateUser":                 schemas.UpdateUserSchema,
		"Role":                       schemas.RoleSchema,
		"Roles":                      schemas.RolesSchema,
		"AssignRole":                 schemas.AssignRoleSchema,
		"AssignRoles":                schemas.AssignRolesSchema,
		"CreateRole":                 schemas.CreateRoleSchema,
		"UpdateRole":                 schemas.UpdateRoleSchema,
		"Material":                   schemas.MaterialSchema,
		"Materials":                  schemas.MaterialsSchema,
		"AssignMaterial":             schemas.AssignMaterialSchema,
		"AssignMaterials":            schemas.AssignMaterialsSchema,
		"CreateMaterial":             schemas.CreateMaterialSchema,
		"UpdateMaterial":             schemas.UpdateMaterialSchema,
		"Collection":                 schemas.CollectionSchema,
		"Collections":                schemas.CollectionsSchema,
		"CreateCollection":           schemas.CreateCollectionSchema,
		"UpdateCollection":           schemas.UpdateCollectionSchema,
		"CollectionMaterial":         schemas.CollectionMaterialSchema,
		"CollectionMaterials":        schemas.CollectionMaterialsSchema,
		"AssignCollectionMaterial":   schemas.AssignCollectionMaterialSchema,
		"AssignCollectionMaterials":  schemas.AssignCollectionMaterialsSchema,
		"CreateCollectionMaterial":   schemas.CreateCollectionMaterialSchema,
		"UpdateCollectionMaterial":   schemas.UpdateCollectionMaterialSchema,
		"Transaction":                schemas.TransactionSchema,
		"Transactions":               schemas.TransactionsSchema,
		"CreateTransaction":          schemas.CreateTransactionSchema,
		"UpdateTransaction":          schemas.UpdateTransactionSchema,
		"TransactionMaterial":        schemas.TransactionMaterialSchema,
		"TransactionMaterials":       schemas.TransactionMaterialsSchema,
		"AssignTransactionMaterial":  schemas.AssignTransactionMaterialSchema,
		"AssignTransactionMaterials": schemas.AssignTransactionMaterialsSchema,
		"CreateTransactionMaterial":  schemas.CreateTransactionMaterialSchema,
		"UpdateTransactionMaterial":  schemas.UpdateTransactionMaterialSchema,
		"Business":                   schemas.BusinessSchema,
		"Businesses":                 schemas.BusinessesSchema,
		"AssignBusiness":             schemas.AssignBusinessSchema,
		"AssignBusinesses":           schemas.AssignBusinessesSchema,
		"CreateBusiness":             schemas.CreateBusinessSchema,
		"UpdateBusiness":             schemas.UpdateBusinessSchema,
		"Permission":                 schemas.PermissionSchema,
		"Permissions":                schemas.PermissionsSchema,
		"PermissionGroup":            schemas.PermissionGroupSchema,
		"PermissionGroups":           schemas.PermissionGroupsSchema,
	}

	for _, route := range h.routes {
		pathItem := &openapi3.PathItem{}

		switch route.Method {
		case routing.GET:
			pathItem.Get = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				Responses:   route.Responses,
			}
		case routing.POST:
			pathItem.Post = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.PUT:
			pathItem.Put = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.DELETE:
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
			case routing.GET:
				existingPathItem.Get = pathItem.Get
			case routing.POST:
				existingPathItem.Post = pathItem.Post
			case routing.PUT:
				existingPathItem.Put = pathItem.Put
			case routing.DELETE:
				existingPathItem.Delete = pathItem.Delete
			}
		} else {
			paths.Set(path, pathItem)
		}
	}

	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   common.EnvString("APP_NAME", "Dynamic CRUD API"),
			Version: common.EnvString("APP_VERSION", "1.0.0"),
		},
		Servers: openapi3.Servers{
			{
				URL:         fmt.Sprintf("http://localhost:%s", common.EnvString("APP_PORT", "6173")),
				Description: "Development",
			},
			{
				URL:         common.EnvString("APP_BASE_URL", "https://example.com"),
				Description: "Production",
			},
		},
		Tags:  openapi3.Tags{},
		Paths: paths,
		Components: &openapi3.Components{
			Schemas:       schemas,
			RequestBodies: bodies,
		},
	}
}
