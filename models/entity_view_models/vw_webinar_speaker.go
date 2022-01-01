package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityWebinarSpeakerView struct {
	models.WebinarSpeaker
	OrganizationName   string `json:"organization_name"`
	SpeakerFirstName   string `json:"speaker_first_name"`
	SpeakerLastName    string `json:"speaker_last_name"`
	SpeakerFullName    string `json:"speaker_full_name"`
	SpeakerInitialName string `json:"speaker_initial_name"`
	SpeakerTitle       string `json:"speaker_title"`
	UserCreate         string `json:"user_create"`
	UserUpdate         string `json:"user_update"`
}

func (EntityWebinarSpeakerView) TableName() string {
	return "vw_webinar_speaker"
}

func (EntityWebinarSpeakerView) ViewModel() string {
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
	sql.WriteString("  r.owner_id,")
	sql.WriteString("  r.entity_id,")
	sql.WriteString("  r.webinar_id,")
	sql.WriteString("  w.title AS webinar_title,")
	sql.WriteString("  r.speaker_id,")
	sql.WriteString("  o.name AS organization_name,")
	sql.WriteString("  u3.first_name AS speaker_first_name,")
	sql.WriteString("  u3.last_name AS speaker_last_name,")
	sql.WriteString("  u3.title AS speaker_title,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS speaker_full_name,")
	sql.WriteString("  CONCAT(UPPER(LEFT(u3.first_name, 1)), '', UPPER(LEFT(u3.last_name, 1))) AS speaker_initial_name,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString(" FROM webinar_speaker r ")
	sql.WriteString("  LEFT JOIN webinar w ON w.id = r.webinar_id")
	sql.WriteString("  LEFT JOIN organization o ON o.id = r.speaker_id")
	sql.WriteString("  LEFT JOIN application_user u3 ON u3.id = r.speaker_id")
	sql.WriteString("  LEFT JOIN application_user u1 ON u1.id = r.created_by")
	sql.WriteString("  LEFT JOIN application_user u2 ON u2.id = r.updated_by")
	return sql.String()
}
func (EntityWebinarSpeakerView) Migration() map[string]string {
	var view = EntityWebinarSpeakerView{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
