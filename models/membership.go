package models

import (
	"database/sql"
)

type Membership struct {
	Id          string       `gorm:"type:varchar(50);primary_key:pk_id" json:"id"`
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

	Name        string  `gorm:"type:varchar(50);not null" json:"name"`
	Price       float64 `gorm:"type:decimal(18,2)" json:"price"`
	Duration    float64 `gorm:"type:decimal(18,2)" json:"duration"`
	TotalSaving float64 `gorm:"type:decimal(18,2)" json:"total_saving"`
	Description string  `gorm:"type:text" json:"description"`
}

func (Membership) TableName() string {
	return "membership"
}
