package dto

import "time"

//WebinarDto is a model that client use when updating a book
type WebinarDto struct {
	Id                      string    `json:"id" form:"id"`
	WebinarCategoryId       string    `json:"webinar_category_id" form:"webinar_category_id" binding:"required"`
	Title                   string    `json:"title" form:"title"`
	Description             string    `json:"description" form:"description"`
	WebinarStartDate        time.Time `json:"webinar_start_date" form:"webinar_start_date"`
	WebinarEndDate          time.Time `json:"webinar_end_date" form:"webinar_end_date"`
	OrganizerOrganizationId string    `json:"organizer_organization_id" form:"organizer_organization_id"`
	WebinarSpeaker          string    `json:"webinar_speaker" form:"webinar_speaker"`
	MinAge                  int       `json:"min_age" form:"min_age"`
	WebinarLevel            string    `json:"webinar_level" form:"webinar_level"`
	Price                   float64   `json:"price" form:"price"`
	Discount                float64   `json:"discount" form:"discount"`
	IsCertificate           bool      `json:"is_certificate" form:"is_certificate"`
	Reward                  int       `json:"reward" form:"reward"`
	Status                  int       `json:"status" form:"status"`
	SpeakerType             int       `json:"speaker_type" form:"speaker_type"`

	EntityId  string `json:"-"`
	UpdatedBy string
}

//WebinarDto is a model that client use when updating a book
type WebinarSepakerDto struct {
	Id             string `json:"id" form:"id"`
	WebinarId      string `json:"webinar_id" form:"webinar_id"`
	WebinarSpeaker string `json:"webinar_speaker" form:"webinar_speaker"`
	SpeakerType    int    `json:"speaker_type" form:"speaker_type"`
	EntityId       string `json:"-"`
	UpdatedBy      string
}
