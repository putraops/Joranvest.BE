package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityFilemasterView struct {
	models.Filemaster
}

func (EntityFilemasterView) TableName() string {
	return "vw_filemaster"
}

func (EntityFilemasterView) ViewModel() string {
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
	sql.WriteString("  r.record_id,")
	sql.WriteString("  r.filepath,")
	sql.WriteString("  r.filepath_thumbnail,")
	sql.WriteString("  r.filename,")
	sql.WriteString("  r.extension,")
	sql.WriteString("  r.size,")
	sql.WriteString("  r.description,")
	sql.WriteString("  r.file_type ")
	sql.WriteString("FROM filemaster r")
	return sql.String()
}
func (EntityFilemasterView) Migration() map[string]string {
	var view = EntityFilemasterView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
