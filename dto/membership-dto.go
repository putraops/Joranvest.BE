package dto

//BookUpdateDTO is a model that client use when updating a book
type MembershipDto struct {
	Id          string  `json:"id" form:"id"`
	Name        string  `json:"name" form:"name" binding:"required"`
	Price       float64 `json:"price" form:"price"`
	Duration    float64 `json:"duration" form:"duration"`
	TotalSaving float64 `json:"total_saving" form:"total_saving"`
	Description string  `json:"description" form:"description"`
	EntityId    string  `json:"-"`
	UpdatedBy   string
}

type MembershipRecommendationDto struct {
	Id        string `json:"id" form:"id"`
	IsChecked bool   `json:"is_checked" form:"is_checked"`
}
