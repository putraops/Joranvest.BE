package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityApplicationUserView struct {
	models.ApplicationUser
	RoleName string `json:"role_name"`
}

func (EntityApplicationUserView) TableName() string {
	return "vw_application_user"
}

func (EntityApplicationUserView) ViewModel() string {
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
	sql.WriteString("  r.first_name,")
	sql.WriteString("  r.last_name,")
	sql.WriteString("  r.username,")
	sql.WriteString("  r.password,")
	sql.WriteString("  r.address,")
	sql.WriteString("  r.phone,")
	sql.WriteString("  r.email,")
	sql.WriteString("  r.is_admin,")
	sql.WriteString("  ro.name AS role_name ")
	sql.WriteString("FROM application_user r ")
	sql.WriteString("LEFT JOIN role_member rm ON rm.application_user_id = r.id ")
	sql.WriteString("LEFT JOIN role AS ro ON ro.id = rm.role_id ")
	return sql.String()
}
func (EntityApplicationUserView) Migration() map[string]string {
	var view = EntityMembershipView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
