package dto

import (
	"time"
)

//PaymentDto is a model that client use when updating a book
type PaymentDto struct {
	Id            string    `json:"id" form:"id"`
	RecordId      string    `json:"record_id" form:"record_id"`
	OrderNumber   string    `json:"order_number" form:"order_number"`
	PaymentDate   time.Time `json:"payment_date" form:"payment_date"`
	PaymentType   string    `json:"payment_type" form:"payment_type"`
	PaymentStatus int       `json:"payment_status" form:"payment_status"`
	Price         int       `json:"price" form:"price"`
	UniqueNumber  int       `json:"unique_number" form:"unique_number"`
	AccountName   string    `json:"account_name" form:"account_name"`
	AccountNumber string    `json:"account_number" form:"account_number"`
	BankName      string    `json:"bank_name" form:"bank_name"`
	CardNumber    string    `json:"card_number" form:"card_number"`
	ExpMonth      int       `json:"exp_month" form:"exp_month"`
	ExpYear       int       `json:"exp_year" form:"exp_year"`
	EntityId      string    `json:"-"`
	UpdatedBy     string
}
