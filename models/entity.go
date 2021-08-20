package models

import (
	"database/sql"
)

type Entity struct {
	Id          string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive    bool         `gorm:"type:bool" json:"is_active"`
	IsLocked    bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault   bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt   sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy   string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt   sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy   string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt  sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy  string       `gorm:"type:varchar(50)" json:"approved_by"`
	Name        string       `gorm:"type:varchar(255)" json:"name"`
	Address     string       `gorm:"type:text" json:"address"`
	Phone       string       `gorm:"type:varchar(12)" json:"phone"`
	Email       string       `gorm:"type:varchar(255);unique" json:"email"`
	Description string       `gorm:"type:text" json:"description"`
}

func (Entity) TableName() string {
	return "entity"
}
