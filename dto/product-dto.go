package dto

//BookUpdateDTO is a model that client use when updating a book
type ProductDto struct {
	Id          string  `json:"id" form:"id"`
	Name        string  `json:"name" form:"name" binding:"required"`
	EntityId    string  `json:"-"`
	Price       float64 `json:"price" form:"price" binding:"required"`
	Description string  `json:"description" form:"description"`
	IsUnit      bool    `json:"is_unit" form:"is_unit"`
	UpdatedBy   string  `json:"-"`
	CategoryId  string  `json:"category_id" form:"category_id" binding:"required"`
}
