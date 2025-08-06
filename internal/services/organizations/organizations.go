package organizations

import (
	"github.com/connor-davis/threereco-nextgen/internal/models"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// OrganizationsService provides methods to manage organization-related operations.
// It interacts with the underlying storage layer to perform CRUD operations and other business logic for organizations.
type OrganizationsService struct {
	Storage *storage.Storage
}

// NewOrganizationsService creates and returns a new OrganizationsService instance
// using the provided storage. This service provides methods to manage organizations.
//
// Parameters:
//
//	storage - a pointer to the Storage instance used for data persistence.
//
// Returns:
//
//	A pointer to the newly created OrganizationsService.
func NewOrganizationsService(storage *storage.Storage) *OrganizationsService {
	return &OrganizationsService{
		Storage: storage,
	}
}

// Upsert creates a new organization record or updates an existing one in the database.
// It uses the provided auditId for tracking the user performing the operation.
// Additional GORM clause expressions can be passed to customize the query.
// Returns an error if the operation fails.
func (s *OrganizationsService) Upsert(auditId uuid.UUID, organization *models.Organization, clauses ...clause.Expression) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Clauses(clauses...).Assign(&organization).FirstOrCreate(&organization).Error; err != nil {
		return err
	}

	return nil
}

// Delete removes an organization record from the database by its ID.
// It associates the deletion with the provided auditId for auditing purposes.
// Returns an error if the deletion fails.
func (s *OrganizationsService) Delete(auditId uuid.UUID, id string) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).Where("id = ?", id).Delete(&models.Organization{}).Error; err != nil {
		return err
	}

	return nil
}

// GetById retrieves an organization by its unique identifier from the database.
// It returns a pointer to the Organization model and an error if the operation fails.
//
// Parameters:
//   - id: The unique identifier of the organization.
//
// Returns:
//   - *models.Organization: Pointer to the retrieved organization, or nil if not found.
//   - error: Error encountered during the retrieval, or nil if successful.
func (s *OrganizationsService) GetById(id string) (*models.Organization, error) {
	var organization models.Organization

	if err := s.Storage.Postgres.Where("id = ?", id).Find(&organization).Error; err != nil {
		return nil, err
	}

	return &organization, nil
}

// GetAll retrieves all organizations from the database, applying the provided GORM clause expressions as filters.
// It returns a slice of Organization models and an error if the query fails.
//
// Parameters:
//
//	clauses - Optional GORM clause expressions to customize the query.
//
// Returns:
//
//	[]models.Organization - A slice containing the retrieved organizations.
//	error                - An error if the database query fails, otherwise nil.
func (s *OrganizationsService) GetAll(clauses ...clause.Expression) ([]models.Organization, error) {
	var organizations []models.Organization

	if err := s.Storage.Postgres.Clauses(clauses...).Find(&organizations).Error; err != nil {
		return nil, err
	}

	return organizations, nil
}
