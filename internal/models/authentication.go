package models

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterPayload struct {
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Type     UserType `json:"type"`
}

type PolicyType string

const (
	CollectionsPolicy  PolicyType = "collections"
	TransactionsPolicy PolicyType = "transactions"
	SystemPolicy       PolicyType = "system"
)

type VerifyMfaPayload struct {
	Code string `json:"code"`
}
