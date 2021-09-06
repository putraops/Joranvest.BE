package dto

import "time"

//WebinarDto is a model that client use when updating a book
type WebinarDto struct {
	Id                      string    `json:"id" form:"id"`
	WebinarCategoryId       string    `json:"webinar_category_id" form:"webinar_category_id" binding:"required"`
	Title                   string    `json:"title" form:"title"`
	Description             string    `json:"description" form:"description"`
	WebinarFirstStartDate   time.Time `json:"webinar_first_start_date" form:"webinar_first_start_date"`
	WebinarFirstEndDate     time.Time `json:"webinar_first_end_date" form:"webinar_first_end_date"`
	WebinarLastStartDate    time.Time `json:"webinar_last_start_date" form:"webinar_last_start_date"`
	WebinarLastEndDate      time.Time `json:"webinar_last_end_date" form:"webinar_last_end_date"`
	OrganizerOrganizationId string    `json:"organizer_organization_id" form:"organizer_organization_id"`
	WebinarSpeaker          string    `json:"webinar_speaker" form:"webinar_speaker"`
	MinAge                  int       `json:"min_age" form:"min_age"`
	WebinarLevel            string    `json:"webinar_level" form:"webinar_level"`
	Price                   float64   `json:"price" form:"price"`
	Discount                float64   `json:"discount" form:"discount"`
	IsCertificate           bool      `json:"is_certificate" form:"is_certificate"`
	Reward                  int       `json:"reward" form:"reward"`
	Status                  string    `json:"status" form:"status"`

	EntityId  string `json:"-"`
	UpdatedBy string
}