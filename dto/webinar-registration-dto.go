package dto

type WebinarRegistrationDto struct {
	Id                string `json:"id" form:"id"`
	WebinarId         string `json:"webinar_id" form:"webinar_id" binding:"required"`
	ApplicationUserId string `json:"application_user_id" form:"application_user_id" binding:"required"`
	PaymentId         string `json:"payment_id" form:"payment_id"`

	EntityId  string `json:"-"`
	UpdatedBy string
}

type WebinarRegistrationUpdatePaymentDto struct {
	Id            string `json:"id" form:"id"`
	WebinarId     string `json:"webinar_id" form:"webinar_id" binding:"required"`
	PaymentStatus int    `json:"payment_status" form:"payment_status"`

	EntityId  string `json:"-"`
	UpdatedBy string
}
