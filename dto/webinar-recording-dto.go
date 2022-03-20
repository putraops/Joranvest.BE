package dto

import "github.com/gin-gonic/gin"

type WebinarRecordingDto struct {
	Id        string  `json:"id" form:"id"`
	WebinarId *string `json:"webinar_id" form:"webinar_id" binding:"required"`
	VideoUrl  *string `json:"video_url" form:"video_url" binding:"required"`
	Context   *gin.Context
}
