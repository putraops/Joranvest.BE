package dto

import "github.com/gin-gonic/gin"

type EducationDto struct {
	Id                  string `json:"id" form:"id"`
	EducationCategoryId string `json:"education_category_id" form:"education_category_id" binding:"required"`
	Title               string `json:"title" form:"title" binding:"required"`
	Level               string `json:"level" form:"level" binding:"required"`
	Description         string `json:"description" form:"description"`
	EntityId            string `json:"-"`
	Context             *gin.Context
}

type EducationPlaylistDto struct {
	Id          string `json:"id" form:"id"`
	EducationId string `json:"education_id" form:"education_id" binding:"required"`
	Title       string `json:"title" form:"title" binding:"required"`
	FileUrl     string `json:"file_url" form:"file_url" binding:"required"`
	Description string `json:"description" form:"description"`
	EntityId    string `json:"-"`
	Context     *gin.Context
}
