package dto

import "github.com/gin-gonic/gin"

type EducationCategoryDto struct {
	Id          string `json:"id" form:"id"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
	ParentId    string `json:"parent_id" form:"parent_id"`
	EntityId    string `json:"-"`
	Context     *gin.Context
}
