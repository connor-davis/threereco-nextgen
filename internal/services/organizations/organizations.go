package organizations

import (
	"errors"
	"fmt"

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

func (s *OrganizationsService) SendInvite(auditId uuid.UUID, organizationId uuid.UUID, userId uuid.UUID) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Create(&models.Notification{
			Title:   "Organization Invite",
			Message: "You have been invited to join an organization.",
			Action: &models.NotificationAction{
				Link:     fmt.Sprintf("/organizations/%s/invite", organizationId),
				LinkText: "View Invitation",
			},
			UserId:           userId,
			ModifiedByUserId: auditId,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (s *OrganizationsService) AcceptInvite(auditId uuid.UUID, organizationId uuid.UUID, userId uuid.UUID) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Model(&models.Organization{
			Id:               organizationId,
			ModifiedByUserId: auditId,
		}).
		Association("Users").
		Append(&models.User{
			Id:               userId,
			ModifiedByUserId: auditId,
		}); err != nil {
		return err
	}

	return nil
}

// Create adds a new organization record to the database with the specified audit ID.
// It sets the audit user ID for tracking purposes before creating the organization.
// Returns an error if the operation fails.
func (s *OrganizationsService) Create(auditId uuid.UUID, organization models.CreateOrganizationPayload) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Create(&models.Organization{
			Name:             organization.Name,
			Domain:           organization.Domain,
			OwnerId:          organization.OwnerId,
			ModifiedByUserId: auditId,
		}).Error; err != nil {
		return err
	}

	return nil
}

// Update updates an existing organization record in the database identified by the given id.
// It uses the provided auditId for tracking the user performing the update.
// The organization parameter contains the new data to be updated.
// Returns an error if the update operation fails.
func (s *OrganizationsService) Update(auditId uuid.UUID, id uuid.UUID, organization models.UpdateOrganizationPayload) error {
	var existingOrganization models.Organization

	if err := s.Storage.Postgres.Where("id = $1", id).
		Find(&existingOrganization).Error; err != nil {
		return err
	}

	if existingOrganization.Id == uuid.Nil {
		return errors.New("organization not found")
	}

	if organization.Name != nil {
		existingOrganization.Name = *organization.Name
	}

	if organization.Domain != nil {
		existingOrganization.Domain = *organization.Domain
	}

	if organization.OwnerId != nil {
		existingOrganization.OwnerId = *organization.OwnerId
	}

	existingOrganization.ModifiedByUserId = auditId

	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where("id = $1", id).
		Updates(&map[string]any{
			"name":                existingOrganization.Name,
			"domain":              existingOrganization.Domain,
			"owner_id":            existingOrganization.OwnerId,
			"modified_by_user_id": existingOrganization.ModifiedByUserId,
		}).Error; err != nil {
		return err
	}

	return nil
}

// Delete removes an organization record from the database by its ID.
// It associates the deletion with the provided auditId for auditing purposes.
// Returns an error if the deletion fails.
func (s *OrganizationsService) Delete(auditId uuid.UUID, id uuid.UUID) error {
	if err := s.Storage.Postgres.Set("one:audit_user_id", auditId).
		Where("id = $1", id).
		Delete(&models.Organization{}).Error; err != nil {
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
func (s *OrganizationsService) GetById(id uuid.UUID) (*models.Organization, error) {
	var organization models.Organization

	if err := s.Storage.Postgres.
		Where("id = $1", id).
		Preload("ModifiedByUser").
		Find(&organization).Error; err != nil {
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

	organizationsClauses := []clause.Expression{
		clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "created_at",
					},
					Desc: true,
				},
			},
		},
	}

	organizationsClauses = append(organizationsClauses, clauses...)

	if err := s.Storage.Postgres.
		Clauses(organizationsClauses...).
		Find(&organizations).Error; err != nil {
		return nil, err
	}

	return organizations, nil
}

// GetTotal returns the total number of Organization records that match the provided GORM clause expressions.
// It accepts a variable number of clause.Expression arguments to filter the query.
// The function returns the count of matching records and any error encountered during the query execution.
func (s *OrganizationsService) GetTotal(clauses ...clause.Expression) (int64, error) {
	var count int64

	if err := s.Storage.Postgres.
		Model(&models.Organization{}).
		Clauses(clauses...).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
