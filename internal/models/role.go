package models

import "github.com/lib/pq"

type Role struct {
	Base
	Name        string         `json:"name" gorm:"type:text;not null;uniqueIndex"`
	Description *string        `json:"description" gorm:"type:text"`
	Permissions pq.StringArray `json:"permissions" gorm:"type:text[];default:'{}'"`
	Default     bool           `json:"default" gorm:"type:boolean;default:false"`
}
