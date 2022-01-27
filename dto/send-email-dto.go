package dto

type SendWebinarInformationDto struct {
	WebinarId  string `json:"webinar_id" form:"webinar_id"`
	MeetingUrl string `json:"meeting_url" form:"meeting_url"`
	MeetingId  string `json:"meeting_id" form:"meeting_id"`
	Passcode   string `json:"passcode" form:"passcode"`
}
