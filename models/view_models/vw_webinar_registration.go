package entity_view_models

import (
	"database/sql"
	"joranvest/models"
	"strings"
)

type EntityWebinarRegistrationView struct {
	models.WebinarRegistration
	WebinarTitle        string       `json:"webinar_title"`
	WebinarDescription  string       `json:"webinar_description"`
	WebinarCategoryId   string       `json:"webinar_category_id"`
	WebinarCategoryName string       `json:"webinar_category_name"`
	WebinarStartDate    sql.NullTime `json:"webinar_start_date"`
	WebinarEndDate      sql.NullTime `json:"webinar_end_date"`
	MinAge              int          `json:"min_age"`
	WebinarLevel        string       `json:"webinar_level"`
	Price               float64      `json:"price"`
	Discount            float64      `json:"discount"`
	IsCertificate       bool         `json:"is_certificate"`
	Reward              int          `json:"reward"`
	PaymentDate         sql.NullTime `json:"payment_date"`
	PaymentType         string       `json:"payment_type"`
	PaymentStatus       int          `json:"payment_status"`
	SpeakerTitle        string       `json:"speaker_title"`

	UserFullname      string `json:"user_fullname"`
	CreatedByFullname string `json:"created_by_fullname"`
	UserCreateTitle   string `json:"user_create_title"`
	UpdatedByFullname string `json:"updated_by_fullname"`
	SubmittedFullname string `json:"submitted_by_fullname"`
}

func (EntityWebinarRegistrationView) TableName() string {
	return "vw_webinar_registration"
}

func (EntityWebinarRegistrationView) ViewModel() string {
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
	sql.WriteString("  r.submitted_at,")
	sql.WriteString("  r.submitted_by,")
	sql.WriteString("  r.approved_at,")
	sql.WriteString("  r.approved_by,")
	sql.WriteString("  r.owner_id,")
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.webinar_id,")
	sql.WriteString("  r.application_user_id,")
	sql.WriteString("  CONCAT(u4.first_name, ' ', u4.last_name) AS user_fullname,")
	sql.WriteString("  r.payment_date,")
	sql.WriteString("  r.payment_type,")
	sql.WriteString("  r.payment_status,")
	sql.WriteString("  w.webinar_category_id,")
	sql.WriteString("  c.name AS webinar_category_name,")
	sql.WriteString("  w.title AS webinar_title,")
	sql.WriteString("  w.description AS webinar_description,")
	sql.WriteString("  w.webinar_start_date,")
	sql.WriteString("  w.webinar_end_date,")
	sql.WriteString("  w.min_age,")
	sql.WriteString("  w.webinar_level,")
	sql.WriteString("  w.price,")
	sql.WriteString("  w.discount,")
	sql.WriteString("  w.is_certificate,")
	sql.WriteString("  w.reward,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS created_by_fullname,")
	sql.WriteString("  u1.title AS user_create_title,")
	sql.WriteString("  u4.title AS speaker_title,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS updated_by_fullname,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS submitted_by_fullname ")
	sql.WriteString("FROM webinar_registration r ")
	sql.WriteString("LEFT JOIN webinar w ON w.id = r.webinar_id ")
	sql.WriteString("LEFT JOIN webinar_category c ON c.id = w.webinar_category_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	sql.WriteString("LEFT JOIN application_user u4 ON u4.id = r.application_user_id ")
	return sql.String()
}
func (EntityWebinarRegistrationView) Migration() map[string]string {
	var view = EntityWebinarRegistrationView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
