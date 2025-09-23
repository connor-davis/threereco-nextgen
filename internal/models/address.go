package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/goccy/go-json"
)

type Address struct {
	LineOne  string  `json:"lineOne" gorm:"type:text;not null"`
	LineTwo  *string `json:"lineTwo" gorm:"type:text"`
	City     string  `json:"city" gorm:"type:text;not null"`
	Province string  `json:"province" gorm:"type:text;not null"`
	Country  string  `json:"country" gorm:"type:text;not null"`
	ZipCode  string  `json:"zipCode" gorm:"type:text;not null"`
}

// Value implements the driver.Valuer interface
func (a *Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface
func (a *Address) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %v to Address", value)
	}
	return json.Unmarshal(bytes, a)
}
