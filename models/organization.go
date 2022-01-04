package models

import (
	"database/sql"
)

type Organization struct {
	Id                string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive          bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked          bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault         bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt         sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy         string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt         sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy         string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt        sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy        string       `gorm:"type:varchar(50)" json:"approved_by"`
	SubmittedAt       sql.NullTime `gorm:"type:timestamp" json:"submitted_at"`
	SubmittedBy       string       `gorm:"type:varchar(50)" json:"submitted_by"`
	EntityId          string       `gorm:"type:varchar(50);null" json:"entity_id"`
	Name              string       `gorm:"type:varchar(50);unique" json:"name"`
	Description       string       `gorm:"type:text" json:"description"`
	Filepath          string       `gorm:"type:varchar(200)" json:"filepath"`
	FilepathThumbnail string       `gorm:"type:varchar(200)" json:"filepath_thumbnail"`
	Filename          string       `gorm:"type:varchar(200)" json:"filename"`
	Extension         string       `gorm:"type:varchar(10)" json:"extension"`
	Size              string       `gorm:"type:varchar(100)" json:"size"`
}

func (Organization) TableName() string {
	return "organization"
}
