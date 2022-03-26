package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityPaymentView struct {
	models.Payment
	PaymentUserFullname   string   `json:"payment_user_fullname"`
	MembershipId          string   `json:"membership_id"`
	MembershipName        string   `json:"membership_name"`
	MembershipDuration    *float64 `json:"membership_duration"`
	WebinarId             string   `json:"webinar_id"`
	WebinarTitle          string   `json:"webinar_title"`
	WebinarRegistrationId string   `json:"webinar_registration_id"`
	ProductId             string   `json:"product_id"`
	ProductName           string   `json:"product_name"`
	ProductDuration       *float64 `json:"product_duration"`
	CreatedByFullname     string   `json:"created_by_fullname"`
	UserCreateTitle       string   `json:"user_create_title"`
	UpdatedByFullname     string   `json:"updated_by_fullname"`
	SubmittedFullname     string   `json:"submitted_by_fullname"`
}

func (EntityPaymentView) TableName() string {
	return "vw_payment"
}

func (EntityPaymentView) ViewModel() string {
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
	sql.WriteString("  r.record_id,")
	sql.WriteString("  r.is_extend_membership,")
	sql.WriteString("  r.application_user_id,")
	sql.WriteString("  CONCAT(u4.first_name, ' ', u4.last_name) AS payment_user_fullname,")
	sql.WriteString("  r.coupon_id,")
	sql.WriteString("  r.order_number,")
	sql.WriteString("  m.id AS membership_id,")
	sql.WriteString("  m.name AS membership_name,")
	sql.WriteString("  m.duration AS membership_duration,")
	sql.WriteString("  w.id AS webinar_id,")
	sql.WriteString("  w.title AS webinar_title,")
	sql.WriteString("  wr.id AS webinar_registration_id,")
	sql.WriteString("  p.id AS product_id,")
	sql.WriteString("  p.name AS product_name,")
	sql.WriteString("  p.duration AS product_duration,")
	sql.WriteString("  r.price,")
	sql.WriteString("  r.currency,")
	sql.WriteString("  r.payment_date,")
	sql.WriteString("  r.payment_date_expired,")
	sql.WriteString("  r.payment_type,")
	sql.WriteString("  r.payment_status,")
	sql.WriteString("  r.unique_number,")
	sql.WriteString("  r.account_name,")
	sql.WriteString("  r.account_number,")
	sql.WriteString("  r.bank_name,")
	sql.WriteString("  r.card_number,")
	sql.WriteString("  r.card_type,")
	sql.WriteString("  r.exp_month,")
	sql.WriteString("  r.exp_year,")
	sql.WriteString("  r.provider_name,")
	sql.WriteString("  r.provider_record_id,")
	sql.WriteString("  r.provider_reference_id,")
	sql.WriteString("  r.provider_business_id,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS created_by_fullname,")
	sql.WriteString("  u1.title AS user_create_title,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS updated_by_fullname,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS submitted_by_fullname ")
	sql.WriteString("FROM public.payment r ")
	sql.WriteString("LEFT JOIN membership m ON m.id = r.record_id ")
	sql.WriteString("LEFT JOIN webinar w ON w.id = r.record_id ")
	sql.WriteString("LEFT JOIN webinar_registration wr ON wr.webinar_id = r.record_id AND wr.application_user_id = r.created_by ")
	sql.WriteString("LEFT JOIN product p ON p.id = r.record_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	sql.WriteString("LEFT JOIN application_user u4 ON u4.id = r.application_user_id ")
	return sql.String()
}
func (EntityPaymentView) Migration() map[string]string {
	var view = EntityPaymentView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
