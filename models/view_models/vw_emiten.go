package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityEmitenView struct {
	models.Emiten
	SectorName                string `json:"sector_name"`
	SectorDescription         string `json:"sector_description"`
	EmitenCategoryName        string `json:"emiten_category_name"`
	EmitenCategoryDescription string `json:"emiten_category_description"`
}

func (EntityEmitenView) TableName() string {
	return "vw_emiten"
}

func (EntityEmitenView) ViewModel() string {
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
	sql.WriteString("  r.emiten_name,")
	sql.WriteString("  r.emiten_code,")
	sql.WriteString("  r.current_price,")
	sql.WriteString("  r.description, ")
	sql.WriteString("  r.sector_id,")
	sql.WriteString("  s.name AS sector_name,")
	sql.WriteString("  s.description AS sector_description,")
	sql.WriteString("  r.emiten_category_id,")
	sql.WriteString("  c.name AS emiten_category_name,")
	sql.WriteString("  c.description AS emiten_category_description,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString("FROM emiten r ")
	sql.WriteString("LEFT JOIN sector s ON s.id = r.sector_id ")
	sql.WriteString("LEFT JOIN emiten_category c ON c.id = r.emiten_category_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	return sql.String()
}
func (EntityEmitenView) Migration() map[string]string {
	var view = EntityEmitenView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
