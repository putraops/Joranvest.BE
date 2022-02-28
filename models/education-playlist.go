package models

import (
	"database/sql"
)

type EducationPlaylist struct {
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
	OwnerId     string       `gorm:"type:varchar(50)" json:"owner_id"`
	EntityId    string       `gorm:"type:varchar(50);null" json:"entity_id"`

	EducationId string `gorm:"type:varchar(50);not null" json:"education_id"`
	Title       string `gorm:"type:text;not null" json:"title"`
	FileUrl     string `gorm:"type:varchar(200)" json:"file_url"`
	Description string `gorm:"type:text" json:"description"`
	OrderIndex  int    `gorm:"type:int" json:"order_index"`

	Education Education `gorm:"foreignkey:EducationId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"education"`
}

func (EducationPlaylist) TableName() string {
	return "education_playlist"
}
