package entity_view_models

import (
	"database/sql"
	"joranvest/models"
	"strings"
)

type EntityMembershipUserView struct {
	models.MembershipUser
	PaymenyDate        sql.NullTime `json:"payment_date"`
	PaymentType        string       `json:"payment_type"`
	CreatedByFullname  string       `json:"created_by_fullname"`
	UserCreateTitle    string       `json:"user_create_title"`
	UpdatedByFullname  string       `json:"updated_by_fullname"`
	SubmittedFullname  string       `json:"submitted_by_fullname"`
	MembershipFullname string       `json:"membership_user_fullname"`
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
	sql.WriteString("  CONCAT(u4.first_name, ' ', u4.last_name) AS membership_user_fullname,")
	sql.WriteString("  r.application_user_id,")
	sql.WriteString("  r.membership_payment_id,")
	sql.WriteString("  p.payment_date,")
	sql.WriteString("  p.payment_type,")
	sql.WriteString("  r.expired_date,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS created_by_fullname,")
	sql.WriteString("  u1.title AS user_create_title,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS updated_by_fullname,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS submitted_by_fullname ")
	sql.WriteString("FROM membership_user r ")
	sql.WriteString("LEFT JOIN membership m ON m.id = r.membership_id ")
	sql.WriteString("LEFT JOIN membership_payment p ON p.id = r.membership_payment_id ")
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
