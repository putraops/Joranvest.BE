package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type FundamentalAnalysisController interface {
	GetDatatables(context *gin.Context)
	GetPagination(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type fundamentalAnalysisController struct {
	fundamentalAnalysisService service.FundamentalAnalysisService
	jwtService                 service.JWTService
}

func NewFundamentalAnalysisController(fundamentalAnalysisService service.FundamentalAnalysisService, jwtService service.JWTService) FundamentalAnalysisController {
	return &fundamentalAnalysisController{
		fundamentalAnalysisService: fundamentalAnalysisService,
		jwtService:                 jwtService,
	}
}

func (c *fundamentalAnalysisController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.fundamentalAnalysisService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *fundamentalAnalysisController) GetPagination(context *gin.Context) {
	var req commons.PaginationRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.fundamentalAnalysisService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

func (c *fundamentalAnalysisController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.FundamentalAnalysisDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		fmt.Println("not error")
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.FundamentalAnalysis{}
		newRecord.FundamentalAnalysisTag = []models.FundamentalAnalysisTag{}

		var arrtempTagId []string
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		json.Unmarshal([]byte(recordDto.Tag), &arrtempTagId)
		json.Unmarshal([]byte(recordDto.Tag), &newRecord.FundamentalAnalysisTag)
		newRecord.EntityId = userIdentity.EntityId

		//-- Mapping Tag Id in array into array of struct of FundamentalAnalysisTag
		for i := 0; i < len(arrtempTagId); i++ {
			newRecord.FundamentalAnalysisTag[i].TagId = arrtempTagId[i]
		}

		// myDateString := "2018-01-20 00:00:00"
		//var temp = time.Time{}
		fmt.Println("not error")
		fmt.Println(recordDto.ResearchDate)
		fmt.Println(recordDto.ResearchDate.String())
		tempDate, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", recordDto.ResearchDate.String())
		if err != nil {
			panic(err)
		}
		newRecord.ResearchDate.Time = tempDate
		newRecord.ResearchDate.Valid = true
		// record.ResearchDate = "2021-08-20 17:31:54.911026"

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			newRecord.OwnerId = userIdentity.UserId
			result = c.fundamentalAnalysisService.Insert(newRecord)
		} else {
			newRecord.UpdatedBy = userIdentity.UserId
			result = c.fundamentalAnalysisService.Update(newRecord)
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

func (c *fundamentalAnalysisController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.fundamentalAnalysisService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *fundamentalAnalysisController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.fundamentalAnalysisService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
