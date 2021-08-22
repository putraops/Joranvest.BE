package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityWebinarCategoryView struct {
	models.WebinarCategory
	ParentName string `json:"parent_name"`
}

func (EntityWebinarCategoryView) TableName() string {
	return "vw_webinar_category"
}

func (EntityWebinarCategoryView) ViewModel() string {
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
	sql.WriteString("  r.parent_id,")
	sql.WriteString("  w.name AS parent_name,")
	sql.WriteString("  r.description,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString("FROM webinar_category r ")
	sql.WriteString("LEFT JOIN webinar_category w ON w.id = r.parent_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	return sql.String()
}
func (EntityWebinarCategoryView) Migration() map[string]string {
	var view = EntityWebinarCategoryView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
