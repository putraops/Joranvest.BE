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
	log "github.com/sirupsen/logrus"
)

type MembershipController interface {
	Lookup(context *gin.Context)
	GetDatatables(context *gin.Context)
	GetPagination(context *gin.Context)
	GetAll(context *gin.Context)
	GetById(context *gin.Context)
	GetViewById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
	SetRecommendation(context *gin.Context)
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

// @Tags         Membership
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body helper.ReactSelectRequest true "body"
// @Param        q query string false "q"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /membership/lookup [post]
func (c *membershipController) Lookup(context *gin.Context) {
	var request helper.ReactSelectRequest

	errDTO := context.Bind(&request)
	if errDTO != nil {
		res := helper.StandartResult(false, errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	var result = c.membershipService.Lookup(request)
	response := helper.StandartResult(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
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

func (c *membershipController) GetPagination(context *gin.Context) {
	commons.Logger()

	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		log.Error(errDTO.Error())
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.membershipService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

func (c *membershipController) GetAll(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = c.membershipService.GetAll(filter)
	response := helper.BuildResponse(true, "Ok", result)
	context.JSON(http.StatusOK, response)
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
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			result = c.membershipService.Insert(newRecord)
		} else {
			newRecord.UpdatedBy = userIdentity.UserId
			result = c.membershipService.Update(newRecord)
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
		return
	}

	result := c.membershipService.GetById(id)
	context.JSON(http.StatusOK, result)
}

func (c *membershipController) GetViewById(context *gin.Context) {
	commons.Logger()
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		log.Error("membershipController: Failed to Get Id")
		context.JSON(http.StatusBadRequest, response)
		return
	}

	result := c.membershipService.GetViewById(id)
	context.JSON(http.StatusOK, result)
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
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *membershipController) SetRecommendation(context *gin.Context) {
	var recordRecommendationDto dto.MembershipRecommendationDto
	errDto := context.Bind(&recordRecommendationDto)
	if errDto != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDto.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}

	if recordRecommendationDto.Id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.membershipService.SetRecommendationById(recordRecommendationDto.Id, recordRecommendationDto.IsChecked)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
