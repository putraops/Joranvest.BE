package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityRoleMemberView struct {
	models.RoleMember
	RoleName   string `json:"role_name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	FullName   string `json:"fullname"`
	UserInital string `json:"user_initial"`
	UserCreate string `json:"user_create"`
	UserUpdate string `json:"user_update"`
}

func (EntityRoleMemberView) TableName() string {
	return "vw_role_member"
}

func (EntityRoleMemberView) ViewModel() string {
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
	sql.WriteString("  r.application_user_id,")
	sql.WriteString("  a.first_name,")
	sql.WriteString("  a.last_name,")
	sql.WriteString("  CONCAT(a.first_name, ' ', a.last_name) AS fullname,")
	sql.WriteString("  CONCAT(LEFT(a.first_name, 1), '', LEFT(a.last_name, 1)) AS user_initial,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString("FROM role_member r ")
	sql.WriteString("LEFT JOIN role ro ON ro.id = r.role_id ")
	sql.WriteString("LEFT JOIN application_user a ON a.id = r.application_user_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	return sql.String()
}
func (EntityRoleMemberView) Migration() map[string]string {
	var view = EntityRoleMemberView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
