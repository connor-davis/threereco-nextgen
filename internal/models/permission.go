package models

type Permission struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type PermissionGroup struct {
	Name        string            `json:"name"`
	Permissions []Permission      `json:"permissions"`
	SubGroups   []PermissionGroup `json:"subGroups"`
}
