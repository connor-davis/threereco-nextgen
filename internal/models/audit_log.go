package models

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AuditLog represents a record of changes or actions performed on database tables.
// It captures metadata such as the table name, type of operation (e.g., insert, update, delete),
// the affected object ID, the data involved in the operation, and the user responsible for the action.
// The struct includes timestamps for creation and last update, and maintains a relation to the User entity.
//
// Fields:
//   - Id: Unique identifier for the audit log entry, automatically generated.
//   - TableName: Name of the table where the operation occurred.
//   - OperationType: Type of operation performed (e.g., insert, update, delete).
//   - ObjectId: Unique identifier of the object affected by the operation.
//   - Data: JSON representation of the data involved in the operation.
//   - UserId: Unique identifier of the user who performed the operation.
//   - User: Reference to the User who performed the operation, with cascading delete behavior.
//   - CreatedAt: Timestamp when the audit log entry was created, automatically set.
//   - UpdatedAt: Timestamp when the audit log entry was last updated, automatically set.
type AuditLog struct {
	Id             uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;"`
	TableName      string         `json:"tableName"`
	OperationType  string         `json:"operationType"`
	ObjectId       uuid.UUID      `json:"objectId"`
	Data           datatypes.JSON `json:"data"`
	OrganizationId uuid.UUID      `json:"organizationId" gorm:"type:uuid;"`
	UserId         uuid.UUID      `json:"userId"`
	User           *User          `json:"user" gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`
	CreatedAt      time.Time      `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt      time.Time      `json:"updatedAt" gorm:"autoUpdateTime;"`
}

// BeforeCreate is a GORM hook that runs before persisting an AuditLog.
// It ensures the AuditLog has a valid associated User (via UserId) and that
// the User has a PrimaryOrganizationId set. On success, it assigns the User's
// primary organization to AuditLog.OrganizationId. Any lookup or validation
// failure is logged and returned to abort the insert.
//
// Parameters:
//   - tx: the current GORM transaction used for database lookups.
//
// Returns:
//   - error: non-nil if the user cannot be found or lacks a primary organization.
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	var existingAuditUser *User

	if err := tx.Where(&User{Id: a.UserId}).Find(&existingAuditUser).Error; err != nil {
		log.Errorf("❌ Failed to find existing audit user: %s", err.Error())

		return err
	}

	if existingAuditUser == nil || existingAuditUser.Id == uuid.Nil {
		log.Errorf("❌ Failed to find existing audit user")

		return errors.New("failed to find existing audit user")
	}

	if existingAuditUser.PrimaryOrganizationId == nil {
		log.Errorf("❌ Failed to find existing audit user's primary organization")

		return errors.New("failed to find existing audit user's primary organization")
	}

	a.OrganizationId = *existingAuditUser.PrimaryOrganizationId

	return nil
}
