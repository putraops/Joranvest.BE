package dto

//-- RoleMenuDto is a model that client use when updating a book
type RoleMenuDto struct {
	Id                string `json:"id" form:"id"`
	RoleId            string `json:"role_id" form:"role_id" binding:"required"`
	ApplicationMenuId string `json:"application_menu_id" form:"application_menu_id" binding:"required"`
	EntityId          string `json:"-"`
	UpdatedBy         string
}

type InsertRoleMenuDto struct {
	Id                string `json:"id" form:"id"`
	RoleId            string `json:"role_id" form:"role_id" binding:"required"`
	ApplicationMenuId string `json:"application_menu_id" form:"application_menu_id" binding:"required"`
	IsParent          bool   `json:"is_parent" form:"is_parent"`
	HasChildren       bool   `json:"has_children" form:"has_children"`
	Children          string `json:"children" form:"children"`
	ChildrenIds       []string
	EntityId          string `json:"-"`
	CreatedBy         string
}

type DeleteRoleMenuDto struct {
	RoleId            string `json:"role_id" form:"role_id" binding:"required"`
	ApplicationMenuId string `json:"application_menu_id" form:"application_menu_id" binding:"required"`
	IsParent          bool   `json:"is_parent" form:"is_parent"`
}
