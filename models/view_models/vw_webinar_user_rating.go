package view_models

import (
	"joranvest/models"
)

type WebinarUserRatingViewModel struct {
	models.Webinar
	RatingMasterId            string `json:"rating_master_id"`
	Rating                    int    `json:"rating"`
	Comment                   string `json:"comment"`
	OrganizerOrganizationName string `json:"organizer_organization_name"`
	SpeakerName               string `json:"speaker_name"`
	WebinarCategoryName       string `json:"webinar_category_name"`
	CreatedByFullname         string `json:"created_by_fullname"`
	UpdatedByFullname         string `json:"updated_by_fullname"`
	SubmittedFullname         string `json:"submitted_by_fullname"`
}
