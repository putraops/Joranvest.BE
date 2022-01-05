package view_models

import (
	"strings"
)

type WebinarSpeakerReviewViewModel struct {
	Id                string  `json:"id"`
	SpeakerName       string  `json:"speaker_name"`
	SpeakerTitle      string  `json:"speaker_title"`
	Rating            float32 `json:"rating"`
	TotalRating       int     `json:"total_rating"`
	Filepath          string  `json:"filepath"`
	FilepathThumbnail string  `json:"filepath_thumbnail"`
	Filename          string  `json:"filename"`
	IsOrganization    bool    `json:"is_organization"`
	Description       string  `json:"description"`
}

func (WebinarSpeakerReviewViewModel) TableName() string {
	return "vw_webinar_speaker_review"
}

func (WebinarSpeakerReviewViewModel) ViewModel() string {
	var sql strings.Builder
	// sql.WriteString("SELECT")
	sql.WriteString("SELECT * ")
	sql.WriteString("FROM ( ")
	sql.WriteString("	SELECT")
	sql.WriteString("		r.id,")
	sql.WriteString("		CONCAT(r.first_name, ' ', r.last_name) AS speaker_name,")
	sql.WriteString("		COALESCE(m.rating, 0)::float AS rating,")
	sql.WriteString("		COALESCE(m.total_rating, 0)::integer AS total_rating,")
	sql.WriteString("		r.filepath,")
	sql.WriteString("		r.filepath_thumbnail,")
	sql.WriteString("		r.filename,")
	sql.WriteString("		r.extension,")
	sql.WriteString("		r.title AS speaker_title,")
	sql.WriteString("		0::boolean AS is_organization,")
	sql.WriteString("		r.description")
	sql.WriteString("	FROM application_user r")
	sql.WriteString("	LEFT JOIN LATERAL get_webinar_speaker_rating(r.id) m(rating, total_rating) ON true")
	sql.WriteString("	UNION ALL")
	sql.WriteString("	SELECT ")
	sql.WriteString("		r.id,")
	sql.WriteString("		r.name  AS speaker_name,")
	sql.WriteString("		COALESCE(m.rating, 0)::float AS rating,")
	sql.WriteString("		COALESCE(m.total_rating, 0)::integer AS total_rating,")
	sql.WriteString("		r.filepath,")
	sql.WriteString("		r.filepath_thumbnail,")
	sql.WriteString("		r.filename,")
	sql.WriteString("		r.extension,")
	sql.WriteString("		null AS speaker_title,")
	sql.WriteString("		1::boolean AS is_organization,")
	sql.WriteString("		r.description")
	sql.WriteString("	FROM organization r")
	sql.WriteString("	LEFT JOIN LATERAL get_webinar_speaker_rating(r.id) m(rating, total_rating) ON true ")
	sql.WriteString(") AS r")
	return sql.String()
}
func (WebinarSpeakerReviewViewModel) Migration() map[string]string {
	var view = WebinarSpeakerReviewViewModel{}
	var m = make(map[string]string)
	m["view_name"] = view.TableName()
	m["query"] = view.ViewModel()
	return m
}
