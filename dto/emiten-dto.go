package dto

//EmitenDto is a model that client use when updating a book
type EmitenDto struct {
	Id           string  `json:"id" form:"id"`
	EmitenName   string  `json:"emiten_name" form:"emiten_name" binding:"required"`
	EmitenCode   string  `json:"emiten_code" form:"emiten_code" binding:"required"`
	CurrentPrice float64 `json:"price" form:"current_price"`
	Description  string  `json:"description" form:"description"`
	EntityId     string  `json:"-"`
	UpdatedBy    string
}
