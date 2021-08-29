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

type ApplicationMenuController interface {
	Lookup(context *gin.Context)
	GetDatatables(context *gin.Context)
	GetTree(context *gin.Context)
	GetTreeByRoleId(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type applicationMenuController struct {
	applicationMenuService service.ApplicationMenuService
	jwtService             service.JWTService
}

func NewApplicationMenuController(applicationMenuService service.ApplicationMenuService, jwtService service.JWTService) ApplicationMenuController {
	return &applicationMenuController{
		applicationMenuService: applicationMenuService,
		jwtService:             jwtService,
	}
}

func (c *applicationMenuController) Lookup(context *gin.Context) {
	var request helper.Select2Request
	qry := context.Request.URL.Query()

	if _, found := qry["q"]; found {
		request.Q = fmt.Sprint(qry["q"][0])
	}
	request.Field = helper.StringifyToArray(fmt.Sprint(qry["field"]))

	var result = c.applicationMenuService.Lookup(request)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *applicationMenuController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.applicationMenuService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *applicationMenuController) GetTree(context *gin.Context) {
	var result = c.applicationMenuService.GetTree()
	context.JSON(http.StatusOK, result)
}

func (c *applicationMenuController) GetTreeByRoleId(context *gin.Context) {
	roleId := context.Param("roleId")
	var result = c.applicationMenuService.GetTreeByRoleId(roleId)
	context.JSON(http.StatusOK, result)
}

func (c *applicationMenuController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.ApplicationMenuDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.ApplicationMenu{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			result = c.applicationMenuService.Insert(newRecord)
		} else {
			newRecord.UpdatedBy = userIdentity.UserId
			result = c.applicationMenuService.Update(newRecord)
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

func (c *applicationMenuController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.applicationMenuService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *applicationMenuController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.applicationMenuService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
