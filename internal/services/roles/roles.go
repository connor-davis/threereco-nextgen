package roles

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// RolesService provides methods for managing roles within the application.
// It interacts with the underlying storage layer to perform CRUD operations on roles.
type RolesService struct {
	Storage *storage.Storage
}

// NewRolesService creates and returns a new instance of RolesService using the provided storage.
// It initializes the RolesService with the given storage backend for role management operations.
//
// Parameters:
//
//	storage - a pointer to the Storage instance used for persisting role data.
//
// Returns:
//
//	A pointer to the newly created RolesService.
func NewRolesService(storage *storage.Storage) *RolesService {
	return &RolesService{
		Storage: storage,
	}
}

// Upsert inserts or updates a role in the database. If the role already exists, it is updated;
// otherwise, a new role is created. The operation is performed within the context of the provided
// auditId for auditing purposes. Additional GORM clauses can be specified to customize the query.
// Returns an error if the operation fails.
func (s *RolesService) Upsert(auditId uuid.UUID, role *models.Role, clauses ...clause.Expression) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Clauses(clauses...).Assign(&role).FirstOrCreate(&role).Error; err != nil {
		return err
	}

	return nil
}

// Delete removes a role from the database by its ID.
// It also sets the audit user ID for tracking who performed the deletion.
// Returns an error if the deletion fails.
//
// Parameters:
//   - auditId: UUID of the user performing the audit.
//   - id: String identifier of the role to be deleted.
//
// Returns:
//   - error: Non-nil if the deletion fails, nil otherwise.
func (s *RolesService) Delete(auditId uuid.UUID, id string) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Where("id = ?", id).Delete(&models.Role{}).Error; err != nil {
		return err
	}

	return nil
}

// GetById retrieves a role from the database by its unique identifier.
// It returns a pointer to the Role model and an error if the operation fails.
//
// Parameters:
//   - id: The unique identifier of the role to retrieve.
//
// Returns:
//   - *models.Role: Pointer to the retrieved Role model.
//   - error: Error encountered during the retrieval, or nil if successful.
func (s *RolesService) GetById(id string) (*models.Role, error) {
	var role models.Role

	if err := s.Storage.Postgres.Where("id = ?", id).Find(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// GetAll retrieves all roles from the database, applying the provided GORM clause expressions as filters.
// It returns a slice of models.Role and an error if the query fails.
//
// Parameters:
//
//	clauses - Optional GORM clause expressions to customize the query.
//
// Returns:
//
//	[]models.Role - A slice containing the retrieved roles.
//	error         - An error if the database query fails, otherwise nil.
func (s *RolesService) GetAll(clauses ...clause.Expression) ([]models.Role, error) {
	var roles []models.Role

	if err := s.Storage.Postgres.Clauses(clauses...).Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}
