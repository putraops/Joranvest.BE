package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityMembershipPaymentView struct {
	models.FundamentalAnalysis
	CreatedByFullname string `json:"created_by_fullname"`
	UserCreateTitle   string `json:"user_create_title"`
	UpdatedByFullname string `json:"updated_by_fullname"`
	SubmittedFullname string `json:"submitted_by_fullname"`
}

func (EntityMembershipPaymentView) TableName() string {
	return "vw_membership_payment"
}

func (EntityMembershipPaymentView) ViewModel() string {
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
	sql.WriteString("  r.payment_date,")
	sql.WriteString("  r.payment_type,")
	sql.WriteString("  r.payment_status,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS created_by_fullname,")
	sql.WriteString("  u1.title AS user_create_title,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS updated_by_fullname,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS submitted_by_fullname ")
	sql.WriteString("FROM membership_payment r ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	return sql.String()
}
func (EntityMembershipPaymentView) Migration() map[string]string {
	var view = EntityMembershipPaymentView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
