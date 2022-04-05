package models

type Role struct {
	BaseModel
	Name        string `gorm:"type:varchar(50);unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

func (Role) TableName() string {
	return "role"
}
