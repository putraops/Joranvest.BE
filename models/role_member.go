package models

type RoleMember struct {
	BaseModel
	RoleId            string `gorm:"type:varchar(50);not null;uniqueIndex:idx_role_member" binding:"required" json:"role_id" `
	ApplicationUserId string `gorm:"type:varchar(50);not null;uniqueIndex:idx_role_member" binding:"required" json:"application_user_id"`

	Role            Role            `gorm:"foreignkey:RoleId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
	ApplicationUser ApplicationUser `gorm:"foreignkey:ApplicationUserId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
}

func (RoleMember) TableName() string {
	return "role_member"
}
