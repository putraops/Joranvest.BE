package dto

//MembershipUserDto is a model that client use when updating a book
type MembershipUserDto struct {
	Id                string `json:"id" form:"id"`
	ApplicationUserId string `json:"application_user_id" form:"application_user_id" binding:"required"`
	MembershipId      string `json:"membership_id" form:"membership_id" binding:"required"`
	PaymentStatus     int    `json:"payment_status" form:"payment_status"`
	PaymentType       string `json:"payment_type" form:"payment_type"`
}
