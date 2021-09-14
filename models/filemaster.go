package models

import (
	"database/sql"
)

type Filemaster struct {
	Id         string       `gorm:"type:varchar(50);primary_key:pk_id" json:"id"`
	IsActive   bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked   bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault  bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt  sql.NullTime `gorm:"type:timestamp;default:null" json:"created_at"`
	CreatedBy  string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt  sql.NullTime `gorm:"type:timestamp;default:null;" json:"updated_at"`
	UpdatedBy  string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy string       `gorm:"type:varchar(50)" json:"approved_by"`
	EntityId   string       `gorm:"type:varchar(50);null" json:"entity_id"`

	RecordId    string `gorm:"type:varchar(50);not null" json:"record_id"`
	Filepath    string `gorm:"type:varchar(200)" json:"filepath"`
	Filename    string `gorm:"type:varchar(200)" json:"filename"`
	FileType    int    `gorm:"type:int" json:"file_type"`
	Extension   string `gorm:"type:varchar(10)" json:"extension"`
	Size        string `gorm:"type:varchar(100)" json:"size"`
	Description string `gorm:"type:text" json:"description"`
}

func (Filemaster) TableName() string {
	return "filemaster"
}
