package dto

//-- RoleMemberDto is a model that client use when updating a book
type RoleMemberDto struct {
	Id                string `json:"id" form:"id"`
	RoleId            string `json:"role_id" form:"role_id" binding:"required"`
	ApplicationUserId string `json:"application_user_id" form:"application_user_id" binding:"required"`
	EntityId          string `json:"-"`
	UpdatedBy         string
}
