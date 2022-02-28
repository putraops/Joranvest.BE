package entity_view_models

import (
	"joranvest/models/view_models"
	"strings"
)

type EntityEducationView struct {
	view_models.BaseViewModel
	EducationCategoryId   string `json:"education_category_id"`
	Title                 string `json:"title"`
	Level                 string `json:"level"`
	Description           string `json:"description"`
	Filepath              string `json:"filepath"`
	FilepathThumbnail     string `json:"filepath_thumbnail"`
	Filename              string `json:"filename"`
	Extension             string `json:"extension"`
	Size                  string `json:"size"`
	EducationCategoryName string `json:"education_category_name"`
	ParentId              string `json:"parent_id"`
	ParentName            string `json:"parent_name"`
}

func (EntityEducationView) TableName() string {
	return "vw_education"
}

func (EntityEducationView) ViewModel() string {
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
	sql.WriteString("  r.submitted_at,")
	sql.WriteString("  r.submitted_by,")
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.education_category_id,")
	sql.WriteString("  c.name AS education_category_name,")
	sql.WriteString("  r.title,")
	sql.WriteString("  r.level,")
	sql.WriteString("  r.description,")
	sql.WriteString("  r.filepath,")
	sql.WriteString("  r.filepath_thumbnail,")
	sql.WriteString("  r.filename,")
	sql.WriteString("  r.extension,")
	sql.WriteString("  r.size,")
	sql.WriteString("  CASE WHEN u1.first_name IS NULL OR u1.first_name = '' THEN u1.username ELSE concat(u1.first_name, ' ', u1.last_name) END AS created_user_fullname,")
	sql.WriteString("  CASE WHEN u2.first_name IS NULL OR u2.first_name = '' THEN u2.username ELSE concat(u2.first_name, ' ', u2.last_name) END AS updated_user_fullname,")
	sql.WriteString("  CASE WHEN u3.first_name IS NULL OR u3.first_name = '' THEN u3.username ELSE concat(u3.first_name, ' ', u3.last_name) END AS submitted_user_fullname,")
	sql.WriteString("  CASE WHEN u4.first_name IS NULL OR u4.first_name = '' THEN u4.username ELSE concat(u4.first_name, ' ', u4.last_name) END AS approved_user_fullname ")
	sql.WriteString("FROM education r ")
	sql.WriteString("LEFT JOIN education_category c ON c.id = r.education_category_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	sql.WriteString("LEFT JOIN application_user u4 ON u4.id = r.approved_by ")
	return sql.String()
}
func (EntityEducationView) Migration() map[string]string {
	var view = EntityEducationView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
