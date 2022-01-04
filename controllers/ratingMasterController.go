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

type RatingMasterController interface {
	GetPagination(context *gin.Context)
	GetAll(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type ratingMasterController struct {
	ratingMasterService service.RatingMasterService
	jwtService          service.JWTService
}

func NewRatingMasterController(ratingMasterService service.RatingMasterService, jwtService service.JWTService) RatingMasterController {
	return &ratingMasterController{
		ratingMasterService: ratingMasterService,
		jwtService:          jwtService,
	}
}

func (c *ratingMasterController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.ratingMasterService.GetPagination(req)
	if result.Status {
		response := helper.BuildResponse(true, result.Message, result.Data)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildResponse(false, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *ratingMasterController) GetAll(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = c.ratingMasterService.GetAll(filter)
	response := helper.BuildResponse(true, "Ok", result)
	context.JSON(http.StatusOK, response)
}

func (c *ratingMasterController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.RatingMasterDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.RatingMaster{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			result = c.ratingMasterService.Insert(newRecord)
		} else {
			newRecord.UpdatedBy = userIdentity.UserId
			result = c.ratingMasterService.Update(newRecord)
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

func (c *ratingMasterController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.ratingMasterService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *ratingMasterController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.ratingMasterService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
