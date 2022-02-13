package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityEmailLoggingView struct {
	models.EmailLogging
}

func (EntityEmailLoggingView) TableName() string {
	return "vw_email_logging"
}

func (EntityEmailLoggingView) ViewModel() string {
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
	sql.WriteString("  r.email,")
	sql.WriteString("  r.last_sent,")
	sql.WriteString("  r.mail_type ")
	sql.WriteString("FROM email_logging r ")
	return sql.String()
}
func (EntityEmailLoggingView) Migration() map[string]string {
	var view = EntityEmailLoggingView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
