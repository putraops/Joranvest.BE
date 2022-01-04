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

type OrganizationController interface {
	Lookup(context *gin.Context)
	GetDatatables(context *gin.Context)
	GetPagination(context *gin.Context)
	GetById(context *gin.Context)
	GetViewById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type organizationController struct {
	organizationService service.OrganizationService
	jwtService          service.JWTService
}

func NewOrganizationController(organizationService service.OrganizationService, jwtService service.JWTService) OrganizationController {
	return &organizationController{
		organizationService: organizationService,
		jwtService:          jwtService,
	}
}

func (c *organizationController) Lookup(context *gin.Context) {
	var request helper.ReactSelectRequest
	qry := context.Request.URL.Query()

	if _, found := qry["q"]; found {
		request.Q = fmt.Sprint(qry["q"][0])
	}
	request.Field = helper.StringifyToArray(fmt.Sprint(qry["field"]))

	var result = c.organizationService.Lookup(request)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *organizationController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.organizationService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *organizationController) GetPagination(context *gin.Context) {
	commons.Logger()

	var req commons.PaginationRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		log.Error(errDTO.Error())
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.organizationService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

func (c *organizationController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.OrganizationDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.Organization{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			result = c.organizationService.Insert(newRecord)
		} else {
			newRecord.UpdatedBy = userIdentity.UserId
			result = c.organizationService.Update(newRecord)
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

func (c *organizationController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.organizationService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *organizationController) GetViewById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.organizationService.GetViewById(id)
	if result.Status {
		response := helper.BuildResponse(result.Status, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildResponse(result.Status, result.Message, helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *organizationController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.organizationService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
