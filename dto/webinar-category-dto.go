package dto

//WebinarCategoryDto is a model that client use when updating a book
type WebinarCategoryDto struct {
	Id                string `json:"id" form:"id"`
	Name              string `json:"name" form:"name" binding:"required"`
	Description       string `json:"description" form:"description"`
	WebibarCategoryId string `json:"webinar_category_id" form:"webinar_category_id"`
	EntityId          string `json:"-"`
	UpdatedBy         string
}
