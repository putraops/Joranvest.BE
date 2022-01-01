package view_models

import (
	"joranvest/models"
)

type WebinarUserViewModel struct {
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
