package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityMembershipView struct {
	models.Membership
}

func (EntityMembershipView) TableName() string {
	return "vw_membership"
}

func (EntityMembershipView) ViewModel() string {
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
	sql.WriteString("  r.name,")
	sql.WriteString("  r.price,")
	sql.WriteString("  r.description ")
	sql.WriteString("FROM membership r")
	return sql.String()
}
func (EntityMembershipView) Migration() map[string]string {
	var view = EntityMembershipView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
