package dto

type WebinarRegistrationDto struct {
	Id            string `json:"id" form:"id"`
	WebinarId     string `json:"webinar_id" form:"webinar_id" binding:"required"`
	PaymentType   string `json:"payment_type" form:"payment_type"`
	PaymentStatus int    `json:"payment_status" form:"payment_status"`

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
