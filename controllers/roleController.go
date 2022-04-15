package controllers

import (
	"net/http"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/repository"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleController interface {
	GetPagination(context *gin.Context)
	GetAll(context *gin.Context)
	Save(context *gin.Context)
	GetById(context *gin.Context)
	GetViewById(context *gin.Context)
	DeleteById(context *gin.Context)

	//-- Notification
	SetDashboardAccess(c *gin.Context)
	SetFullAccess(c *gin.Context)
	SetPaymentNotification(c *gin.Context)
	GetNotificationByRoleId(context *gin.Context)
}

type roleController struct {
	roleService    service.RoleService
	roleRepository repository.RoleRepository
	jwtService     service.JWTService
	db             *gorm.DB
}

func NewRoleController(db *gorm.DB, jwtService service.JWTService) RoleController {
	return &roleController{
		db:             db,
		jwtService:     jwtService,
		roleService:    service.NewRoleService(db, jwtService),
		roleRepository: repository.NewRoleRepository(db),
	}
}

// @Tags         Role
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body commons.Pagination2ndRequest true "body"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role/getPagination [post]
func (c roleController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	var result = c.roleService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

// @Tags         Role
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body models.Role true "request"
// @Success      200  {object}  object
// @Failure      400,404,500  {object}  object
// @Router       /role/save [post]
func (r roleController) Save(c *gin.Context) {
	var request models.Role
	err := c.Bind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.StandartResult(false, err.Error(), nil))
		return
	}

	result := r.roleService.Save(request, c)
	c.JSON(http.StatusOK, result)
	return
}

// @Tags         Role
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.PaymentNotificationDto true "request"
// @Success      200  {object}  object
// @Failure      400,404,500  {object}  object
// @Router       /role/set/paymentNotification [post]
func (r roleController) SetPaymentNotification(c *gin.Context) {
	var request dto.PaymentNotificationDto

	err := c.Bind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.StandartResult(false, err.Error(), nil))
		return
	}

	result := r.roleService.SetPaymentNotification(request, c)
	c.JSON(http.StatusOK, result)
	return
}

// @Tags         Role
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.RoleAccessDto true "request"
// @Success      200  {object}  object
// @Failure      400,404,500  {object}  object
// @Router       /role/set/dashboardAccess [post]
func (r roleController) SetDashboardAccess(c *gin.Context) {
	var request dto.RoleAccessDto

	err := c.Bind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.StandartResult(false, err.Error(), nil))
		return
	}

	result := r.roleService.SetDashboardAccess(request, c)
	c.JSON(http.StatusOK, result)
	return
}

// @Tags         Role
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.RoleAccessDto true "request"
// @Success      200  {object}  object
// @Failure      400,404,500  {object}  object
// @Router       /role/set/fullAccess [post]
func (r roleController) SetFullAccess(c *gin.Context) {
	var request dto.RoleAccessDto

	err := c.Bind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.StandartResult(false, err.Error(), nil))
		return
	}

	result := r.roleService.SetFullAccess(request, c)
	c.JSON(http.StatusOK, result)
	return
}

// @Tags         Role
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        roleId path string true "roleId"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role/notification/getByRoleId/{roleId} [get]
func (c roleController) GetNotificationByRoleId(context *gin.Context) {
	roleId := context.Param("roleId")
	if roleId == "" {
		response := helper.BuildErrorResponse("Failed to get roleId", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	result := c.roleService.GetNotificationConfigurationByRoleId(roleId)
	context.JSON(http.StatusOK, result)
}

// @Tags         Role
// @Security 	 BearerAuth
// @Summary 	 Get All
// @Accept       json
// @Produce      json
// @Param        name query string false "name"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role/getAll [get]
func (r roleController) GetAll(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = r.roleService.GetAll(filter)
	context.JSON(http.StatusOK, helper.StandartResult(true, "Ok", result))
}

// @Tags         Role
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role/getById/{id} [get]
func (c roleController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	result := c.roleService.GetById(id)
	context.JSON(http.StatusOK, result)
}

// @Tags         Role
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role/getViewById/{id} [get]
func (c roleController) GetViewById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	result := c.roleService.GetViewById(id)
	context.JSON(http.StatusOK, result)
}

// @Tags         Role
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "id"
// @Success      200 {object} object
// @Failure      400,404,500  {object}  object
// @Router       /role/deleteById/{id} [delete]
func (c roleController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}

	result := c.roleService.DeleteById(id)
	context.JSON(http.StatusOK, result)
}
