package models

type Role struct {
	BaseModel
	Name               string `gorm:"type:varchar(50);unique" json:"name"`
	HasDashboardAccess *bool  `gorm:"type:bool;default:0" json:"has_dashboard_access"`
	HasFullAccess      *bool  `gorm:"type:bool;default:0" json:"has_full_access"`
	Description        string `gorm:"type:text" json:"description"`
}

func (Role) TableName() string {
	return "role"
}
