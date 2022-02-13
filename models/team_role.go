package models

import (
	"database/sql"
)

type TeamRole struct {
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

	RoleId string `gorm:"type:varchar(50);not null" json:"role_id"`
	TeamId string `gorm:"type:varchar(50);not null" json:"team_id"`

	Role Role `gorm:"foreignkey:RoleId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"role"`
	Team Team `gorm:"foreignkey:TeamId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"team"`
}

func (TeamRole) TableName() string {
	return "team_role"
}
