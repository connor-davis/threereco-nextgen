package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeCollection TransactionType = "collection"
	TransactionTypeTransfer   TransactionType = "transfer"
)

type Transaction struct {
	Id             uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"`
	Type           TransactionType `json:"type" gorm:"type:transaction_type;"`
	Weight         float64         `json:"weight" gorm:"type:decimal(10,2);not null;default:0.0;"`
	Amount         float64         `json:"amount" gorm:"type:decimal(10,2);not null;default:0.0;"`
	SellerAccepted bool            `json:"sellerAccepted" gorm:"type:boolean;not null;default:false;"`
	SellerDeclined bool            `json:"sellerDeclined" gorm:"type:boolean;not null;default:false;"`
	Products       []Product       `json:"products" gorm:"many2many:transactions_products;constraint:OnDelete:CASCADE;"`
	CreatedAt      time.Time       `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt      time.Time       `json:"updatedAt" gorm:"autoUpdateTime;"`
}
