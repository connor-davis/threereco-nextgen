package services

import (
	"github.com/connor-davis/threereco-nextgen/internal/services/materials"
	"github.com/connor-davis/threereco-nextgen/internal/services/organizations"
	"github.com/connor-davis/threereco-nextgen/internal/services/products"
	"github.com/connor-davis/threereco-nextgen/internal/services/roles"
	"github.com/connor-davis/threereco-nextgen/internal/services/transactions"
	"github.com/connor-davis/threereco-nextgen/internal/services/users"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

// Services aggregates various service dependencies used throughout the application.
// It provides centralized access to storage, user management, role management,
// and organization management services.
type Services struct {
	Storage       *storage.Storage
	Users         *users.UsersService
	Roles         *roles.RolesService
	Organizations *organizations.OrganizationsService
	Materials     *materials.MaterialsService
	Products      *products.ProductsService
	Transactions  *transactions.TransactionsService
}

// NewServices initializes and returns a new Services struct, wiring together
// the provided storage with user, role, and organization services.
// It ensures that all dependent services share the same storage backend.
//
// Parameters:
//
//	storage - a pointer to the Storage instance used by all services.
//
// Returns:
//
//	A pointer to the newly constructed Services struct.
func NewServices(storage *storage.Storage) *Services {
	users := users.NewUsersService(storage)
	roles := roles.NewRolesService(storage)
	organizations := organizations.NewOrganizationsService(storage)

	return &Services{
		Storage:       storage,
		Users:         users,
		Roles:         roles,
		Organizations: organizations,
	}
}
