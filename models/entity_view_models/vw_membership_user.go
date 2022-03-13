package entity_view_models

import (
	"joranvest/models"
	"strings"
	"time"
)

type EntityMembershipUserView struct {
	models.MembershipUser
	MembershipName         string     `json:"membership_name"`
	MembershipDuration     float64    `json:"membership_duration"`
	MembershipUserFullname string     `json:"membership_user_fullname"`
	MembershipPaymentDate  *time.Time `json:"membership_payment_date"`
	MembershipStartDate    *time.Time `json:"membership_start_date"`
	MembershipExpiredDate  *time.Time `json:"membership_expired_date"`
	IsExpired              bool       `json:"is_expired"`
	PaymentPrice           float64    `json:"payment_price"`
	PaymentUniqueNumber    float64    `json:"payment_unique_number"`
	PaymentType            string     `json:"payment_type"`
	CreatedByFullname      string     `json:"created_by_fullname"`
	UserCreateTitle        string     `json:"user_create_title"`
	UpdatedByFullname      string     `json:"updated_by_fullname"`
	SubmittedFullname      string     `json:"submitted_by_fullname"`
}

func (EntityMembershipUserView) TableName() string {
	return "vw_membership_user"
}

func (EntityMembershipUserView) ViewModel() string {
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
	sql.WriteString("  r.membership_id,")
	sql.WriteString("  m.name AS membership_name,")
	sql.WriteString("  m.duration AS membership_duration,")
	sql.WriteString("  CONCAT(u4.first_name, ' ', u4.last_name) AS membership_user_fullname,")
	sql.WriteString("  r.application_user_id,")
	sql.WriteString("  r.payment_id,")
	sql.WriteString("  p.payment_date AS membership_payment_date,")
	sql.WriteString("  p.payment_type,")
	sql.WriteString("  p.price AS payment_price,")
	sql.WriteString("  p.unique_number AS payment_unique_number,")
	sql.WriteString("  r.started_date,")
	sql.WriteString("  r.expired_date,")
	sql.WriteString("  r.started_date AS membership_started_date,")
	sql.WriteString("  r.expired_date AS membership_expired_date,")
	sql.WriteString("  CASE WHEN (r.started_date <= now() AND r.expired_date >= now()) THEN false ELSE TRUE END AS is_expired,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS created_by_fullname,")
	sql.WriteString("  u1.title AS user_create_title,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS updated_by_fullname,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS submitted_by_fullname ")
	sql.WriteString("FROM membership_user r ")
	sql.WriteString("LEFT JOIN membership m ON m.id = r.membership_id ")
	sql.WriteString("LEFT JOIN payment p ON p.id = r.payment_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	sql.WriteString("LEFT JOIN application_user u4 ON u4.id = r.application_user_id ")
	return sql.String()
}
func (EntityMembershipUserView) Migration() map[string]string {
	var view = EntityMembershipUserView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
