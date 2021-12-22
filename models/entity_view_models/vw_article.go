package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityArticleView struct {
	models.Article
	ArticleCategoryName string `json:"article_category_name"`
	Filepath            string `json:"filepath"`
	FilepathThumbnail   string `json:"filepath_thumbnail"`
	Filename            string `json:"filename"`
	Extension           string `json:"extension"`
	CreatedByFullname   string `json:"created_by_fullname"`
	UpdatedByFullname   string `json:"updated_by_fullname"`
	SubmittedFullname   string `json:"submitted_by_fullname"`
}

func (EntityArticleView) TableName() string {
	return "vw_article"
}

func (EntityArticleView) ViewModel() string {
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
	sql.WriteString("  r.title,")
	sql.WriteString("  r.sub_title,")
	sql.WriteString("  r.body,")
	sql.WriteString("  r.source,")
	sql.WriteString("  r.article_type,")
	sql.WriteString("  r.article_category_id,")
	sql.WriteString("  a.name As article_category_name,")
	sql.WriteString("  r.description,")
	sql.WriteString("  f.filepath,")
	sql.WriteString("  f.filepath_thumbnail,")
	sql.WriteString("  f.filename,")
	sql.WriteString("  f.extension,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS created_by_fullname,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS updated_by_fullname,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS submitted_by_fullname ")
	sql.WriteString("FROM public.article r ")
	sql.WriteString("LEFT JOIN filemaster f ON f.record_id = r.id AND f.file_type = 1 ")
	sql.WriteString("LEFT JOIN article_category a ON a.id = r.article_category_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	return sql.String()
}
func (EntityArticleView) Migration() map[string]string {
	var view = EntityArticleView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
