package controllers

import (
	"net/http"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/repository"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebinarRecordingController interface {
	GetPagination(context *gin.Context)
	Save(context *gin.Context)
	Submit(context *gin.Context)
	GetById(context *gin.Context)
	GetViewById(context *gin.Context)
	GetByWebinarId(context *gin.Context)
	GetByPathUrl(context *gin.Context)
	DeleteById(context *gin.Context)
}

type webinarRecordingController struct {
	webinarRecordingService    service.WebinarRecordingService
	webinarRecordingRepository repository.WebinarRecordingRepository
	jwtService                 service.JWTService
	db                         *gorm.DB
}

func NewWebinarRecordingController(db *gorm.DB, jwtService service.JWTService) WebinarRecordingController {
	return &webinarRecordingController{
		db:                         db,
		jwtService:                 jwtService,
		webinarRecordingService:    service.NewWebinarRecordingService(db, jwtService),
		webinarRecordingRepository: repository.NewWebinarRecordingRepository(db),
	}
}

// @Tags         WebinarRecording
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body commons.Pagination2ndRequest true "body"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /webinar_recording/getPagination [post]
func (c webinarRecordingController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.webinarRecordingService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

// @Tags         WebinarRecording
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.WebinarRecordingDto true "dto"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /webinar_recording/save [post]
func (r webinarRecordingController) Save(c *gin.Context) {
	var result helper.Result
	var dto dto.WebinarRecordingDto
	dto.Context = c

	errDto := c.Bind(&dto)
	if errDto != nil {
		res := helper.StandartResult(false, errDto.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result = r.webinarRecordingService.Save(dto)
	c.JSON(http.StatusOK, helper.StandartResult(result.Status, result.Message, result.Data))
	return
}

// @Tags         WebinarRecording
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /webinar_recording/submit/{id} [get]
func (c webinarRecordingController) Submit(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarRecordingService.Submit(id, context)
	context.JSON(http.StatusOK, result)
}

// @Tags         WebinarRecording
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /webinar_recording/getById/{id} [get]
func (c webinarRecordingController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarRecordingService.GetById(id)
	context.JSON(http.StatusOK, result)
}

// @Tags         WebinarRecording
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /webinar_recording/getViewById/{id} [get]
func (c webinarRecordingController) GetViewById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarRecordingService.GetViewById(id)
	context.JSON(http.StatusOK, result)
}

// @Tags         WebinarRecording
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        webinarId path string true "webinarId"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /webinar_recording/getByWebinarId/{webinarId} [get]
func (c webinarRecordingController) GetByWebinarId(context *gin.Context) {
	webinarId := context.Param("webinarId")
	if webinarId == "" {
		response := helper.BuildErrorResponse("Failed to get webinarId", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarRecordingService.GetByWebinarId(webinarId)
	context.JSON(http.StatusOK, result)
}

// @Tags         WebinarRecording
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        path_url path string true "path_url"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /webinar_recording/getByPathUrl/{path_url} [get]
func (c webinarRecordingController) GetByPathUrl(context *gin.Context) {
	path_url := context.Param("path_url")
	if path_url == "" {
		response := helper.BuildErrorResponse("Failed to get by path_url", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarRecordingRepository.GetByPathUrl(path_url)
	context.JSON(http.StatusOK, result)
}

// @Tags         WebinarRecording
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /webinar_recording/deleteById/{id} [delete]
func (c webinarRecordingController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarRecordingService.DeleteById(id)
	context.JSON(http.StatusOK, result)
}
