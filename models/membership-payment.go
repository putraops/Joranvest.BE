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
	EntityId    string       `gorm:"type:varchar(50);null" json:"entity_id"`

	PaymenyDate sql.NullTime `gorm:"type:timestamp" json:"payment_date"`
	PaymentType string       `gorm:"type:varchar(50);not null" json:"payment_type"`

	Membership      Membership      `gorm:"foreignkey:MembershipId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"membership"`
	ApplicationUser ApplicationUser `gorm:"foreignkey:ApplicationUserId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"application_user"`
}

func (MembershipPayment) TableName() string {
	return "membership_payment"
}