package models

import (
	"database/sql"
)

type ApplicationUser struct {
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
	EntityId          string       `gorm:"type:varchar(50);null" json:"entity_id"`
	FirstName         string       `gorm:"type:varchar(50)" json:"first_name"`
	LastName          string       `gorm:"type:varchar(50)" json:"last_name"`
	Title             string       `gorm:"type:varchar(200)" json:"title"`
	Username          string       `gorm:"type:varchar(50)" json:"username"`
	Password          string       `gorm:"->;<-; not null" json:"-"`
	Address           string       `gorm:"type:text" json:"address"`
	Phone             string       `gorm:"type:varchar(12)" json:"phone"`
	Email             string       `gorm:"type:varchar(255);unique" json:"email"`
	FirebaseId        string       `gorm:"type:varchar(50)" json:"firebase_id"`
	TotalPoint        string       `gorm:"type:int;default:0" json:"total_point"`
	IsEmailVerified   bool         `gorm:"type:bool;default:0" json:"is_email_verified"`
	IsPhoneVerified   bool         `gorm:"type:bool;default:0" json:"is_phone_verified"`
	IsAdmin           bool         `gorm:"type:bool;default:0" json:"is_admin"`
	Gender            bool         `gorm:"type:bool" json:"gender"`
	Filepath          string       `gorm:"type:varchar(200)" json:"filepath"`
	FilepathThumbnail string       `gorm:"type:varchar(200)" json:"filepath_thumbnail"`
	Filename          string       `gorm:"type:varchar(200)" json:"filename"`
	Extension         string       `gorm:"type:varchar(10)" json:"extension"`
	Size              string       `gorm:"type:varchar(100)" json:"size"`
	Description       string       `gorm:"type:text" json:"description"`
	Token             string       `gorm:"-" json:"token,omitempty"`
}

func (ApplicationUser) TableName() string {
	return "application_user"
}
