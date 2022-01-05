package entity_view_models

import (
	"database/sql"
	"joranvest/models"
	"strings"
)

type EntityApplicationUserView struct {
	models.ApplicationUser
	RoleName           string       `json:"role_name"`
	FullName           string       `json:"full_name"`
	InitialName        string       `json:"initial_name"`
	HasRole            bool         `json:"has_role"`
	IsMembership       bool         `json:"is_membership"`
	MembershipId       string       `json:"membership_id"`
	MembershipName     string       `json:"membership_name"`
	MembershipDuration string       `json:"membership_duration"`
	MembershipDate     sql.NullTime `json:"membership_date"`
	MembershipExpired  sql.NullTime `json:"membership_expired"`
	Rating             float32      `json:"rating"`
	TotalRating        int          `json:"total_rating"`
	UserCreate         string       `json:"user_create"`
	UserUpdate         string       `json:"user_update"`
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
	sql.WriteString("  CASE")
	sql.WriteString("	   WHEN m.membership_id IS NOT NULL THEN TRUE")
	sql.WriteString("  	   ELSE FALSE")
	sql.WriteString("  END AS is_membership,")
	sql.WriteString("  m.membership_id,")
	sql.WriteString("  m.membership_name,")
	sql.WriteString("  m.membership_duration,")
	sql.WriteString("  m.membership_date,")
	sql.WriteString("  m.membership_expired,")
	sql.WriteString("  false AS has_role,")
	sql.WriteString("  r.is_admin,")
	sql.WriteString("  r.gender,")
	sql.WriteString("  r.filepath,")
	sql.WriteString("  r.filepath_thumbnail,")
	sql.WriteString("  r.filename,")
	sql.WriteString("  r.extension,")
	sql.WriteString("  r.size,")
	sql.WriteString("  r.description,")
	sql.WriteString("  COALESCE(w.rating, 0) AS rating,")
	sql.WriteString("  COALESCE(w.total_rating, 0) AS total_rating,")
	sql.WriteString("  CONCAT(r.first_name, ' ', r.last_name) AS full_name,")
	sql.WriteString("  CONCAT(UPPER(LEFT(r.first_name, 1)), '', UPPER(LEFT(r.last_name, 1))) AS initial_name,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString("FROM application_user r ")
	sql.WriteString("LEFT JOIN LATERAL get_membership_status(r.id) m ON true ")
	sql.WriteString("LEFT JOIN LATERAL get_webinar_speaker_rating(r.id) w(rating, total_rating) ON true ")
	//sql.WriteString("LEFT JOIN filemaster f ON f.record_id = r.id AND f.file_type = 1 ")
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
