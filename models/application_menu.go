package models

import (
	"database/sql"
)

type ApplicationMenu struct {
	Id                        string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive                  bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked                  bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault                 bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt                 sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy                 string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt                 sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy                 string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt                sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy                string       `gorm:"type:varchar(50)" json:"approved_by"`
	EntityId                  string       `gorm:"type:varchar(50);null" json:"entity_id"`
	Name                      string       `gorm:"type:varchar(50)" json:"name"`
	OrderIndex                int          `gorm:"type:int" json:"order_index"`
	ActionUrl                 string       `gorm:"type:varchar(50)" json:"action_url"`
	IconClass                 string       `gorm:"type:varchar(50)" json:"icon_class"`
	ParentId                  string       `gorm:"type:varchar(50)" json:"parent_id"`
	Description               string       `gorm:"type:text" json:"description"`
	ApplicationMenuCategoryId string       `gorm:"type:varchar(50);not null" json:"application_menu_category_id"`

	ApplicationMenuCategory ApplicationMenuCategory `gorm:"foreignkey:ApplicationMenuCategoryId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"application_menu_category"`
}

func (ApplicationMenu) TableName() string {
	return "application_menu"
}
