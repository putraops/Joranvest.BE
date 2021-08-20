package models

import (
	"database/sql"
)

type RoleMember struct {
	Id         string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive   bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked   bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault  bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt  sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy  string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt  sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy  string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	OwnerName  string       `gorm:"type:varchar(50)" json:"owner_name"`
	ApprovedBy string       `gorm:"type:varchar(50)" json:"approved_by"`
	EntityId   string       `gorm:"type:varchar(50);null" json:"entity_id"`

	RoleId            string `gorm:"type:varchar(50);not null" json:"role_id"`
	ApplicationUserId string `gorm:"type:varchar(50);not null" json:"application_user_id"`

	Role            Role            `gorm:"foreignkey:RoleId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"role"`
	ApplicationUser ApplicationUser `gorm:"foreignkey:ApplicationUserId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"application_user"`
}

func (RoleMember) TableName() string {
	return "role_member"
}
