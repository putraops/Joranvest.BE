package dto

type ApplicationUserRegisterDto struct {
	// Name     string `json:"name" form:"name" binding:"required"`
	EntityId  string `json:"-"`
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	// Address   string `json:"address" form:"address" binding:"required"`
	// Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	// Phone    string `json:"phone" form:"phone" binding:"required"`
	UserType string `json:"user_type" form:"user_type"`
	IsAdmin  bool   `json:"is_admin" form:"is_admin"`
}

type ApplicationUserUpdateDto struct {
	Id       string `json:"id" form:"id"`
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password,omitempty" form:"password,omitempty"`
}
