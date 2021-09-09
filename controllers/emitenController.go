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

type EmitenController interface {
	GetDatatables(context *gin.Context)
	Lookup(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type emitenController struct {
	emitenService service.EmitenService
	jwtService    service.JWTService
}

func NewEmitenController(emitenService service.EmitenService, jwtService service.JWTService) EmitenController {
	return &emitenController{
		emitenService: emitenService,
		jwtService:    jwtService,
	}
}

func (c *emitenController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.emitenService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *emitenController) Lookup(context *gin.Context) {
	var request helper.Select2Request
	qry := context.Request.URL.Query()

	if _, found := qry["q"]; found {
		request.Q = fmt.Sprint(qry["q"][0])
	}
	request.Field = helper.StringifyToArray(fmt.Sprint(qry["field"]))

	var result = c.emitenService.Lookup(request)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *emitenController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.EmitenDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.Emiten{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			result = c.emitenService.Insert(newRecord)
		} else {
			recordDto.UpdatedBy = userIdentity.UserId
			result = c.emitenService.Update(newRecord)
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

// @Tags Emiten
// @Summary Get Emiten by Id
// @Description Get Emiten By Id
// @Param id path string true "User Id"
// @Success 200 {object} models.Emiten
// @Failure 400,404 {object} object
// @Router /emiten/getById/{id} [get]
func (c *emitenController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.emitenService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

// @Tags Emiten
// @Summary Delete Emiten By Id
// @Param id path string true "id"
// @Router /emiten/deleteById [delete]
// @Success 200 {object} helper.Response
// @Failure 400,404 {object} object
// @Router /emiten/getById/{id} [get]
func (c *emitenController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.emitenService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
