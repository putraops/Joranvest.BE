package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityWebinarView struct {
	models.Webinar
	OrganizerOrganizationName string `json:"organizer_organization_name"`
	OrganizerFirstName        string `json:"organizer_first_name"`
	OrganizerLastName         string `json:"organizer_last_name"`
	OrganizerFullName         string `json:"organizer_full_name"`
	OrganizerInitialName      string `json:"organizer_initial_name"`
	UserCreate                string `json:"user_create"`
	UserUpdate                string `json:"user_update"`
}

func (EntityWebinarView) TableName() string {
	return "vw_webinar"
}

func (EntityWebinarView) ViewModel() string {
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
	sql.WriteString("  r.owner_id,")
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.webinar_category_id,")
	sql.WriteString("  r.organizer_organization_id,")
	sql.WriteString("  o.name AS organizer_organization_name,")
	sql.WriteString("  r.organizer_user_id,")
	sql.WriteString("  u3.first_name AS organizer_first_name,")
	sql.WriteString("  u3.last_name AS organizer_last_name,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS organizer_full_name,")
	sql.WriteString("  CONCAT(UPPER(LEFT(u3.first_name, 1)), '', UPPER(LEFT(u3.last_name, 1))) AS organizer_initial_name,")
	sql.WriteString("  r.title,")
	sql.WriteString("  r.description,")
	sql.WriteString("  r.webinar_first_start_date,")
	sql.WriteString("  r.webinar_first_end_date,")
	sql.WriteString("  r.webinar_last_start_date,")
	sql.WriteString("  r.webinar_last_end_date,")
	sql.WriteString("  r.min_age,")
	sql.WriteString("  r.webinar_level,")
	sql.WriteString("  r.price,")
	sql.WriteString("  r.discount,")
	sql.WriteString("  r.is_certificate,")
	sql.WriteString("  r.reward,")
	sql.WriteString("  r.status,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString("FROM webinar r ")
	sql.WriteString("  LEFT JOIN webinar_category c ON c.id = r.webinar_category_id")
	sql.WriteString("  LEFT JOIN organization o ON o.id = r.organizer_organization_id")
	sql.WriteString("  LEFT JOIN application_user u3 ON u3.id = r.organizer_user_id")
	sql.WriteString("  LEFT JOIN application_user u1 ON u1.id = r.created_by")
	sql.WriteString("  LEFT JOIN application_user u2 ON u2.id = r.updated_by")
	return sql.String()
}
func (EntityWebinarView) Migration() map[string]string {
	var view = EntityWebinarView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
