package models

type Product struct {
	BaseModel
	Name        *string `gorm:"type:text" json:"name"`
	Price       float64 `gorm:"type:decimal(18,2)" json:"price"`
	Duration    *int    `gorm:"type:int" json:"duration"`
	ProductType *string `gorm:"type:text" json:"product_type"`
	Description *string `gorm:"type:text" json:"description"`
}

func (Product) TableName() string {
	return "product"
}
