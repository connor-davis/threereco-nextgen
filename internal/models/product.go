package models

import (
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product represents an item offered by the system with a unique UUID identifier,
// human-readable name, and a monetary value (stored with two decimal precision).
// It may be associated with multiple organizations via a many-to-many relationship
// through the join table `organizations_products`. Timestamps for creation and
// last update are automatically managed by GORM.
type Product struct {
	Id               uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;"`
	Name             string     `json:"name" gorm:"type:text;not null;"`
	Value            float64    `json:"value" gorm:"type:decimal(10,2);not null;default:0.0;"`
	Materials        []Material `json:"materials" gorm:"many2many:products_materials;constraint:OnDelete:CASCADE;"`
	ModifiedByUserId uuid.UUID  `json:"modifiedById" gorm:"type:uuid;"`
	ModifiedByUser   *User      `json:"modifiedBy"`
	CreatedAt        time.Time  `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt        time.Time  `json:"updatedAt" gorm:"autoUpdateTime;"`
}

// CreateProductPayload represents the incoming data required to create a new product.
// Name is the human-readable product name.
// Value is the monetary value (price) of the product in the smallest currency unit or as a floating-point amount, depending on system conventions.
// Materials is an optional list of UUIDs representing materials associated with the product.
type CreateProductPayload struct {
	Name      string      `json:"name"`
	Value     float64     `json:"value"`
	Materials []uuid.UUID `json:"materials"` // Optional materials to associate with the product
}

// UpdateProductPayload represents the partial update input for a product.
// Pointer fields allow differentiating between omitted fields (nil) and
// explicit zero-value updates.
// - Name: Optional new product name.
// - Value: Optional new product value (e.g., price); nil means no change.
// - Materials: Optional new list of material UUIDs.
type UpdateProductPayload struct {
	Name      *string     `json:"name"`
	Value     *float64    `json:"value"`
	Materials []uuid.UUID `json:"materials"`
}

// AfterCreate is a GORM hook that is triggered before a Product record is created in the database.
// It retrieves the audit user ID from the transaction context, marshals the Product object to JSON,
// and creates an audit log entry recording the creation event. If any step fails, it logs the error
// and returns it to GORM, which may abort the transaction.
func (p *Product) AfterCreate(tx *gorm.DB) error {
	if _, ok := tx.Get("one:ignore_audit_log"); ok {
		return nil
	}

	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	productJSON, err := json.Marshal(p)

	if err != nil {
		log.Errorf("❌ Failed to marshal product: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "products",
		OperationType: "INSERT",
		ObjectId:      p.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          productJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for product creation: %s", err.Error())

		return err
	}

	log.Infof("✅ Product %s created successfully with ID %s", p.Name, p.Id)

	return nil
}

// AfterUpdate is a GORM hook that is triggered after a Product record is updated.
// It retrieves the audit user ID from the transaction context, marshals the updated Product
// into JSON, and creates an audit log entry recording the update operation.
func (p *Product) AfterUpdate(tx *gorm.DB) error {
	if _, ok := tx.Get("one:ignore_audit_log"); ok {
		return nil
	}

	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	productJSON, err := json.Marshal(p)

	if err != nil {
		log.Errorf("❌ Failed to marshal product: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "products",
		OperationType: "UPDATE",
		ObjectId:      p.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          productJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for product update: %s", err.Error())

		return err
	}

	log.Infof("✅ Product %s updated successfully with ID %s", p.Name, p.Id)

	return nil
}

// BeforeDelete is a GORM hook that is triggered after a Product record is deleted.
// It retrieves the audit user ID from the transaction context, marshals the Product object to JSON,
// and creates an audit log entry recording the deletion event. If any step fails, it logs the error
// and returns it to GORM, which may abort the transaction.
func (p *Product) BeforeDelete(tx *gorm.DB) error {
	if _, ok := tx.Get("one:ignore_audit_log"); ok {
		return nil
	}

	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	productJSON, err := json.Marshal(p)

	if err != nil {
		log.Errorf("❌ Failed to marshal product: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "products",
		OperationType: "DELETE",
		ObjectId:      p.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          productJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for product deletion: %s", err.Error())

		return err
	}

	log.Infof("✅ Product %s deleted successfully with ID %s", p.Name, p.Id)

	return nil
}
