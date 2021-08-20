package dto

//LoginDTO is a model that used by client when POST from /login url
type LoginDTO struct {
	Credential string `json:"credential" form:"credential" binding:"required"`
	Password   string `json:"password" form:"password" binding:"required"`
}
