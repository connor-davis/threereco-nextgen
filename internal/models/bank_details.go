package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/goccy/go-json"
)

type BankDetails struct {
	AccountHolder string `json:"accountHolder" gorm:"type:text;not null"`
	AccountNumber string `json:"accountNumber" gorm:"type:text;not null"`
	BankName      string `json:"bankName" gorm:"type:text;not null"`
	BranchCode    string `json:"branchCode" gorm:"type:text;not null"`
}

// Value implements the driver.Valuer interface
func (b *BankDetails) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan implements the sql.Scanner interface
func (b *BankDetails) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %v to BankDetails", value)
	}
	return json.Unmarshal(bytes, b)
}
