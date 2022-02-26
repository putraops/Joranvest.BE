package controllers

import (
	"fmt"
	"net/http"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleController interface {
	Save(context *gin.Context)
	GetPagination(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
}

type roleController struct {
	roleService service.RoleService
	jwtService  service.JWTService
	db          *gorm.DB
}

func NewRoleController(db *gorm.DB, jwtService service.JWTService) RoleController {
	return &roleController{
		db:          db,
		jwtService:  jwtService,
		roleService: service.NewRoleService(db, jwtService),
	}
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

// @Tags         Role
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.RoleDto true "dto"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /role/save [post]
func (r roleController) Save(c *gin.Context) {
	var result helper.Result
	var dto dto.RoleDto
	r.db = c.MustGet("db_trx").(*gorm.DB)
	dto.Context = c

	errDto := c.Bind(&dto)
	if errDto != nil {
		res := helper.StandartResult(false, errDto.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if dto.Id == "" {
		result = r.roleService.OpenTransaction(r.db).Insert(dto)
		if !result.Status {
			c.JSON(http.StatusInternalServerError, helper.StandartResult(result.Status, result.Message, result.Data))
			r.db.Rollback()
			return
		}
	} else {
		fmt.Println("Do Update")
	}

	if err := r.db.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, helper.StandartResult(result.Status, result.Message, result.Data))
		return
	}

	c.JSON(http.StatusOK, helper.StandartResult(result.Status, result.Message, result.Data))
	return
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
