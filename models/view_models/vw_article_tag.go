package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityArticleTagView struct {
	models.ArticleTag
	TagName        string `json:"tag_name"`
	TagDescription string `json:"tag_description"`
	UserCreate     string `json:"user_create"`
	UserUpdate     string `json:"user_update"`
}

func (EntityArticleTagView) TableName() string {
	return "vw_article_tag"
}

func (EntityArticleTagView) ViewModel() string {
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
	sql.WriteString("  r.article_id,")
	sql.WriteString("  r.tag_id,")
	sql.WriteString("  t.name AS tag_name,")
	sql.WriteString("  t.description AS tag_description,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString(" FROM article_tag r ")
	sql.WriteString("  LEFT JOIN article a ON a.id = r.article_id")
	sql.WriteString("  LEFT JOIN tag t ON t.id = r.tag_id")
	sql.WriteString("  LEFT JOIN application_user u1 ON u1.id = r.created_by")
	sql.WriteString("  LEFT JOIN application_user u2 ON u2.id = r.updated_by")
	return sql.String()
}
func (EntityArticleTagView) Migration() map[string]string {
	var view = EntityArticleTagView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
