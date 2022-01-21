package controllers

import (
	"fmt"
	"net/http"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type WebinarRegistrationController interface {
	GetDatatables(context *gin.Context)
	GetPagination(context *gin.Context)
	GetById(context *gin.Context)
	GetViewById(context *gin.Context)
	Save(context *gin.Context)
	IsWebinarRegistered(context *gin.Context)
	DeleteById(context *gin.Context)

	SendWebinarInformationByWebinarId(context *gin.Context)
}

type webinarRegistrationController struct {
	webinarRegistrationService service.WebinarRegistrationService
	jwtService                 service.JWTService
}

func NewWebinarRegistrationController(webinarRegistrationService service.WebinarRegistrationService, jwtService service.JWTService) WebinarRegistrationController {
	return &webinarRegistrationController{
		webinarRegistrationService: webinarRegistrationService,
		jwtService:                 jwtService,
	}
}

func (c *webinarRegistrationController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.webinarRegistrationService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *webinarRegistrationController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.webinarRegistrationService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

func (c *webinarRegistrationController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarRegistrationService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *webinarRegistrationController) GetViewById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarRegistrationService.GetViewById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *webinarRegistrationController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.WebinarRegistrationDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.WebinarRegistration{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))

		newRecord.ApplicationUserId = userIdentity.UserId
		newRecord.EntityId = userIdentity.EntityId
		newRecord.CreatedBy = userIdentity.UserId
		newRecord.OwnerId = userIdentity.UserId
		result = c.webinarRegistrationService.Insert(newRecord)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *webinarRegistrationController) IsWebinarRegistered(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	authHeader := context.GetHeader("Authorization")
	userIdentity := c.jwtService.GetUserByToken(authHeader)

	result := c.webinarRegistrationService.IsWebinarRegistered(id, userIdentity.UserId)

	var message = "Ok"
	if !result.Status {
		message = "Error"
	}

	response := helper.BuildResponse(result.Status, message, result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *webinarRegistrationController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.webinarRegistrationService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *webinarRegistrationController) SendWebinarInformationByWebinarId(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}

	result := c.webinarRegistrationService.SendWebinarInformationViaEmail(id)
	response := helper.BuildResponse(result.Status, result.Message, helper.EmptyObj{})
	context.JSON(http.StatusOK, response)
}
