package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityApplicationMenuView struct {
	models.ApplicationMenu
	ParentName                  string `json:"parent_name"`
	ApplicationMenuCategoryName string `json:"application_menu_category_name"`
	IsChecked                   bool   `json:"is_checked"`
}

func (EntityApplicationMenuView) TableName() string {
	return "vw_application_menu"
}

func (EntityApplicationMenuView) ViewModel() string {
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
	sql.WriteString("  r.name,")
	sql.WriteString("  r.order_index,")
	sql.WriteString("  r.action_url,")
	sql.WriteString("  r.icon_class,")
	sql.WriteString("  r.parent_id,")
	sql.WriteString("  a.name AS parent_name,")
	sql.WriteString("  COALESCE(r.is_divider, false) AS is_divider,")
	sql.WriteString("  COALESCE(r.is_header, false) AS is_header,")
	sql.WriteString("  r.description,")
	sql.WriteString("  r.application_menu_category_id,")
	sql.WriteString("  c.name AS application_menu_category_name,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString("FROM application_menu r ")
	sql.WriteString("LEFT JOIN application_menu_category c ON c.id = r.application_menu_category_id ")
	sql.WriteString("LEFT JOIN application_menu a ON a.id = r.parent_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	return sql.String()
}
func (EntityApplicationMenuView) Migration() map[string]string {
	var view = EntityApplicationMenuView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
