package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityApplicationUserView struct {
	models.ApplicationUser
	RoleName    string `json:"role_name"`
	FullName    string `json:"full_name"`
	InitialName string `json:"initial_name"`
	HasRole     bool   `json:"has_role"`
	UserCreate  string `json:"user_create"`
	UserUpdate  string `json:"user_update"`
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
	sql.WriteString("  r.title,")
	sql.WriteString("  r.username,")
	sql.WriteString("  r.password,")
	sql.WriteString("  r.address,")
	sql.WriteString("  r.phone,")
	sql.WriteString("  r.email,")
	sql.WriteString("  r.firebase_id,")
	sql.WriteString("  r.total_point,")
	sql.WriteString("  r.is_email_verified,")
	sql.WriteString("  r.is_phone_verified,")
	sql.WriteString("  r.is_membership,")
	sql.WriteString("  false AS has_role,")
	sql.WriteString("  r.is_admin,")
	sql.WriteString("  CONCAT(r.first_name, ' ', r.last_name) AS full_name,")
	sql.WriteString("  CONCAT(UPPER(LEFT(r.first_name, 1)), '', UPPER(LEFT(r.last_name, 1))) AS initial_name,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString("FROM application_user r ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	return sql.String()
}
func (EntityApplicationUserView) Migration() map[string]string {
	var view = EntityApplicationUserView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
