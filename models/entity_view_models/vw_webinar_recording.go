package entity_view_models

import (
	"joranvest/models/view_models"
	"strings"
)

type EntityWebinarRecordingView struct {
	view_models.BaseViewModel
	WebinarId           string  `json:"webinar_id"`
	VideoUrl            string  `json:"video_url"`
	PathUrl             string  `json:"path_url"`
	WebinarCategoryId   string  `json:"webinar_category_id"`
	WebinarCategoryName string  `json:"webinar_category_name"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	Price               float64 `json:"price"`
	WebinarLevel        string  `json:"webinar_level"`
	Filepath            string  `json:"filepath"`
	FilepathThumbnail   string  `json:"filepath_thumbnail"`
	Filename            string  `json:"filename"`
	Extension           string  `json:"extension"`
	Size                string  `json:"size"`
}

func (EntityWebinarRecordingView) TableName() string {
	return "vw_webinar_recording"
}

func (EntityWebinarRecordingView) ViewModel() string {
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
	sql.WriteString("  CASE WHEN r.submitted_at IS NOT NULL THEN true ELSE false END AS is_submitted,")
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.webinar_id,")
	sql.WriteString("  r.video_url,")
	sql.WriteString("  r.path_url,")
	sql.WriteString("  wc.id AS webinar_category_id,")
	sql.WriteString("  wc.name AS webinar_category_name,")
	sql.WriteString("  w.title,")
	sql.WriteString("  w.description,")
	sql.WriteString("  w.webinar_level,")
	sql.WriteString("  w.price,")
	sql.WriteString("  w.filepath,")
	sql.WriteString("  w.filepath_thumbnail,")
	sql.WriteString("  w.filename,")
	sql.WriteString("  w.extension,")
	sql.WriteString("  w.size,")
	sql.WriteString("  CASE WHEN u1.first_name IS NULL OR u1.first_name = '' THEN u1.username ELSE concat(u1.first_name, ' ', u1.last_name) END AS created_user_fullname,")
	sql.WriteString("  CASE WHEN u2.first_name IS NULL OR u2.first_name = '' THEN u2.username ELSE concat(u2.first_name, ' ', u2.last_name) END AS updated_user_fullname,")
	sql.WriteString("  CASE WHEN u3.first_name IS NULL OR u3.first_name = '' THEN u3.username ELSE concat(u3.first_name, ' ', u3.last_name) END AS submitted_user_fullname,")
	sql.WriteString("  CASE WHEN u4.first_name IS NULL OR u4.first_name = '' THEN u4.username ELSE concat(u4.first_name, ' ', u4.last_name) END AS approved_user_fullname ")
	sql.WriteString("FROM webinar_recording r ")
	sql.WriteString("LEFT JOIN webinar w ON w.id = r.webinar_id ")
	sql.WriteString("LEFT JOIN webinar_category wc ON wc.id = w.webinar_category_id ")
	sql.WriteString("LEFT JOIN application_user u1 ON u1.id = r.created_by ")
	sql.WriteString("LEFT JOIN application_user u2 ON u2.id = r.updated_by ")
	sql.WriteString("LEFT JOIN application_user u3 ON u3.id = r.submitted_by ")
	sql.WriteString("LEFT JOIN application_user u4 ON u4.id = r.approved_by ")
	return sql.String()
}
func (EntityWebinarRecordingView) Migration() map[string]string {
	var view = EntityWebinarRecordingView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
