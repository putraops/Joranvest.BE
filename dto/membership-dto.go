package dto

//BookUpdateDTO is a model that client use when updating a book
type MembershipDto struct {
	Id          string  `json:"id" form:"id"`
	Name        string  `json:"name" form:"name" binding:"required"`
	Price       float64 `json:"price" form:"price"`
	Description string  `json:"description" form:"description"`
	EntityId    string  `json:"-"`
	UpdatedBy   string
}
