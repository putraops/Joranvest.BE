package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityEmailBlacklistView struct {
	models.EmailLogging
}

func (EntityEmailBlacklistView) TableName() string {
	return "vw_email_blacklist"
}

func (EntityEmailBlacklistView) ViewModel() string {
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
	sql.WriteString("  r.reason ")
	sql.WriteString("FROM email_blacklist r ")
	return sql.String()
}
func (EntityEmailBlacklistView) Migration() map[string]string {
	var view = EntityEmailBlacklistView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
