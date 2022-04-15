package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityRoleMemberView struct {
	models.RoleMember
	ApplicationUserFirstName   string  `json:"application_user_first_name"`
	ApplicationUserLastName    string  `json:"application_user_last_name"`
	ApplicationUserFullname    string  `json:"application_user_fullname"`
	ApplicationUserInitialName string  `json:"application_user_initial_name"`
	RoleName                   string  `json:"role_name"`
	HasDashboardAccess         bool    `json:"has_dashboard_access"`
	HasFullAccess              bool    `json:"has_full_access"`
	CreatedUserFullname        *string `json:"created_by_fullname"`
	UpdatedUserFullname        *string `json:"updated_by_fullname"`
	SubmittedUserFullname      *string `json:"submitted_by_fullname"`
	ApprovedUserFullname       *string `json:"approved_by_fullname"`
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
	sql.WriteString("  r.submitted_at,")
	sql.WriteString("  r.submitted_by,")
	sql.WriteString("  r.owner_id,")
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.application_user_id,")
	sql.WriteString("  u5.last_name AS application_user_last_name,")
	sql.WriteString("  u5.first_name AS application_user_first_name,")
	sql.WriteString("  CASE WHEN u5.first_name IS NULL OR u5.first_name = '' THEN u1.username ELSE concat(u5.first_name, ' ', u5.last_name) END AS application_user_fullname,")
	sql.WriteString("  CONCAT(UPPER(LEFT(u5.first_name, 1)), '', UPPER(LEFT(u5.last_name, 1))) AS application_user_initial_name,")
	sql.WriteString("  r.role_id,")
	sql.WriteString("  rl.name AS role_name,")
	sql.WriteString("  rl.has_dashboard_access,")
	sql.WriteString("  rl.has_full_access,")
	sql.WriteString("  CASE WHEN u1.first_name IS NULL OR u1.first_name = '' THEN u1.username ELSE concat(u1.first_name, ' ', u1.last_name) END AS created_by_fullname,")
	sql.WriteString("  CASE WHEN u2.first_name IS NULL OR u2.first_name = '' THEN u2.username ELSE concat(u2.first_name, ' ', u2.last_name) END AS updated_by_fullname,")
	sql.WriteString("  CASE WHEN u3.first_name IS NULL OR u3.first_name = '' THEN u3.username ELSE concat(u3.first_name, ' ', u3.last_name) END AS submitted_by_fullname,")
	sql.WriteString("  CASE WHEN u4.first_name IS NULL OR u4.first_name = '' THEN u4.username ELSE concat(u4.first_name, ' ', u4.last_name) END AS approved_by_fullname ")
	sql.WriteString("FROM role_member r ")
	sql.WriteString("LEFT JOIN role rl ON rl.id = r.role_id ")
	sql.WriteString("LEFT JOIN application_user u5 ON u5.id = r.application_user_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	sql.WriteString("LEFT JOIN application_user u4 ON u4.id = r.approved_by ")
	return sql.String()
}
func (EntityRoleMemberView) Migration() map[string]string {
	var view = EntityRoleMemberView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
