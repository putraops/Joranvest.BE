package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityOrganizationView struct {
	models.Organization
	Rating      float32 `json:"rating"`
	TotalRating int     `json:"total_rating"`
}

func (EntityOrganizationView) TableName() string {
	return "vw_organization"
}

func (EntityOrganizationView) ViewModel() string {
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
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.name,")
	sql.WriteString("  r.description,")
	sql.WriteString("  r.filepath,")
	sql.WriteString("  r.filepath_thumbnail,")
	sql.WriteString("  r.filename,")
	sql.WriteString("  r.extension,")
	sql.WriteString("  r.size,")
	sql.WriteString("  COALESCE(m.rating, 0) AS rating,")
	sql.WriteString("  COALESCE(m.total_rating, 0) AS total_rating,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString("FROM organization r ")
	sql.WriteString("LEFT JOIN LATERAL get_webinar_speaker_rating(r.id) m(rating, total_rating) ON true ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	return sql.String()
}
func (EntityOrganizationView) Migration() map[string]string {
	var view = EntityOrganizationView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
