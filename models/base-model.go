package models

import (
	"time"
)

type BaseModel struct {
	Id          *string    `gorm:"type:varchar(50);primary_key:pk_id" json:"id"`
	IsActive    *bool      `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked    *bool      `gorm:"type:bool" json:"is_locked"`
	IsDefault   *bool      `gorm:"type:bool" json:"is_default"`
	CreatedAt   *time.Time `gorm:"type:timestamp" json:"created_at"`
	CreatedBy   *string    `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt   *time.Time `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy   *string    `gorm:"type:varchar(50)" json:"updated_by"`
	SubmittedAt *time.Time `gorm:"type:timestamp" json:"submitted_at"`
	SubmittedBy *string    `gorm:"type:varchar(50)" json:"submitted_by"`
	ApprovedAt  *time.Time `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy  *string    `gorm:"type:varchar(50)" json:"approved_by"`
	OwnerId     *string    `gorm:"type:varchar(50);null" json:"owner_id"`
	EntityId    *string    `gorm:"type:varchar(50);null" json:"entity_id"`
}
