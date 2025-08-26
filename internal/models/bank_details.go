package models

import (
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BankDetails represents the banking information associated with an account holder.
// Fields:
//
//	Id            - Unique UUID primary key (generated via uuid_generate_v4()).
//	AccountHolder - Full name of the person or entity that owns the bank account (required).
//	AccountNumber - Bank account number as provided by the financial institution (required).
//	BankName      - Name of the banking institution (required).
//	BranchCode    - Code identifying the specific branch of the bank (required).
//	CreatedAt     - Timestamp automatically set when the record is first created.
//	UpdatedAt     - Timestamp automatically updated whenever the record is modified.
type BankDetails struct {
	Id            uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	AccountHolder string    `json:"accountHolder" gorm:"type:text;not null"`
	AccountNumber string    `json:"accountNumber" gorm:"type:text;not null"`
	BankName      string    `json:"bankName" gorm:"type:text;not null"`
	BranchCode    string    `json:"branchCode" gorm:"type:text;not null"`
	UserId        uuid.UUID `json:"-" gorm:"type:uuid;not null"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// CreateBankDetailsPayload represents the request payload used to create or submit
// bank account details. It contains the information required to identify the
// account and its banking institution.
//
// Fields:
//   - AccountHolder: Full name of the person or entity that owns the bank account.
//   - AccountNumber: The bank account number as provided by the financial institution.
//   - BankName: The official name of the bank where the account is held.
//   - BranchCode: The branch or routing code identifying the specific bank branch.
type CreateBankDetailsPayload struct {
	AccountHolder string `json:"accountHolder"`
	AccountNumber string `json:"accountNumber"`
	BankName      string `json:"bankName"`
	BranchCode    string `json:"branchCode"`
}

// UpdateBankDetailsPayload represents a partial update request for a user's bank details.
// Any field set to a non-nil pointer will be considered for update; nil fields are ignored.
// Fields:
//   - AccountHolder: Optional updated name of the account holder.
//   - AccountNumber: Optional updated bank account number (string form to preserve leading zeros).
//   - BankName:      Optional updated name of the financial institution.
//   - BranchCode:    Optional updated branch/sort/routing code associated with the account.
type UpdateBankDetailsPayload struct {
	AccountHolder *string `json:"accountHolder"`
	AccountNumber *string `json:"accountNumber"`
	BankName      *string `json:"bankName"`
	BranchCode    *string `json:"branchCode"`
}

// AfterCreate is a GORM hook that runs immediately after a BankDetails record is
// inserted into the database. Unless the transaction contains the key
// "one:ignore_audit_log", it creates a corresponding AuditLog entry capturing
// the INSERT operation, the acting user (fetched from "one:audit_user_id"), and
// a JSON snapshot of the newly created BankDetails. It returns an error if the
// audit user ID is missing, if marshaling the BankDetails struct fails, or if
// persisting the AuditLog record fails. On success, it logs a confirmation
// message.
func (b *BankDetails) AfterCreate(tx *gorm.DB) error {
	if _, ok := tx.Get("one:ignore_audit_log"); ok {
		return nil
	}

	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	bankDetailsJSON, err := json.Marshal(b)

	if err != nil {
		log.Errorf("❌ Failed to marshal bank details: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "bank_details",
		OperationType: "INSERT",
		ObjectId:      b.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          bankDetailsJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for bank details creation: %s", err.Error())

		return err
	}

	log.Infof("✅ Bank details for account %s created successfully with ID %s", b.AccountHolder, b.Id)

	return nil
}

// AfterUpdate is a GORM hook that runs automatically after a BankDetails record
// has been successfully updated. It creates an audit log entry capturing the full
// state of the updated BankDetails row in JSON form.
//
// Behavior:
//   - If the GORM transaction context contains the key "one:ignore_audit_log",
//     the hook returns immediately and no audit entry is written.
//   - It expects the transaction context to contain a UUID under the key
//     "one:audit_user_id" identifying the user responsible for the change.
//     If this value is missing, the hook logs the issue and returns an error,
//     preventing silent audit omissions.
//   - On success, it persists an AuditLog record with:
//   - TableName: "bank_details"
//   - OperationType: "UPDATE"
//   - ObjectId: the BankDetails primary key
//   - UserId: the auditing user ID from the transaction
//   - Data: a JSON snapshot of the updated BankDetails struct
//   - Any failure to marshal the BankDetails struct or to insert the audit log
//     results in an error being returned.
//
// Logging:
//   - Logs descriptive errors for missing audit user ID, JSON marshal failures,
//     and audit log persistence failures.
//   - Logs a success message including the account holder and record ID when the
//     audit log is written successfully.
//
// Return:
//   - nil on success or when auditing is explicitly skipped.
//   - An error if required audit data is missing, JSON marshaling fails, or the
//     audit log cannot be created.
func (b *BankDetails) AfterUpdate(tx *gorm.DB) error {
	if _, ok := tx.Get("one:ignore_audit_log"); ok {
		return nil
	}

	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	bankDetailsJSON, err := json.Marshal(b)

	if err != nil {
		log.Errorf("❌ Failed to marshal bank details: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "bank_details",
		OperationType: "UPDATE",
		ObjectId:      b.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          bankDetailsJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for bank details update: %s", err.Error())

		return err
	}

	log.Infof("✅ Bank details for account %s updated successfully with ID %s", b.AccountHolder, b.Id)

	return nil
}

// BeforeDelete is a GORM lifecycle hook that runs just before a BankDetails record is deleted.
// It creates an audit log entry (unless the transaction context includes the key "one:ignore_audit_log").
// The hook:
//  1. Checks for a transaction flag "one:ignore_audit_log"; if present, auditing is skipped.
//  2. Retrieves the auditing user ID from the transaction context key "one:audit_user_id"; returns an error if absent.
//  3. Marshals the BankDetails instance into JSON to persist a snapshot of the deleted state.
//  4. Inserts an AuditLog row describing the DELETE operation (table, object ID, user ID, and serialized data).
//  5. Logs success or returns any encountered error (missing user ID, JSON marshalling failure, or audit log persistence failure).
//
// Return value:
//   - nil on success (including when auditing is intentionally ignored).
//   - A non-nil error if required context data is missing or persistence/marshalling fails.
//
// Expected transaction context keys:
//   - "one:ignore_audit_log" (optional, any value): when present, suppresses audit logging.
//   - "one:audit_user_id" (required unless auditing is ignored): must be a uuid.UUID identifying the acting user.
func (b *BankDetails) BeforeDelete(tx *gorm.DB) error {
	if _, ok := tx.Get("one:ignore_audit_log"); ok {
		return nil
	}

	auditUserId, ok := tx.Get("one:audit_user_id")

	if !ok {
		log.Errorf("❌ Failed to get audit user ID")

		return errors.New("failed to get audit user ID")
	}

	bankDetailsJSON, err := json.Marshal(b)

	if err != nil {
		log.Errorf("❌ Failed to marshal bank details: %s", err.Error())

		return err
	}

	auditLog := &AuditLog{
		TableName:     "bank_details",
		OperationType: "DELETE",
		ObjectId:      b.Id,
		UserId:        auditUserId.(uuid.UUID),
		Data:          bankDetailsJSON,
	}

	if err := tx.Create(auditLog).Error; err != nil {
		log.Errorf("❌ Failed to create audit log for bank details deletion: %s", err.Error())

		return err
	}

	log.Infof("✅ Bank details for account %s deleted successfully with ID %s", b.AccountHolder, b.Id)

	return nil
}
