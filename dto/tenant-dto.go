package dto

type TenantCreateDTO struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Address     string `json:"address" form:"address" binding:"required"`
	Phone       string `json:"phone" form:"phone" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
}

type TenantUpdateDTO struct {
	Id          string `json:"id" form:"id" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	Address     string `json:"address" form:"address" binding:"required"`
	Phone       string `json:"phone" form:"phone" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
}

type TenantDto struct {
	Id          string `json:"id" form:"id"`
	Name        string `json:"name" form:"name" binding:"required"`
	Address     string `json:"address" form:"address" binding:"required"`
	Phone       string `json:"phone" form:"phone" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
}
