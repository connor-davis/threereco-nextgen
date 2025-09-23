package models

type Material struct {
	Base
	Name         string `json:"name" gorm:"type:text;not null;uniqueIndex"`
	GWCode       string `json:"gwCode" gorm:"type:text;not null;uniqueIndex"`
	CarbonFactor string `json:"carbonFactor" gorm:"type:text;not null"`
}
