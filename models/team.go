package models

import (
	"database/sql"
)

type Team struct {
	Id           string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive     bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked     bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault    bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt    sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy    string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt    sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy    string       `gorm:"type:varchar(50)" json:"updated_by"`
	SubmittedAt  sql.NullTime `gorm:"type:timestamp" json:"submitted_at"`
	SubmittedBy  string       `gorm:"type:varchar(50)" json:"submitted_by"`
	ApprovedAt   sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy   string       `gorm:"type:varchar(50)" json:"approved_by"`
	TeamName     string       `gorm:"type:varchar(255)" json:"team_name"`
	TeamLeaderId string       `gorm:"type:varchar(50)" json:"team_leader_id"`
	Description  string       `gorm:"type:varchar(255)" json:"description"`
}

func (Team) TableName() string {
	return "team"
}
