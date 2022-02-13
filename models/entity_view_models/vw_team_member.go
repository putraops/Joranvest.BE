package entity_view_models

import (
	"joranvest/models/view_models"
	"strings"
)

type EntityTeamMemberView struct {
	view_models.BaseViewModel
	TeamId                string `json:"name"`
	ApplicationUserId     string `json:"application_user_id"`
	TeamName              string `json:"team_name"`
	TeamMemberFirstName   string `json:"team_member_first_name"`
	TeamMemberLastName    string `json:"team_member_last_name"`
	TeamMemberFullname    string `json:"team_member_fullname"`
	TeamMemberInitialName string `json:"team_member_initial_name"`
}

func (EntityTeamMemberView) TableName() string {
	return "vw_team_member"
}

func (EntityTeamMemberView) ViewModel() string {
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
	sql.WriteString("  r.organization_id,")
	sql.WriteString("  o.name AS organization_name,")
	sql.WriteString("  t.id AS tenant_id,")
	sql.WriteString("  r.team_id,")
	sql.WriteString("  r.application_user_id,")
	sql.WriteString("  t.name AS team_name,")
	sql.WriteString("  m.first_name AS team_member_first_name,")
	sql.WriteString("  m.last_name AS team_member_last_name,")
	sql.WriteString("  CONCAT(m.first_name, ' ', m.last_name) AS team_member_fullname,")
	sql.WriteString("  CONCAT(UPPER(LEFT(m.first_name, 1)), '', UPPER(LEFT(m.last_name, 1))) AS team_member_initial_name,")
	sql.WriteString("  CASE WHEN u1.first_name IS NULL OR u1.first_name = '' THEN u1.username ELSE concat(u1.first_name, ' ', u1.last_name) END AS created_user_fullname,")
	sql.WriteString("  CASE WHEN u2.first_name IS NULL OR u2.first_name = '' THEN u2.username ELSE concat(u2.first_name, ' ', u2.last_name) END AS updated_user_fullname,")
	sql.WriteString("  CASE WHEN u3.first_name IS NULL OR u3.first_name = '' THEN u3.username ELSE concat(u3.first_name, ' ', u3.last_name) END AS submitted_user_fullname,")
	sql.WriteString("  CASE WHEN u4.first_name IS NULL OR u4.first_name = '' THEN u4.username ELSE concat(u4.first_name, ' ', u4.last_name) END AS approved_user_fullname ")
	sql.WriteString("FROM team_member r ")
	sql.WriteString("LEFT JOIN team t ON t.id = r.team_id ")
	sql.WriteString("LEFT JOIN organization o ON o.id = r.organization_id ")
	sql.WriteString("LEFT JOIN application_user m ON m.id = r.application_user_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	sql.WriteString("LEFT JOIN application_user u4 ON u4.id = r.approved_by ")
	return sql.String()
}
func (EntityTeamMemberView) Migration() map[string]string {
	var view = EntityTeamMemberView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
