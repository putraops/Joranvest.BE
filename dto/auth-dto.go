package dto

//LoginDTO is a model that used by client when POST from /login url
type LoginDto struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password" binding:"required"`
}

//ChangePasswordDto is a model that used by client when POST from /login url
type ChangePasswordDto struct {
	Username    string `json:"username" form:"username"`
	Email       string `json:"email" form:"email"`
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required"`
}

type RecoverPasswordDto struct {
	Id          string `json:"id" form:"id"`
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
}
