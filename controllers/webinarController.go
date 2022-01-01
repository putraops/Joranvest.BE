package controllers

import (
	"fmt"
	"net/http"
	"time"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
	log "github.com/sirupsen/logrus"
)

type WebinarController interface {
	GetDatatables(context *gin.Context)
	GetPagination(context *gin.Context)
	GetPaginationRegisteredByUser(context *gin.Context)
	DeleteById(context *gin.Context)
	GetById(context *gin.Context)
	GetWebinarWithRatingByUserId(context *gin.Context)
	Save(context *gin.Context)
	Submit(context *gin.Context)
}

type webinarController struct {
	webinarService service.WebinarService
	jwtService     service.JWTService
}

func NewWebinarController(webinarService service.WebinarService, jwtService service.JWTService) WebinarController {
	return &webinarController{
		webinarService: webinarService,
		jwtService:     jwtService,
	}
}

func (c *webinarController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.webinarService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

// @Tags Webinar
// @Summary Get Pagination
// @Param id path string true "id"
// @Router /webinar/getPagination [post]
// @Success 200 {obsject} object
// @Failure 400,404 {object} object
// @Router /webinar/getPagination [get]
func (c *webinarController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.webinarService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

// @Tags Webinar
// @Summary Get Pagination
// @Param id path string true "id"
// @Router /webinar/getPagination [post]
// @Success 200 {obsject} object
// @Failure 400,404 {object} object
// @Router /webinar/getPagination [get]
func (c *webinarController) GetPaginationRegisteredByUser(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}

	userId := context.Param("user_id")
	if userId == "" {
		response := helper.BuildErrorResponse("Failed to get userId", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	fmt.Println("======================")
	fmt.Println(userId)
	fmt.Println("======================")

	var result = c.webinarService.GetPaginationRegisteredByUser(req, userId)
	context.JSON(http.StatusOK, result)
}

func (c *webinarController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.WebinarDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.Webinar{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId

		// myDateString := "2018-01-20 00:00:00"
		//-- Mapping WebinarStartDate
		fmt.Println(recordDto.WebinarStartDate.String())
		fmt.Println(recordDto.WebinarStartDate.String())
		startDate, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", recordDto.WebinarStartDate.String())
		if err != nil {
			panic(err)
		}
		newRecord.WebinarStartDate.Time = startDate
		newRecord.WebinarStartDate.Valid = true

		endDate, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", recordDto.WebinarEndDate.String())
		if err != nil {
			panic(err)
		}
		newRecord.WebinarEndDate.Time = endDate
		newRecord.WebinarEndDate.Valid = true

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			newRecord.OwnerId = userIdentity.UserId
			result = c.webinarService.Insert(newRecord)
		} else {
			newRecord.UpdatedBy = userIdentity.UserId
			result = c.webinarService.Update(newRecord)
		}

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

// @Tags Webinar
// @Summary Get Webinar By Id
// @Param id path string true "id"
// @Router /webinar/getById [delete]
// @Success 200 {object} helper.Response
// @Failure 400,404 {object} object
// @Router /webinar/getById/{id} [get]
func (c *webinarController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *webinarController) GetWebinarWithRatingByUserId(context *gin.Context) {
	commons.Logger()

	webinar_id := context.Param("webinar_id")
	user_id := context.Param("user_id")
	if user_id == "" || webinar_id == "" {
		log.Error("Failed to get User Id or Webinar Id :: Function Name: GetWebinarWithRatingByUserId")
		response := helper.BuildErrorResponse("Failed to get User Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarService.GetWebinarWithRatingByUserId(webinar_id, user_id)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

// @Tags Webinar
// @Summary Delete Webinar By Id
// @Param id path string true "id"
// @Router /webinar/deleteById [delete]
// @Success 200 {object} helper.Response
// @Failure 400,404 {object} object
// @Router /webinar/getById/{id} [get]
func (c *webinarController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.webinarService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

// @Tags Webinar
// @Summary Submit Webinar By Id
// @Param id path string true "id"
// @Router /webinar/submit/{id} [post]
// @Success 200 {object} helper.Response
// @Failure 400,404 {object} object
func (c *webinarController) Submit(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}

	authHeader := context.GetHeader("Authorization")
	userIdentity := c.jwtService.GetUserByToken(authHeader)
	var result = c.webinarService.Submit(id, userIdentity.UserId)

	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
