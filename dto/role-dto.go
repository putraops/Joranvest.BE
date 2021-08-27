package dto

//-- RoleDto is a model that client use when updating a book
type RoleDto struct {
	Id          string `json:"id" form:"id"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
	EntityId    string `json:"-"`
	UpdatedBy   string
}
