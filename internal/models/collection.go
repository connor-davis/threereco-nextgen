package models

import "github.com/google/uuid"

type Collection struct {
	Base
	SellerId  uuid.UUID            `json:"sellerId" gorm:"type:uuid;not null"`
	Seller    User                 `json:"seller" gorm:"foreignKey:SellerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BuyerId   uuid.UUID            `json:"buyerId" gorm:"type:uuid;not null"`
	Buyer     Business             `json:"buyer" gorm:"foreignKey:BuyerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Materials []CollectionMaterial `json:"materials" gorm:"many2many:collections_materials;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CollectionMaterial struct {
	Base
	Name         string  `json:"name" gorm:"type:text;not null"`
	GWCode       string  `json:"gwCode" gorm:"type:text;not null"`
	CarbonFactor string  `json:"carbonFactor" gorm:"type:text;not null"`
	Weight       float64 `json:"weight" gorm:"type:decimal(10,2);not null"`
	Value        float64 `json:"value" gorm:"type:decimal(10,2);not null"`
}
