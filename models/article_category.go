package models

import (
	"database/sql"
)

type ArticleCategory struct {
	Id          string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive    bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked    bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault   bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt   sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy   string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt   sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy   string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt  sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy  string       `gorm:"type:varchar(50)" json:"approved_by"`
	EntityId    string       `gorm:"type:varchar(50);null" json:"entity_id"`
	Name        string       `gorm:"type:varchar(50);unique" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
}

func (ArticleCategory) TableName() string {
	return "article_category"
}
