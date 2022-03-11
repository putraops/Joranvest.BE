package entity_view_models

import (
	"joranvest/models/view_models"
	"strings"
)

type EntityEducationPlaylistUserView struct {
	view_models.BaseViewModel
	EducationPlaylistId     string `json:"education_playlist_id"`
	EducationId             string `json:"education_id"`
	Title                   string `json:"title"`
	FileUrl                 string `json:"file_url"`
	Description             string `json:"description"`
	OrderIndex              int    `json:"order_index"`
	ApplicationUserId       string `json:"application_user_id"`
	ApplicationUserFullname string `json:"application_user_fullname"`
}

func (EntityEducationPlaylistUserView) TableName() string {
	return "vw_education_playlist_user"
}

func (EntityEducationPlaylistUserView) ViewModel() string {
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
	sql.WriteString("  r.application_user_id,")
	sql.WriteString("  r.education_playlist_id,")
	sql.WriteString("  e.education_id,")
	sql.WriteString("  e.title AS education_playlist_title,")
	sql.WriteString("  e.file_url AS education_Playlist_file_url,")
	sql.WriteString("  e.description AS education_playlist_description,")
	sql.WriteString("  e.order_index education_playlist_order_index,")
	sql.WriteString("  CASE WHEN u5.first_name IS NULL OR u5.first_name = '' THEN u5.username ELSE concat(u5.first_name, ' ', u5.last_name) END AS application_user_fullname, ")
	sql.WriteString("  CASE WHEN u1.first_name IS NULL OR u1.first_name = '' THEN u1.username ELSE concat(u1.first_name, ' ', u1.last_name) END AS created_user_fullname,")
	sql.WriteString("  CASE WHEN u2.first_name IS NULL OR u2.first_name = '' THEN u2.username ELSE concat(u2.first_name, ' ', u2.last_name) END AS updated_user_fullname,")
	sql.WriteString("  CASE WHEN u3.first_name IS NULL OR u3.first_name = '' THEN u3.username ELSE concat(u3.first_name, ' ', u3.last_name) END AS submitted_user_fullname,")
	sql.WriteString("  CASE WHEN u4.first_name IS NULL OR u4.first_name = '' THEN u4.username ELSE concat(u4.first_name, ' ', u4.last_name) END AS approved_user_fullname ")
	sql.WriteString("FROM education_playlist_user r ")
	sql.WriteString("LEFT JOIN education_playlist e ON e.id = r.education_playlist_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	sql.WriteString("LEFT JOIN application_user u4 ON u4.id = r.approved_by ")
	sql.WriteString("LEFT JOIN application_user u5 ON u5.id = r.application_user_id ")
	return sql.String()
}
func (EntityEducationPlaylistUserView) Migration() map[string]string {
	var view = EntityEducationPlaylistUserView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
