package account

// Account represents a customer account.
type Account struct {
	ID             int64  `json:"account_id" example:"1"`
	DocumentNumber string `json:"document_number" example:"12345678900"`
}

// CreateAccountRequest is the payload to create an account.
type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" example:"12345678900"`
}
