package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityRoleMenuView struct {
	models.RoleMember
	RoleName            string `json:"role_name"`
	ApplicationMenuName string `json:"application_menu_name"`
	ParentId            string `json:"parent_id"`
	UserCreate          string `json:"user_create"`
	UserUpdate          string `json:"user_update"`
}

func (EntityRoleMenuView) TableName() string {
	return "vw_role_menu"
}

func (EntityRoleMenuView) ViewModel() string {
	var sql strings.Builder
	sql.WriteString("SELECT")
	sql.WriteString("  r.id,")
	sql.WriteString("  r.is_active,")
	sql.WriteString("  r.is_locked,")
	sql.WriteString("  r.is_default,")
	sql.WriteString("  r.created_at,")
	sql.WriteString("  r.created_by,")
	sql.WriteString("  r.updated_at,")
	sql.WriteString("  r.updated_by,")
	sql.WriteString("  r.approved_at,")
	sql.WriteString("  r.approved_by,")
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.role_id,")
	sql.WriteString("  ro.name AS role_name,")
	sql.WriteString("  r.application_menu_id,")
	sql.WriteString("  m.name AS application_menu_name,")
	sql.WriteString("  m.parent_id,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString("FROM role_menu r ")
	sql.WriteString("LEFT JOIN role ro ON ro.id = r.role_id ")
	sql.WriteString("LEFT JOIN application_menu m ON m.id = r.application_menu_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	return sql.String()
}
func (EntityRoleMenuView) Migration() map[string]string {
	var view = EntityRoleMenuView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
