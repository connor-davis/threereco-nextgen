package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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
