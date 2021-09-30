package models

import (
	"database/sql"
)

type MembershipPayment struct {
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

	PaymentDate   sql.NullTime `gorm:"type:timestamp" json:"payment_date"`
	PaymentType   string       `gorm:"type:varchar(50);not null" json:"payment_type"`
	PaymentStatus int          `gorm:"type:int;default:1" json:"payment_status"`
}

func (MembershipPayment) TableName() string {
	return "membership_payment"
}
