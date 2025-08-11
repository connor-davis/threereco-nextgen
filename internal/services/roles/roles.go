package roles

import (
	"errors"

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

// Create inserts a new role record into the database using the provided Role model.
// It associates the operation with the given auditId for auditing purposes.
// Returns an error if the creation fails.
func (s *RolesService) Create(auditId uuid.UUID, role models.CreateRolePayload) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Create(&models.Role{
			Name:        role.Name,
			Description: role.Description,
			Permissions: role.Permissions,
			Organizations: []models.Organization{
				{
					Id: role.OrganizationId,
				},
			},
			ModifiedByUserId: auditId,
		}).Error; err != nil {
		return err
	}

	return nil
}

// Update updates the role with the specified ID in the database.
// It uses the provided auditId for auditing purposes and applies the changes from the given role model.
// Returns an error if the update operation fails.
func (s *RolesService) Update(auditId uuid.UUID, id uuid.UUID, role models.UpdateRolePayload) error {
	var existingRole models.Role

	if err := s.Storage.Postgres.Where("id = $1", id).Find(&existingRole).Error; err != nil {
		return err
	}

	if existingRole.Id == uuid.Nil {
		return errors.New("role not found")
	}

	if role.Name != nil {
		existingRole.Name = *role.Name
	}

	if role.Description != nil {
		existingRole.Description = role.Description
	}

	if role.Permissions != nil {
		existingRole.Permissions = role.Permissions
	}

	existingRole.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where("id = $1", id).
		Updates(&map[string]any{
			"name":                existingRole.Name,
			"description":         existingRole.Description,
			"permissions":         existingRole.Permissions,
			"modified_by_user_id": existingRole.ModifiedByUserId,
		}).Error; err != nil {
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
func (s *RolesService) Delete(auditId uuid.UUID, id uuid.UUID) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Where("id = $1", id).Delete(&models.Role{}).Error; err != nil {
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
func (s *RolesService) GetById(id uuid.UUID) (*models.Role, error) {
	var role models.Role

	if err := s.Storage.Postgres.
		Where("id = $1", id).
		Preload("ModifiedByUser").
		Find(&role).Error; err != nil {
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

// GetTotal returns the total number of Role records in the database that match the provided GORM clause expressions.
// It accepts a variable number of clause.Expression arguments to filter the query.
// Returns the count as int64 and an error if the query fails.
func (s *RolesService) GetTotal(clauses ...clause.Expression) (int64, error) {
	var total int64

	if err := s.Storage.Postgres.Model(&models.Role{}).Clauses(clauses...).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
