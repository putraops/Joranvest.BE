package controllers

import (
	"fmt"
	"net/http"

	"joranvest/helper"
	"joranvest/models"
	"joranvest/service"

	"github.com/gin-gonic/gin"
)

type FilemasterController interface {
	GetAll(context *gin.Context)
	Insert(context *gin.Context)
	DeleteByRecordId(context *gin.Context)
}

type filemasterController struct {
	filemasterService service.FilemasterService
	jwtService        service.JWTService
}

func NewFilemasterController(_service service.FilemasterService, jwtService service.JWTService) FilemasterController {
	return &filemasterController{
		filemasterService: _service,
		jwtService:        jwtService,
	}
}

func (c *filemasterController) GetAll(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = c.filemasterService.GetAll(filter)
	response := helper.BuildResponse(true, "Ok", result)
	context.JSON(http.StatusOK, response)
}

func (c *filemasterController) Insert(context *gin.Context) {
	result := helper.Response{}
	var record models.Filemaster

	err := context.Bind(&record)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		record.EntityId = userIdentity.EntityId
		record.CreatedBy = userIdentity.UserId
		result = c.filemasterService.Insert(record)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *filemasterController) DeleteByRecordId(context *gin.Context) {
	recordId := context.Param("recordId")
	if recordId == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.filemasterService.DeleteByRecordId(recordId)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
