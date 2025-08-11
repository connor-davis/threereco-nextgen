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
	Id            uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;"`
	TableName     string         `json:"tableName"`
	OperationType string         `json:"operationType"`
	ObjectId      uuid.UUID      `json:"objectId"`
	Data          datatypes.JSON `json:"data"`
	UserId        uuid.UUID      `json:"userId"`
	User          *User          `json:"user" gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"autoUpdateTime;"`
}

// AfterCreate is a GORM hook invoked after an AuditLog record is created.
// It is a no-op when the target table is "audit_logs" to avoid self-referential logging.
// The hook expects the acting user's ID to be present in the transaction context under
// the key "one:audit_user_id". It parses the ID as a UUID, verifies the user exists,
// and associates the new audit log with the user's primary organization via the
// "AuditLogs" association.
// Returns an error if the audit user ID is missing or invalid, the user cannot be found,
// the primary organization is unavailable, or the association append fails.
func (a *AuditLog) AfterCreate(tx *gorm.DB) error {
	if a.TableName == "audit_logs" {
		return nil
	}

	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")
	}

	auditUserIdUUID, err := uuid.Parse(auditUserId.(string))

	if err != nil {
		log.Errorf("❌ Failed to parse audit user ID: %s", err.Error())

		return err
	}

	var existingAuditUser *User

	if err := tx.Where(&User{Id: auditUserIdUUID}).Find(&existingAuditUser).Error; err != nil {
		log.Errorf("❌ Failed to find existing audit user: %s", err.Error())

		return err
	}

	if existingAuditUser == nil || existingAuditUser.Id == uuid.Nil {
		log.Errorf("❌ Failed to find existing audit user")

		return errors.New("failed to find existing audit user")
	}

	if err := tx.
		Model(&Organization{
			Id: *existingAuditUser.PrimaryOrganizationId,
		}).
		Association("AuditLogs").
		Append(a); err != nil {
		log.Errorf("❌ Failed to append audit log to organization: %s", err.Error())

		return err
	}

	return nil
}

// AfterUpdate is a GORM hook for AuditLog that runs after an update completes.
// It no-ops when the target table is "audit_logs" to avoid self-referential auditing.
// Otherwise, it:
//   - Reads the acting user's ID (key "one:audit_user_id") from the transaction.
//   - Parses the ID as a UUID and loads the corresponding User.
//   - Appends the audit log to that User's primary Organization via the "AuditLogs" association.
//
// Any failure (missing/invalid user ID, lookup errors, or association errors) is logged and returned;
// nil is returned on success.
func (a *AuditLog) AfterUpdate(tx *gorm.DB) error {
	if a.TableName == "audit_logs" {
		return nil
	}

	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")
	}

	auditUserIdUUID, err := uuid.Parse(auditUserId.(string))

	if err != nil {
		log.Errorf("❌ Failed to parse audit user ID: %s", err.Error())

		return err
	}

	var existingAuditUser *User

	if err := tx.Where(&User{Id: auditUserIdUUID}).Find(&existingAuditUser).Error; err != nil {
		log.Errorf("❌ Failed to find existing audit user: %s", err.Error())

		return err
	}

	if existingAuditUser == nil || existingAuditUser.Id == uuid.Nil {
		log.Errorf("❌ Failed to find existing audit user")

		return errors.New("failed to find existing audit user")
	}

	if err := tx.
		Model(&Organization{
			Id: *existingAuditUser.PrimaryOrganizationId,
		}).
		Association("AuditLogs").
		Append(a); err != nil {
		log.Errorf("❌ Failed to append audit log to organization: %s", err.Error())

		return err
	}

	return nil
}

// AfterDelete is a GORM hook executed after an AuditLog record is deleted.
// It is a no-op when the log references the "audit_logs" table to avoid self-referential updates.
// Otherwise, it:
// - Retrieves the acting user's ID from the transaction context (key "one:audit_user_id").
// - Parses the ID, loads, and validates the user.
// - Appends the current audit log to the user's primary organization via the "AuditLogs" association.
// The method logs and returns errors on failure (e.g., missing/invalid user ID, user not found,
// or association/DB errors); on success, it returns nil.
func (a *AuditLog) AfterDelete(tx *gorm.DB) error {
	if a.TableName == "audit_logs" {
		return nil
	}

	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")
	}

	auditUserIdUUID, err := uuid.Parse(auditUserId.(string))

	if err != nil {
		log.Errorf("❌ Failed to parse audit user ID: %s", err.Error())

		return err
	}

	var existingAuditUser *User

	if err := tx.Where(&User{Id: auditUserIdUUID}).Find(&existingAuditUser).Error; err != nil {
		log.Errorf("❌ Failed to find existing audit user: %s", err.Error())

		return err
	}

	if existingAuditUser == nil || existingAuditUser.Id == uuid.Nil {
		log.Errorf("❌ Failed to find existing audit user")

		return errors.New("failed to find existing audit user")
	}

	if err := tx.
		Model(&Organization{
			Id: *existingAuditUser.PrimaryOrganizationId,
		}).
		Association("AuditLogs").
		Append(a); err != nil {
		log.Errorf("❌ Failed to append audit log to organization: %s", err.Error())

		return err
	}

	return nil
}
