package dto

//ArticleCategoryDto is a model that client use when updating a book
type ArticleCategoryDto struct {
	Id          string `json:"id" form:"id"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
	ParentId    string `json:"parent_id" form:"parent_id"`
	EntityId    string `json:"-"`
	UpdatedBy   string
}
