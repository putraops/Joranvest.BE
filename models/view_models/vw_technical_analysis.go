package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityTechnicalAnalysisView struct {
	models.TechnicalAnalysis
	EmitenName string `json:"emiten_name"`
	EmitenCode string `json:"emiten_code"`
	UserCreate string `json:"user_create"`
	UserUpdate string `json:"user_update"`
}

func (EntityTechnicalAnalysisView) TableName() string {
	return "vw_technical_analysis"
}

func (EntityTechnicalAnalysisView) ViewModel() string {
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
	sql.WriteString("  r.emiten_id,")
	sql.WriteString("  e.emiten_name,")
	sql.WriteString("  e.emiten_code,")
	sql.WriteString("  r.signal,")
	sql.WriteString("  r.bandarmology_status,")
	sql.WriteString("  r.start_buy,")
	sql.WriteString("  r.end_buy,")
	sql.WriteString("  r.start_sell,")
	sql.WriteString("  r.end_sell,")
	sql.WriteString("  r.start_cut,")
	sql.WriteString("  r.end_cut,")
	sql.WriteString("  r.start_ratio,")
	sql.WriteString("  r.end_ratio,")
	sql.WriteString("  r.timeframe,")
	sql.WriteString("  r.reason_to_buy,")
	sql.WriteString("  r.status,")
	sql.WriteString("  r.description,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString("FROM technical_analysis r ")
	sql.WriteString("LEFT JOIN emiten e ON e.id = r.emiten_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	return sql.String()
}
func (EntityTechnicalAnalysisView) Migration() map[string]string {
	var view = EntityTechnicalAnalysisView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
