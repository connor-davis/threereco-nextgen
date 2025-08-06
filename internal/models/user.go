package models

import (
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents an application user with authentication and profile information.
// It includes fields for unique identification, timestamps, contact details, permissions,
// roles, organizational associations, and auditing information for tracking modifications.
//
// Fields:
//   - Id: Unique identifier for the user, automatically generated.
//   - Email: User's email address, must be unique and not null.
//   - Image: Optional URL to the user's profile image.
//   - Name: Optional name of the user.
//   - Phone: Optional phone number of the user.
//   - Roles: List of roles associated with the user, with cascading delete behavior.
//   - Organizations: List of organizations the user is associated with, with cascading delete behavior.
//   - PrimaryOrganizationId: UUID of the user's primary organization, can be null.
//   - ModifiedByUserId: UUID of the user who last modified this user record.
//   - ModifiedByUser: Reference to the User who last modified this user (self-referencing).
//   - CreatedAt: Timestamp when the user was created, automatically set.
//   - UpdatedAt: Timestamp when the user was last updated, automatically set.
type User struct {
	Id                    uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;"`
	Email                 string         `json:"email" gorm:"type:text;uniqueIndex;not null;"`
	Password              []byte         `json:"-" gorm:"type:bytea;"`
	Image                 *string        `json:"image" gorm:"type:text;"`
	Name                  *string        `json:"name" gorm:"type:text;"`
	Phone                 *string        `json:"phone" gorm:"type:text;"`
	JobTitle              *string        `json:"jobTitle" gorm:"type:text;"`
	MfaSecret             []byte         `json:"-" gorm:"type:bytea;"`
	MfaEnabled            bool           `json:"mfaEnabled" gorm:"default:false;"`
	MfaVerified           bool           `json:"mfaVerified" gorm:"default:false;"`
	Roles                 []Role         `json:"roles" gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;"`
	Organizations         []Organization `json:"organizations" gorm:"many2many:organization_users;constraint:OnDelete:CASCADE;"`
	PrimaryOrganizationId uuid.UUID      `json:"primaryOrganizationId" gorm:"type:uuid;"`
	ModifiedByUserId      uuid.UUID      `json:"modifiedById" gorm:"type:uuid;"`
	ModifiedByUser        *User          `json:"modifiedBy"`
	CreatedAt             time.Time      `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt             time.Time      `json:"updatedAt" gorm:"autoUpdateTime;"`
}

// AfterCreate is a GORM hook that is triggered after a User record is created in the database.
// It retrieves the audit user ID from the transaction context, marshals the User object to JSON,
// and creates an audit log entry recording the creation event. If any step fails, it logs the error
// and returns it to GORM, which may abort the transaction.
func (u *User) AfterCreate(tx *gorm.DB) error {
	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	userJSON, err := json.Marshal(u)

	if err != nil {
		log.Errorf("❌ Failed to marshal user: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "users",
		OperationType: "INSERT",
		ObjectId:      u.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          userJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for user creation: %s", err.Error())

		return err
	}

	log.Infof("✅ User %s created successfully with ID %s", u.Email, u.Id)

	return nil
}

// AfterUpdate is a GORM hook that is triggered after a User record is updated.
// It retrieves the audit user ID from the transaction context, marshals the updated User
// into JSON, and creates an audit log entry recording the update operation.
func (u *User) AfterUpdate(tx *gorm.DB) error {
	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	userJSON, err := json.Marshal(u)

	if err != nil {
		log.Errorf("❌ Failed to marshal user: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "users",
		OperationType: "UPDATE",
		ObjectId:      u.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          userJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for user update: %s", err.Error())

		return err
	}

	log.Infof("✅ User %s updated successfully with ID %s", u.Email, u.Id)

	return nil
}

// AfterDelete is a GORM hook that is triggered after a User record is deleted.
// It retrieves the audit user ID from the transaction context, marshals the User object to JSON,
// and creates an audit log entry recording the deletion event. If any step fails, it logs the error
// and returns it to GORM, which may abort the transaction.
func (u *User) AfterDelete(tx *gorm.DB) error {
	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	userJSON, err := json.Marshal(u)

	if err != nil {
		log.Errorf("❌ Failed to marshal user: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "users",
		OperationType: "DELETE",
		ObjectId:      u.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          userJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for user deletion: %s", err.Error())

		return err
	}

	log.Infof("✅ User %s deleted successfully with ID %s", u.Email, u.Id)

	return nil
}
