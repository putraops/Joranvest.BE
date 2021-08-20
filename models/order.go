package models

import (
	"database/sql"
)

type Order struct {
	Id                  string       `gorm:"type:varchar(50);primary_key:pk_id" json:"id"`
	IsActive            bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked            bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault           bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt           sql.NullTime `gorm:"type:timestamp;default:null;" json:"created_at"`
	CreatedBy           string       `gorm:"type:varchar(50);default:null" json:"created_by"`
	UpdatedAt           sql.NullTime `gorm:"type:timestamp;default:null;" json:"updated_at"`
	UpdatedBy           string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt          sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy          string       `gorm:"type:varchar(50)" json:"approved_by"`
	OrderNumber         string       `gorm:"type:varchar(20);null" json:"order_number"`
	EntityId            string       `gorm:"type:varchar(50);null" json:"entity_id"`
	TenantId            string       `gorm:"type:varchar(50);null" json:"tenant_id"`
	EstimatedDate       sql.NullTime `gorm:"type:timestamp;default:null;" json:"estimated_date"`
	CustomerId          string       `gorm:"type:varchar(50);not null" json:"customer_id"`
	VoucherId           string       `gorm:"type:varchar(50);null" json:"voucher_id"`
	TotalPrice          float64      `gorm:"type:decimal(18,2)" json:"total_price"`
	TotalPayment        float64      `gorm:"type:decimal(18,2)" json:"total_payment"`
	InsufficientPayment float64      `gorm:"type:decimal(18,2)" json:"insufficient_payment"`
	TotalItem           int          `gorm:"type:int" json:"total_item"`
	Description         string       `gorm:"type:text" json:"description"`
	PaymentStatus       int          `gorm:"type:int" json:"payment_status"`
	OrderStatus         int          `gorm:"type:int" json:"order_status"`
	Entity              Entity       `gorm:"foreignkey:EntityId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"entity"`
}

func (Order) TableName() string {
	return "order"
}
