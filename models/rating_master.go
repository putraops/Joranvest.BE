package models

import (
	"database/sql"
)

type RatingMaster struct {
	Id          string       `gorm:"type:varchar(50);primary_key:pk_id" json:"id"`
	IsActive    bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked    bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault   bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt   sql.NullTime `gorm:"type:timestamp;default:null" json:"created_at"`
	CreatedBy   string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt   sql.NullTime `gorm:"type:timestamp;default:null;" json:"updated_at"`
	UpdatedBy   string       `gorm:"type:varchar(50)" json:"updated_by"`
	SubmittedAt sql.NullTime `gorm:"type:timestamp" json:"submitted_at"`
	SubmittedBy string       `gorm:"type:varchar(50)" json:"submitted_by"`
	ApprovedAt  sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy  string       `gorm:"type:varchar(50)" json:"approved_by"`
	EntityId    string       `gorm:"type:varchar(50);null" json:"entity_id"`

	UserId        string `gorm:"type:varchar(50)" json:"user_id"`
	ObjectRatedId string `gorm:"type:varchar(50)" json:"object_rated_id"`
	ReferenceId   string `gorm:"type:varchar(50)" json:"reference_id"`
	Rating        int    `gorm:"type:int" json:"rating"`
	Comment       string `gorm:"type:text" json:"comment"`

	ApplicationUser ApplicationUser `gorm:"foreignkey:UserId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"application_user"`
}

func (RatingMaster) TableName() string {
	return "rating_master"
}
