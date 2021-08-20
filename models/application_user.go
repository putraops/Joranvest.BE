package models

import (
	"database/sql"
)

type ApplicationUser struct {
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
	FirstName  string       `gorm:"type:varchar(50)" json:"first_name"`
	LastName   string       `gorm:"type:varchar(50)" json:"last_name"`
	Username   string       `gorm:"type:varchar(50);unique" json:"username"`
	Password   string       `gorm:"->;<-; not null" json:"-"`
	Address    string       `gorm:"type:text" json:"address"`
	Phone      string       `gorm:"type:varchar(12)" json:"phone"`
	Email      string       `gorm:"type:varchar(255);unique" json:"email"`
	IsAdmin    bool         `gorm:"type:bool;default:0" json:"is_admin"`
	Token      string       `gorm:"-" json:"token,omitempty"`
}

func (ApplicationUser) TableName() string {
	return "application_user"
}
