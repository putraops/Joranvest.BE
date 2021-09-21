package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityTechnicalAnalysisView struct {
	models.TechnicalAnalysis
	EmitenName        string `json:"emiten_name"`
	EmitenCode        string `json:"emiten_code"`
	CreatedByFullname string `json:"created_by_fullname"`
	UserCreateTitle   string `json:"user_create_title"`
	UpdatedByFullname string `json:"updated_by_fullname"`
	SubmittedFullname string `json:"submitted_by_fullname"`
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
	sql.WriteString("  r.submitted_at,")
	sql.WriteString("  r.submitted_by,")
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
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS created_by_fullname,")
	sql.WriteString("  u1.title AS user_create_title,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS updated_by_fullname,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS submitted_by_fullname ")
	sql.WriteString("FROM technical_analysis r ")
	sql.WriteString("LEFT JOIN emiten e ON e.id = r.emiten_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	return sql.String()
}
func (EntityTechnicalAnalysisView) Migration() map[string]string {
	var view = EntityTechnicalAnalysisView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
