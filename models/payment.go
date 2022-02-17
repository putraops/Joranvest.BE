package models

import (
	"database/sql"
)

type Payment struct {
	Id          string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive    bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked    bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault   bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt   sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy   string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt   sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy   string       `gorm:"type:varchar(50)" json:"updated_by"`
	SubmittedAt sql.NullTime `gorm:"type:timestamp" json:"submitted_at"`
	SubmittedBy string       `gorm:"type:varchar(50)" json:"submitted_by"`
	ApprovedAt  sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy  string       `gorm:"type:varchar(50)" json:"approved_by"`
	OwnerId     string       `gorm:"type:varchar(50);null" json:"owner_id"`
	EntityId    string       `gorm:"type:varchar(50);null" json:"entity_id"`

	RecordId           string       `gorm:"type:varchar(50)" json:"record_id"`
	IsExtendMembership bool         `gorm:"type:bool;default:0" json:"is_extend_membership"`
	ApplicationUserId  string       `gorm:"type:varchar(50)" json:"application_user_id"`
	CouponId           string       `gorm:"type:varchar(50)" json:"coupon_id"`
	OrderNumber        string       `gorm:"type:varchar(50)" json:"order_number"`
	PaymentDate        sql.NullTime `gorm:"type:timestamp" json:"payment_date"`
	PaymentDateExpired sql.NullTime `gorm:"type:timestamp" json:"payment_date_expired"`
	PaymentType        string       `gorm:"type:varchar(50);not null" json:"payment_type"`
	PaymentStatus      int          `gorm:"type:int" json:"payment_status"`
	Price              int          `gorm:"type:int" json:"price"`
	Currency           string       `gorm:"type:varchar(20)" json:"currency"`
	UniqueNumber       int          `gorm:"type:int" json:"unique_number"`

	AccountName   string `gorm:"type:varchar(200)" json:"account_name"`
	AccountNumber string `gorm:"type:varchar(30)" json:"account_number"`
	BankName      string `gorm:"type:varchar(50)" json:"bank_name"`
	CardNumber    string `gorm:"type:varchar(30)" json:"card_number"`
	CardType      string `gorm:"type:varchar(30)" json:"card_type"`
	ExpMonth      int    `gorm:"type:int" json:"exp_month"`
	ExpYear       int    `gorm:"type:int" json:"exp_year"`

	ProviderName        string `gorm:"type:varchar(50)" json:"provider_name"`
	ProviderRecordId    string `gorm:"type:varchar(50)" json:"provider_record_id"`
	ProviderReferenceId string `gorm:"type:varchar(50)" json:"provider_reference_id"`
	ProviderBusinessId  string `gorm:"type:varchar(50)" json:"provider_business_id"`
}

func (Payment) TableName() string {
	return "payment"
}
