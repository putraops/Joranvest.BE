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

type MembershipController interface {
	GetDatatables(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type membershipController struct {
	membershipService service.MembershipService
	jwtService        service.JWTService
}

func NewMembershipController(membershipService service.MembershipService, jwtService service.JWTService) MembershipController {
	return &membershipController{
		membershipService: membershipService,
		jwtService:        jwtService,
	}
}

func (c *membershipController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.membershipService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *membershipController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.MembershipDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.Membership{}
		// newRecord.OrderDetail = []models.OrderDetail{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		// json.Unmarshal([]byte(recordDto.Detail), &newRecord.OrderDetail)
		newRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			result = c.membershipService.Insert(newRecord)
		}
		// else {
		// 	json.Unmarshal([]byte(recordDto.Detail), &recordDto.OrderDetail)
		// 	recordDto.UpdatedBy = userIdentity.UserId
		// 	result = c.orderService.Update(recordDto)
		// }

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *membershipController) GetUserSession(context *gin.Context) helper.UserSession {
	var appSession helper.AppSession = helper.NewAppSession(context)
	userSession := appSession.GetUserSession()
	return userSession
}

func (c *membershipController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.membershipService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *membershipController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.membershipService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}