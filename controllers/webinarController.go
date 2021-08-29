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

type WebinarController interface {
	GetDatatables(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type webinarController struct {
	webinarService service.WebinarService
	jwtService     service.JWTService
}

func NewWebinarController(webinarService service.WebinarService, jwtService service.JWTService) WebinarController {
	return &webinarController{
		webinarService: webinarService,
		jwtService:     jwtService,
	}
}

func (c *webinarController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.webinarService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *webinarController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.WebinarDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		fmt.Println("not error")
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.Webinar{}
		newRecord.WebinarSpeaker = []models.WebinarSpeaker{}

		var arrtempWebinarSpeakerId []string
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		json.Unmarshal([]byte(recordDto.WebinarSpeaker), &arrtempWebinarSpeakerId)
		json.Unmarshal([]byte(recordDto.WebinarSpeaker), &newRecord.WebinarSpeaker)
		newRecord.EntityId = userIdentity.EntityId

		fmt.Println(newRecord.WebinarCategoryId)

		//-- Mapping Id in array into array of struct of WebinarSpeaker
		for i := 0; i < len(arrtempWebinarSpeakerId); i++ {
			newRecord.WebinarSpeaker[i].SpeakerId = arrtempWebinarSpeakerId[i]
		}

		// myDateString := "2018-01-20 00:00:00"
		//var temp = time.Time{}
		//-- Mapping WebinarStartDate
		firstStartDate, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", recordDto.WebinarFirstStartDate.String())
		if err != nil {
			panic(err)
		}
		newRecord.WebinarFirstStartDate.Time = firstStartDate
		newRecord.WebinarFirstStartDate.Valid = true

		firstEndDate, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", recordDto.WebinarFirstEndDate.String())
		if err != nil {
			panic(err)
		}
		newRecord.WebinarFirstEndDate.Time = firstEndDate
		newRecord.WebinarFirstEndDate.Valid = true

		//-- Mapping WebinarEndDate
		if recordDto.WebinarLastEndDate.String() != "0001-01-01 00:00:00 +0000 UTC" {
			lastStartDate, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", recordDto.WebinarLastStartDate.String())
			if err != nil {
				panic(err)
			}
			newRecord.WebinarLastStartDate.Time = lastStartDate
			newRecord.WebinarLastStartDate.Valid = true

			lastEndDate, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", recordDto.WebinarLastEndDate.String())
			if err != nil {
				panic(err)
			}
			newRecord.WebinarLastEndDate.Time = lastEndDate
			newRecord.WebinarLastEndDate.Valid = true
		}

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			newRecord.OwnerId = userIdentity.UserId
			result = c.webinarService.Insert(newRecord)
		} else {
			newRecord.UpdatedBy = userIdentity.UserId
			result = c.webinarService.Update(newRecord)
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

func (c *webinarController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *webinarController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.webinarService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
