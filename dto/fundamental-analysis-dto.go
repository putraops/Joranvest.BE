package dto

//FundamentalAnalysisDto is a model that client use when updating a book
type FundamentalAnalysisDto struct {
	Id             string  `json:"id" form:"id"`
	EmitenId       string  `json:"emiten_id" form:"emiten_id" binding:"required"`
	CurrentPrice   float64 `json:"current_price" form:"current_price"`
	NormalPrice    float64 `json:"normal_price" form:"normal_price"`
	MarginOfSafety float64 `json:"margin_of_safety" form:"margin_of_safety"`
	ResearchData   string  `json:"research_data" form:"research_data"`
	Status         string  `json:"status" form:"status"`
	Description    string  `json:"description" form:"description"`
	Tag            string  `json:"tag" form:"tag"`

	EntityId  string `json:"-"`
	UpdatedBy string
}
