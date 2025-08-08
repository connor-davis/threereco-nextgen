package models

import (
	"errors"
	"slices"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Role represents a user role within the system, defining a set of permissions and associations.
// It includes a unique identifier, name, description, a list of permissions, and relationships
// to users and organizations. The struct also tracks creation and update timestamps.
//
// Fields:
//   - Id: Unique identifier for the role, automatically generated.
//   - Name: Name of the role, must not be null.
//   - Description: Optional description of the role.
//   - Permissions: List of permissions associated with the role, stored as a string array.
//   - Users: List of users associated with the role, with cascading delete behavior.
//   - Organizations: List of organizations associated with the role, with cascading delete behavior.
//   - CreatedAt: Timestamp when the role was created, automatically set.
//   - UpdatedAt: Timestamp when the role was last updated, automatically set.
type Role struct {
	Id               uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"`
	Name             string         `json:"name" gorm:"type:text;not null;"`
	Description      *string        `json:"description" gorm:"type:text;"`
	Permissions      pq.StringArray `json:"permissions" gorm:"type:text[];default:array[]::text[];"`
	Users            []User         `json:"users" gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;"`
	Organizations    []Organization `json:"organizations" gorm:"many2many:organization_roles;constraint:OnDelete:CASCADE;"`
	ModifiedByUserId uuid.UUID      `json:"modifiedById" gorm:"type:uuid;"`
	ModifiedByUser   *User          `json:"modifiedBy"`
	CreatedAt        time.Time      `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt        time.Time      `json:"updatedAt" gorm:"autoUpdateTime;"`
}

// CreateRolePayload represents the input payload used to create a new role.
// Name is the required unique name of the role.
// Description is an optional human-readable description; if nil, no description was provided,
// allowing the distinction between an explicitly empty string and absence of a value.
// Permissions is the list of permission identifiers (e.g., action or resource codes) that
// will be assigned to the role; it should typically be non-empty to grant capabilities.
type CreateRolePayload struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Permissions []string `json:"permissions"`
}

// UpdateRolePayload represents the set of mutable fields that can be supplied to update an existing role.
// Fields set to nil (for pointer fields) indicate that the corresponding attribute should not be changed.
// Name: Optional new role name.
// Description: Optional new role description.
// Permissions: Full replacement list of permission identifiers; if provided, it overwrites the existing set.
type UpdateRolePayload struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Permissions []string `json:"permissions"`
}

// HasPermissions checks if the Role has at least one of the specified permissions.
// It returns true if any of the provided permissions are found in the Role's Permissions slice.
// Otherwise, it returns false.
//
// Example usage:
//
//	if role.HasPermissions("read", "write") {
//	    // Role has at least one of the permissions
//	}
func (r *Role) HasPermissions(permissions ...string) bool {
	for _, permission := range permissions {
		found := slices.Contains(r.Permissions, permission)

		if found {
			return true
		}
	}

	return false
}

// AfterCreate is a GORM hook that is triggered after a Role record is created in the database.
// It retrieves the audit user ID from the transaction context, marshals the Role object to JSON,
// and creates an audit log entry recording the creation event. If any step fails, it logs the error
// and returns it to GORM, which may abort the transaction.
func (r *Role) AfterCreate(tx *gorm.DB) error {
	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	roleJSON, err := json.Marshal(r)

	if err != nil {
		log.Errorf("❌ Failed to marshal role: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "roles",
		OperationType: "INSERT",
		ObjectId:      r.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          roleJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for role creation: %s", err.Error())

		return err
	}

	log.Infof("✅ Role %s created successfully with ID %s", r.Name, r.Id)

	return nil
}

// AfterUpdate is a GORM hook that is triggered after a Role record is updated.
// It retrieves the audit user ID from the transaction context, marshals the updated Role
// into JSON, and creates an audit log entry recording the update operation.
// If any step fails (retrieving the user ID, marshaling the Role, or creating the audit log),
// it logs the error and returns it.
func (r *Role) AfterUpdate(tx *gorm.DB) error {
	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	roleJSON, err := json.Marshal(r)

	if err != nil {
		log.Errorf("❌ Failed to marshal role: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "roles",
		OperationType: "UPDATE",
		ObjectId:      r.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          roleJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for role update: %s", err.Error())

		return err
	}

	log.Infof("✅ Role %s updated successfully with ID %s", r.Name, r.Id)

	return nil
}

// AfterDelete is a GORM hook that is triggered after a Role record is deleted from the database.
// It logs the deletion event by creating an audit log entry containing details about the deleted role,
// the operation type, and the user who performed the deletion. If the audit user ID cannot be retrieved
// or if any error occurs during marshalling or audit log creation, the function returns an error.
func (r *Role) AfterDelete(tx *gorm.DB) error {
	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	roleJSON, err := json.Marshal(r)

	if err != nil {
		log.Errorf("❌ Failed to marshal role: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "roles",
		OperationType: "DELETE",
		ObjectId:      r.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          roleJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for role deletion: %s", err.Error())

		return err
	}

	log.Infof("✅ Role %s deleted successfully with ID %s", r.Name, r.Id)

	return nil
}
