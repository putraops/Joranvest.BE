package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityWebinarView struct {
	models.Webinar
	Filepath                  string `json:"filepath"`
	FilepathThumbnail         string `json:"filepath_thumbnail"`
	Filename                  string `json:"filename"`
	Extension                 string `json:"extension"`
	OrganizerOrganizationName string `json:"organizer_organization_name"`
	SpeakerName               string `json:"speaker_name"`
	WebinarCategoryName       string `json:"webinar_category_name"`
	CreatedByFullname         string `json:"created_by_fullname"`
	UpdatedByFullname         string `json:"updated_by_fullname"`
	SubmittedFullname         string `json:"submitted_by_fullname"`
}

func (EntityWebinarView) TableName() string {
	return "vw_webinar"
}

func (EntityWebinarView) ViewModel() string {
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
	sql.WriteString("  r.owner_id,")
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.webinar_category_id,")
	sql.WriteString("  c.name AS webinar_category_name,")
	sql.WriteString("  (SELECT s.speaker FROM (SELECT ws.webinar_id, string_agg(ws.speaker_full_name, ', ') AS speaker FROM vw_webinar_speaker ws WHERE ws.webinar_id = r.id GROUP BY 1) AS s) AS speaker_name,")
	sql.WriteString("  r.title,")
	sql.WriteString("  r.description,")
	sql.WriteString("  r.webinar_start_date,")
	sql.WriteString("  r.webinar_end_date,")
	sql.WriteString("  r.min_age,")
	sql.WriteString("  r.webinar_level,")
	sql.WriteString("  r.price,")
	sql.WriteString("  r.discount,")
	sql.WriteString("  r.is_certificate,")
	sql.WriteString("  r.reward,")
	sql.WriteString("  r.status,")
	sql.WriteString("  r.speaker_type,")
	sql.WriteString("  r.filepath,")
	sql.WriteString("  r.filepath_thumbnail,")
	sql.WriteString("  r.filename,")
	sql.WriteString("  r.extension,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS created_by_fullname,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS updated_by_fullname, ")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS submitted_by_fullname ")
	sql.WriteString("FROM public.webinar r ")
	sql.WriteString("  LEFT JOIN webinar_category c ON c.id = r.webinar_category_id")
	//sql.WriteString("  LEFT JOIN filemaster f ON f.record_id = r.id")
	sql.WriteString("  LEFT JOIN application_user u1 ON u1.id = r.created_by")
	sql.WriteString("  LEFT JOIN application_user u2 ON u2.id = r.updated_by")
	sql.WriteString("  LEFT JOIN application_user u3 ON u3.id = r.submitted_by")
	return sql.String()
}

func (EntityWebinarView) Migration() map[string]string {
	var view = EntityWebinarView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
