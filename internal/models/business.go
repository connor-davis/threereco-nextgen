package models

import "github.com/google/uuid"

type Business struct {
	Base
	Name        string       `json:"name" gorm:"not null"`
	Address     *Address     `json:"address" gorm:"type:jsonb;"`
	BankDetails *BankDetails `json:"bankDetails" gorm:"type:jsonb;"`
	Users       []User       `json:"users" gorm:"many2many:businesses_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OwnerId     uuid.UUID    `json:"ownerId" gorm:"type:uuid;not null"`
}
