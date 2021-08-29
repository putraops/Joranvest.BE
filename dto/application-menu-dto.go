package dto

//ApplicationMenuDto is a model that client use when updating a book
type ApplicationMenuDto struct {
	Id                        string `json:"id" form:"id"`
	Name                      string `json:"name" form:"name"`
	OrderIndex                string `json:"order_index" form:"order_index"`
	ActionUrl                 string `json:"action_url" form:"action_url"`
	IconClass                 string `json:"icon_class" form:"icon_class"`
	ParentId                  string `json:"parent_id" form:"parent_id"`
	IsDivider                 string `json:"is_divider" form:"is_divider"`
	IsHeader                  string `json:"is_header" form:"is_header"`
	Description               string `json:"description" form:"description"`
	ApplicationMenuCategoryId string `json:"application_menu_category_id" form:"application_menu_category_id"`
	EntityId                  string `json:"-"`
	UpdatedBy                 string
}
