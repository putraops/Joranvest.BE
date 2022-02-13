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

type RoleController interface {
	Lookup(context *gin.Context)
	GetPagination(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type roleController struct {
	roleService service.RoleService
	jwtService  service.JWTService
}

func NewRoleController(roleService service.RoleService, jwtService service.JWTService) RoleController {
	return &roleController{
		roleService: roleService,
		jwtService:  jwtService,
	}
}

func (c *roleController) Lookup(context *gin.Context) {
	var request helper.Select2Request
	qry := context.Request.URL.Query()

	if _, found := qry["q"]; found {
		request.Q = fmt.Sprint(qry["q"][0])
	}
	request.Field = helper.StringifyToArray(fmt.Sprint(qry["field"]))

	var result = c.roleService.Lookup(request)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *roleController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.roleService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

func (c *roleController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.RoleDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.Role{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			result = c.roleService.Insert(newRecord)
		} else {
			newRecord.UpdatedBy = userIdentity.UserId
			result = c.roleService.Update(newRecord)
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

func (c *roleController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.roleService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *roleController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.roleService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
