package transaction

import "time"

// Transaction represents a financial movement on an account.
type Transaction struct {
	ID              int64     `json:"transaction_id" example:"1"`
	AccountID       int64     `json:"account_id" example:"1"`
	OperationTypeID int       `json:"operation_type_id" example:"1"`
	Amount          float64   `json:"amount" example:"-100"`
	EventDate       time.Time `json:"event_date" example:"2026-05-18T12:00:00Z"`
}

// CreateTransactionRequest is the payload to create a transaction.
// Operation types: 1=purchase, 2=installment, 3=withdrawal, 4=credit voucher.
type CreateTransactionRequest struct {
	AccountID       int64   `json:"account_id" example:"1"`
	OperationTypeID int     `json:"operation_type_id" example:"1"`
	Amount          float64 `json:"amount" example:"100"`
}
