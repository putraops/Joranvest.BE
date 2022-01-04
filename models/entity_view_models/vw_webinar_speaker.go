package entity_view_models

import (
	"joranvest/models"
	"strings"
)

type EntityWebinarSpeakerView struct {
	models.WebinarSpeaker
	Rating                                     float32 `json:"rating"`
	TotalRating                                int     `json:"total_rating"`
	OrganizationName                           string  `json:"organization_name"`
	OrganizationSpeakerProfilePicture          string  `json:"organization_speaker_profile_picture"`
	OrganizationSpeakerProfilePictureThumbnail string  `json:"organization_speaker_profile_picture_thumbnail"`
	OrganizationSpeakerProfilePictureExtension string  `json:"organization_speaker_profile_picture_extension"`

	SpeakerFirstName                   string `json:"speaker_first_name"`
	SpeakerLastName                    string `json:"speaker_last_name"`
	SpeakerFullName                    string `json:"speaker_full_name"`
	SpeakerInitialName                 string `json:"speaker_initial_name"`
	SpeakerTitle                       string `json:"speaker_title"`
	UserSpeakerProfilePicture          string `json:"user_speaker_profile_picture"`
	UserSpeakerProfilePictureThumbnail string `json:"user_speaker_profile_picture_thumbnail"`
	UserSpeakerProfilePictureExtension string `json:"user_speaker_profile_picture_extension"`

	UserCreate string `json:"user_create"`
	UserUpdate string `json:"user_update"`
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
	sql.WriteString("  COALESCE(m.rating, 0) AS rating,")
	sql.WriteString("  COALESCE(m.total_rating, 0) AS total_rating,")
	sql.WriteString("  o.name AS organization_name,")
	sql.WriteString("  u3.first_name AS speaker_first_name,")
	sql.WriteString("  u3.last_name AS speaker_last_name,")
	sql.WriteString("  u3.title AS speaker_title,")
	sql.WriteString("  o.filepath AS organization_speaker_profile_picture,")
	sql.WriteString("  o.filepath_thumbnail AS organization_speaker_profile_picture_thumbnail,")
	sql.WriteString("  o.extension AS organization_speaker_profile_picture_extension,")
	sql.WriteString("  u3.filepath AS user_speaker_profile_picture,")
	sql.WriteString("  u3.filepath_thumbnail AS user_speaker_profile_picture_thumbnail,")
	sql.WriteString("  u3.extension AS user_speaker_profile_picture_extension,")
	sql.WriteString("  CONCAT(u3.first_name, ' ', u3.last_name) AS speaker_full_name,")
	sql.WriteString("  CONCAT(UPPER(LEFT(u3.first_name, 1)), '', UPPER(LEFT(u3.last_name, 1))) AS speaker_initial_name,")
	sql.WriteString("  CONCAT(u1.first_name, ' ', u1.last_name) AS user_create,")
	sql.WriteString("  CONCAT(u2.first_name, ' ', u2.last_name) AS user_update ")
	sql.WriteString(" FROM webinar_speaker r ")
	sql.WriteString("  LEFT JOIN webinar w ON w.id = r.webinar_id")
	sql.WriteString("  LEFT JOIN LATERAL get_webinar_speaker_rating(r.speaker_id) m(rating, total_rating) ON true")
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
