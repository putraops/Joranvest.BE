package dto

//LoginDTO is a model that used by client when POST from /login url
type LoginDto struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password" binding:"required"`
}
