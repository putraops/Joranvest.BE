package models

import (
	"database/sql"
	"time"
)

type BaseModel struct {
	Id         string       `gorm:"type:nvarchar(50);primary_key:pk_id" json:"id"`
	IsActive   bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked   bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault  bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt  time.Time    `gorm:"type:time" json:"created_at"`
	CreatedBy  string       `gorm:"type:nvarchar(50)" json:"created_by"`
	UpdatedAt  time.Time    `gorm:"type:time" json:"updated_at"`
	UpdatedBy  string       `gorm:"type:nvarchar(50)" json:"updated_by"`
	ApprovedAt sql.NullTime `gorm:"type:datetime2(7);default:null" json:"approved_at"`
	ApprovedBy string       `gorm:"type:nvarchar(50)" json:"approved_by"`
}
