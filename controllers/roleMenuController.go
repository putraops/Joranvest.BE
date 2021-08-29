package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"joranvest/dto"
	"joranvest/helper"
	"joranvest/service"

	"github.com/gin-gonic/gin"
)

type RoleMenuController interface {
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	DeleteByRoleAndMenuId(context *gin.Context)
	Save(context *gin.Context)
}

type roleMenuController struct {
	roleMenuService service.RoleMenuService
	jwtService      service.JWTService
}

func NewRoleMenuController(roleMenuService service.RoleMenuService, jwtService service.JWTService) RoleMenuController {
	return &roleMenuController{
		roleMenuService: roleMenuService,
		jwtService:      jwtService,
	}
}

func (c *roleMenuController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.InsertRoleMenuDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		recordDto.EntityId = userIdentity.EntityId
		recordDto.CreatedBy = userIdentity.UserId
		json.Unmarshal([]byte(recordDto.Children), &recordDto.ChildrenIds)

		result = c.roleMenuService.Insert(recordDto)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *roleMenuController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.roleMenuService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *roleMenuController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.roleMenuService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *roleMenuController) DeleteByRoleAndMenuId(context *gin.Context) {
	var recordDto dto.DeleteRoleMenuDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}

	if recordDto.RoleId == "" || recordDto.ApplicationMenuId == "" {
		response := helper.BuildErrorResponse("Failed to bind data", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.roleMenuService.DeleteByRoleAndMenuId(recordDto.RoleId, recordDto.ApplicationMenuId, recordDto.IsParent)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
