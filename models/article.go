package models

import (
	"database/sql"
)

type Article struct {
	Id         string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive   bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked   bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault  bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt  sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy  string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt  sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy  string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy string       `gorm:"type:varchar(50)" json:"approved_by"`
	EntityId   string       `gorm:"type:varchar(50);null" json:"entity_id"`

	Title       string `gorm:"type:text" json:"title"`
	SubTitle    string `gorm:"type:text" json:"sub_title"`
	Source      string `gorm:"type:varchar(200)" json:"name"`
	EditorId    string `gorm:"type:varchar(50);not null" json:"editor_id"`
	Description string `gorm:"type:text" json:"description"`

	Editor ApplicationUser `gorm:"foreignkey:EditorId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"editor"`
}

func (Article) TableName() string {
	return "article"
}
