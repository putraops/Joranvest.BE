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

type WebinarCategoryController interface {
	Lookup(context *gin.Context)
	GetDatatables(context *gin.Context)
	GetTreeParent(context *gin.Context)
	GetTree(context *gin.Context)
	OrderTree(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type webinarCategoryController struct {
	webinarCategoryService service.WebinarCategoryService
	jwtService             service.JWTService
}

func NewWebinarCategoryController(webinarCategoryService service.WebinarCategoryService, jwtService service.JWTService) WebinarCategoryController {
	return &webinarCategoryController{
		webinarCategoryService: webinarCategoryService,
		jwtService:             jwtService,
	}
}

func (c *webinarCategoryController) Lookup(context *gin.Context) {
	var request helper.Select2Request
	qry := context.Request.URL.Query()

	if _, found := qry["q"]; found {
		request.Q = fmt.Sprint(qry["q"][0])
	}
	request.Field = helper.StringifyToArray(fmt.Sprint(qry["field"]))

	var result = c.webinarCategoryService.Lookup(request)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *webinarCategoryController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.webinarCategoryService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *webinarCategoryController) GetTreeParent(context *gin.Context) {
	var result = c.webinarCategoryService.GetTreeParent()
	context.JSON(http.StatusOK, result)
}

func (c *webinarCategoryController) GetTree(context *gin.Context) {
	var result = c.webinarCategoryService.GetTree()
	context.JSON(http.StatusOK, result)
}

func (c *webinarCategoryController) OrderTree(context *gin.Context) {
	var dto dto.OrderTreeDto
	err := context.Bind(&dto)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.webinarCategoryService.OrderTree(dto.RecordId, dto.ParentId, dto.OrderIndex)
	context.JSON(http.StatusOK, result)
}

func (c *webinarCategoryController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.WebinarCategoryDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.WebinarCategory{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			result = c.webinarCategoryService.Insert(newRecord)
		} else {
			newRecord.UpdatedBy = userIdentity.UserId
			result = c.webinarCategoryService.Update(newRecord)
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

func (c *webinarCategoryController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarCategoryService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *webinarCategoryController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.webinarCategoryService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
