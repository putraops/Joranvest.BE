package dto

type CustomerDto struct {
	Id         string `json:"id" form:"id"`
	EntityId   string `json:"-"`
	TenantId   string `json:"-"`
	FirstName  string `json:"first_name" form:"first_name" binding:"required"`
	LastName   string `json:"last_name" form:"last_name" binding:"required"`
	Address    string `json:"address" form:"address" binding:"required"`
	Username   string `json:"username" form:"username"`
	Password   string `json:"password" form:"password"`
	Email      string `json:"email" form:"email"`
	Phone      string `json:"phone" form:"phone" binding:"required"`
	CreatedBy  string `json:"-"`
	ApprovedBy string `json:"-"`
	UpdatedBy  string `json:"-"`
}

type CustomerUpdatePasswordDto struct {
	Id       string `json:"id" form:"id"`
	Password string `json:"password" form:"password" binding:"required"`
}
