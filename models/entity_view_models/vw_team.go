package entity_view_models

import (
	"joranvest/models/view_models"
	"strings"
)

type EntityTeamView struct {
	view_models.BaseViewModel
	TeamName              string `json:"team_name"`
	TeamLeaderId          string `json:"team_leader_id"`
	TeamLeaderFirstName   string `json:"team_leader_first_name"`
	TeamLeaderLastName    string `json:"team_leader_last_name"`
	TeamLeaderFullname    string `json:"team_leader_fullname"`
	TeamLeaderInitialName string `json:"team_leader_initial_name"`
	Description           string `json:"description"`
}

func (EntityTeamView) TableName() string {
	return "vw_team"
}

func (EntityTeamView) ViewModel() string {
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
	sql.WriteString("  r.team_name,")
	sql.WriteString("  r.team_leader_id,")
	sql.WriteString("  tl.first_name AS team_leader_first_name,")
	sql.WriteString("  tl.last_name AS team_leader_last_name,")
	sql.WriteString("  CONCAT(tl.first_name, ' ', tl.last_name) AS team_leader_fullname,")
	sql.WriteString("  CONCAT(UPPER(LEFT(tl.first_name, 1)), '', UPPER(LEFT(tl.last_name, 1))) AS team_leader_initial_name,")
	sql.WriteString("  CASE WHEN u1.first_name IS NULL OR u1.first_name = '' THEN u1.username ELSE concat(u1.first_name, ' ', u1.last_name) END AS created_user_fullname,")
	sql.WriteString("  CASE WHEN u2.first_name IS NULL OR u2.first_name = '' THEN u2.username ELSE concat(u2.first_name, ' ', u2.last_name) END AS updated_user_fullname,")
	sql.WriteString("  CASE WHEN u3.first_name IS NULL OR u3.first_name = '' THEN u3.username ELSE concat(u3.first_name, ' ', u3.last_name) END AS submitted_user_fullname,")
	sql.WriteString("  CASE WHEN u4.first_name IS NULL OR u4.first_name = '' THEN u4.username ELSE concat(u4.first_name, ' ', u4.last_name) END AS approved_user_fullname ")
	sql.WriteString("FROM team r ")
	sql.WriteString("LEFT JOIN organization o ON o.id = r.organization_id ")
	sql.WriteString("LEFT JOIN application_user tl ON tl.id = r.team_leader_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	sql.WriteString("LEFT JOIN application_user u4 ON u4.id = r.approved_by ")
	return sql.String()
}
func (EntityTeamView) Migration() map[string]string {
	var view = EntityTeamView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
