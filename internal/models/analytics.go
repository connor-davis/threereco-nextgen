package models

type TopBusiness struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

type MaterialType struct {
	Type   string  `json:"type"`
	Weight float64 `json:"weight"`
}
