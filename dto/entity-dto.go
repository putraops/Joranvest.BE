package dto

type EntityCreateDto struct {
	Name        string `json:"name" form:"name" binding:"required"`
	OwnerName   string `json:"owner_name" form:"owner_name" binding:"required"`
	Address     string `json:"address" form:"address" binding:"required"`
	Phone       string `json:"phone" form:"phone" binding:"required"`
	Email       string `json:"email" form:"email" binding:"required"`
	Description string `json:"description" form:"description"`
}

type EntityUpdateDto struct {
	Id          string `json:"id" form:"id" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	OwnerName   string `json:"owner_name" form:"owner_name" binding:"required"`
	Address     string `json:"address" form:"address" binding:"required"`
	Phone       string `json:"phone" form:"phone" binding:"required"`
	Email       string `json:"email" form:"email" binding:"required,email"`
	Description string `json:"description" form:"description"`
}

type EntityDto struct {
	Id            string `json:"id" form:"id"`
	AdminUsername string `json:"admin_username" form:"admin_username" binding:"required"`
	Name          string `json:"name" form:"name" binding:"required"`
	OwnerName     string `json:"owner_name" form:"owner_name" binding:"required"`
	Address       string `json:"address" form:"address" binding:"required"`
	Phone         string `json:"phone" form:"phone" binding:"required"`
	Email         string `json:"email" form:"email" binding:"required,email"`
	Description   string `json:"description" form:"description"`
}
