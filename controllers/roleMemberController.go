package controllers

import (
	"net/http"

	"joranvest/commons"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleMemberController interface {
	GetPagination(context *gin.Context)
	Save(context *gin.Context)
	// Submit(context *gin.Context)
	GetById(context *gin.Context)
	GetViewById(context *gin.Context)
	GetByRoleId(context *gin.Context)
	DeleteById(context *gin.Context)
}

type roleMemberController struct {
	roleMemberService    service.RoleMemberService
	roleMemberRepository repository.RoleMemberRepository
	jwtService           service.JWTService
	db                   *gorm.DB
}

func NewRoleMemberController(db *gorm.DB, jwtService service.JWTService) RoleMemberController {
	return &roleMemberController{
		db:                   db,
		jwtService:           jwtService,
		roleMemberService:    service.NewRoleMemberService(db, jwtService),
		roleMemberRepository: repository.NewRoleMemberRepository(db),
	}
}

// @Tags         RoleMember
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body commons.Pagination2ndRequest true "body"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role_member/getPagination [post]
func (c roleMemberController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	var result = c.roleMemberService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

// @Tags         RoleMember
// @Security 	 BearerAuth
// @Summary 	 Get All
// @Accept       json
// @Produce      json
// @Param        role_id query string false "RoleId"
// @Param        application_user_id query string false "ApplicationUserId"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role_member/getAll [get]
func (r roleMemberController) GetAll(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = r.roleMemberService.GetAll(filter)
	context.JSON(http.StatusOK, helper.StandartResult(true, "Ok", result))
}

// @Tags         RoleMember
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        roleId 		path 		string 	true 	"roleId"
// @Success      200 			{object} 	object
// @Failure      400,404,500	{object}  	object
// @Router       /role_member/getByRoleId/{roleId} [get]
func (r roleMemberController) GetByRoleId(context *gin.Context) {
	roleId := context.Param("roleId")
	if roleId == "" {
		context.JSON(http.StatusBadRequest, helper.StandartResult(false, "Failed to get roleId", nil))
		return
	}

	var result = r.roleMemberService.GetByRoleId(roleId)
	context.JSON(http.StatusOK, helper.StandartResult(true, "Ok", result))
}

// @Tags         RoleMember
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body models.RoleMember true "request"
// @Success      200  {object}  object
// @Failure      400,404,500  {object}  object
// @Router       /role_member/save [post]
func (r roleMemberController) Save(c *gin.Context) {
	var request models.RoleMember
	err := c.Bind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.StandartResult(false, err.Error(), nil))
		return
	}

	result := r.roleMemberService.Save(request, c)
	c.JSON(http.StatusOK, result)
	return
}

// @Tags         RoleMember
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role_member/getById/{id} [get]
func (c roleMemberController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	result := c.roleMemberService.GetById(id)
	context.JSON(http.StatusOK, result)
}

// @Tags         RoleMember
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role_member/getViewById/{id} [get]
func (c roleMemberController) GetViewById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	result := c.roleMemberService.GetViewById(id)
	context.JSON(http.StatusOK, result)
}

// @Tags         RoleMember
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role_member/deleteById/{id} [delete]
func (c roleMemberController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	result := c.roleMemberService.DeleteById(id)
	context.JSON(http.StatusOK, result)
}
