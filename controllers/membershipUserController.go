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

type MembershipUserController interface {
	GetDatatables(context *gin.Context)
	GetPagination(context *gin.Context)
	GetAll(context *gin.Context)
	GetById(context *gin.Context)
	GetByUserLogin(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type membershipUserController struct {
	membershipUserService service.MembershipUserService
	jwtService            service.JWTService
}

func NewMembershipUserController(membershipUserService service.MembershipUserService, jwtService service.JWTService) MembershipUserController {
	return &membershipUserController{
		membershipUserService: membershipUserService,
		jwtService:            jwtService,
	}
}

func (c *membershipUserController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.membershipUserService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *membershipUserController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.membershipUserService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

// @Tags MembershipUser
// @Summary Get All
// @Router /membershipUser/getAll [post]
// @Success 200 {obsject} object
// @Failure 400,404 {object} object
func (c *membershipUserController) GetAll(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = c.membershipUserService.GetAll(filter)
	response := helper.BuildResponse(true, "Ok", result)
	context.JSON(http.StatusOK, response)
}

func (c *membershipUserController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.MembershipUserDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		fmt.Println("not error")
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newMembershipUserRecord = models.MembershipUser{}
		smapping.FillStruct(&newMembershipUserRecord, smapping.MapFields(&recordDto))

		var newPaymentRecord = models.Payment{}
		smapping.FillStruct(&newPaymentRecord, smapping.MapFields(&recordDto))

		newMembershipUserRecord.EntityId = userIdentity.EntityId
		newPaymentRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newMembershipUserRecord.CreatedBy = userIdentity.UserId
			newMembershipUserRecord.OwnerId = userIdentity.UserId

			newPaymentRecord.CreatedBy = userIdentity.UserId
			newPaymentRecord.OwnerId = userIdentity.UserId

			result = c.membershipUserService.Insert(newMembershipUserRecord, newPaymentRecord)
		} else {
			newMembershipUserRecord.UpdatedBy = userIdentity.UserId
			newPaymentRecord.UpdatedBy = userIdentity.UserId
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

func (c *membershipUserController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
		return
	}
	result := c.membershipUserService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *membershipUserController) GetByUserLogin(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	userIdentity := c.jwtService.GetUserByToken(authHeader)

	response := c.membershipUserService.GetByUserId(userIdentity.UserId)
	context.JSON(http.StatusOK, response)
}

func (c *membershipUserController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.membershipUserService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
