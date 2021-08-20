package models

import (
	"database/sql"
)

type Emiten struct {
	Id         string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive   bool         `gorm:"type:bool" json:"is_active"`
	IsLocked   bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault  bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt  sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy  string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt  sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy  string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy string       `gorm:"type:varchar(50)" json:"approved_by"`
	EntityId   string       `gorm:"type:varchar(50);null" json:"entity_id"`

	EmitenName   string  `gorm:"type:varchar(255)" json:"emiten_name"`
	EmitenCode   string  `gorm:"type:varchar(10)" json:"emiten_code"`
	CurrentPrice float64 `gorm:"type:decimal(18,2)" json:"current_price"`
	Description  string  `gorm:"type:text" json:"description"`
}

func (Emiten) TableName() string {
	return "emiten"
}
