package models

import (
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization represents an organization within the application.
// It includes fields for the organization's unique identifier, creation and update timestamps,
// name, domain, owner, and relationships with users and roles.
// It also supports self-referencing for the user who last modified the organization record.
//
// Fields:
//   - Id: Unique identifier for the organization, automatically generated.
//   - Name: Name of the organization, must be unique and not null.
//   - Domain: Domain associated with the organization, must be unique and not null.
//   - OwnerId: Unique identifier of the user who owns the organization, must not be null.
//   - Owner: Reference to the User who owns the organization.
//   - Users: List of users associated with the organization, with cascading delete behavior.
//   - Roles: List of roles associated with the organization, with cascading delete behavior.
//   - ModifiedByUserId: UUID of the user who last modified this organization record.
//   - ModifiedByUser: Reference to the User who last modified this organization (self-referencing).
//   - CreatedAt: Timestamp when the organization was created, automatically set.
//   - UpdatedAt: Timestamp when the organization was last updated, automatically set.
type Organization struct {
	Id               uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;"`
	Name             string    `json:"name" gorm:"type:text;uniqueIndex;not null;"`
	Domain           string    `json:"domain" gorm:"type:text;uniqueIndex;not null;"`
	OwnerId          uuid.UUID `json:"ownerId" gorm:"type:uuid;not null;"`
	Owner            User      `json:"owner"`
	Users            []User    `json:"users" gorm:"many2many:organizations_users;constraint:OnDelete:CASCADE;"`
	Roles            []Role    `json:"roles" gorm:"many2many:organizations_roles;constraint:OnDelete:CASCADE;"`
	ModifiedByUserId uuid.UUID `json:"modifiedById" gorm:"type:uuid;"`
	ModifiedByUser   *User     `json:"modifiedBy"`
	CreatedAt        time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt        time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
}

// CreateOrganizationPayload contains the data required to create a new organization.
// Name is the human‑readable organization name.
// Domain is the unique domain or slug used to identify the organization (e.g. in URLs).
// OwnerId is the UUID of the user who will become the organization owner.
type CreateOrganizationPayload struct {
	Name    string    `json:"name"`
	Domain  string    `json:"domain"`
	OwnerId uuid.UUID `json:"ownerId"`
}

// UpdateOrganizationPayload contains optional fields for partially updating an organization.
// Any field left as nil will be ignored (the existing value is preserved).
// Name: new human‑readable name for the organization.
// Domain: primary domain to associate with the organization; should be unique if enforced by business rules.
// OwnerId: UUID of the user who will become the new owner (transfer of ownership).
type UpdateOrganizationPayload struct {
	Name    *string    `json:"name"`
	Domain  *string    `json:"domain"`
	OwnerId *uuid.UUID `json:"ownerId"`
}

// AfterCreate is a GORM hook that is triggered after a Organization record is created in the database.
// It retrieves the audit user ID from the transaction context, marshals the Organization object to JSON,
// and creates an audit log entry recording the creation event. If any step fails, it logs the error
// and returns it to GORM, which may abort the transaction.
func (o *Organization) AfterCreate(tx *gorm.DB) error {
	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	organizationJSON, err := json.Marshal(o)

	if err != nil {
		log.Errorf("❌ Failed to marshal organization: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "organizations",
		OperationType: "INSERT",
		ObjectId:      o.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          organizationJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for organization creation: %s", err.Error())

		return err
	}

	log.Infof("✅ Organization %s created successfully with ID %s", o.Name, o.Id)

	return nil
}

// AfterUpdate is a GORM hook that is triggered after a Organization record is updated.
// It retrieves the audit user ID from the transaction context, marshals the updated Organization
// into JSON, and creates an audit log entry recording the update operation.
// If any step fails (retrieving the user ID, marshaling the Organization, or creating the audit log),
// it logs the error and returns it.
func (o *Organization) AfterUpdate(tx *gorm.DB) error {
	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	organizationJSON, err := json.Marshal(o)

	if err != nil {
		log.Errorf("❌ Failed to marshal organization: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "organizations",
		OperationType: "UPDATE",
		ObjectId:      o.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          organizationJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for organization update: %s", err.Error())

		return err
	}

	log.Infof("✅ Organization %s updated successfully with ID %s", o.Name, o.Id)

	return nil
}

// AfterDelete is a GORM hook that is triggered after a Organization record is deleted from the database.
// It logs the deletion event by creating an audit log entry containing details about the deleted organization,
// the operation type, and the user who performed the deletion. If the audit user ID cannot be retrieved
// or if any error occurs during marshalling or audit log creation, the function returns an error.
func (o *Organization) AfterDelete(tx *gorm.DB) error {
	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	organizationJSON, err := json.Marshal(o)

	if err != nil {
		log.Errorf("❌ Failed to marshal organization: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "organizations",
		OperationType: "DELETE",
		ObjectId:      o.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          organizationJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for organization deletion: %s", err.Error())

		return err
	}

	log.Infof("✅ Organization %s deleted successfully with ID %s", o.Name, o.Id)

	return nil
}
