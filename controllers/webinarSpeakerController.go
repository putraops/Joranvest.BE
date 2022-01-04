package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type WebinarSpeakerController interface {
	Save(context *gin.Context)
	GetAll(context *gin.Context)
	GetById(context *gin.Context)
	GetSpeakersRatingByWebinarId(context *gin.Context)
	GetSpeakerReviewById(context *gin.Context)
}

type webinarSpeakerController struct {
	webinarSpeakerService service.WebinarSpeakerService
	jwtService            service.JWTService
}

func NewWebinarSpeakerController(_service service.WebinarSpeakerService, jwtService service.JWTService) WebinarSpeakerController {
	return &webinarSpeakerController{
		webinarSpeakerService: _service,
		jwtService:            jwtService,
	}
}

func (c *webinarSpeakerController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.WebinarSepakerDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = []models.WebinarSpeaker{}
		var arrtempWebinarSpeakerId []string
		json.Unmarshal([]byte(recordDto.WebinarSpeaker), &newRecord)
		json.Unmarshal([]byte(recordDto.WebinarSpeaker), &arrtempWebinarSpeakerId)

		//-- Mapping Id in array into array of struct of WebinarSpeaker
		for i := 0; i < len(arrtempWebinarSpeakerId); i++ {
			newRecord[i] = models.WebinarSpeaker{}
			newRecord[i].CreatedBy = userIdentity.UserId
			newRecord[i].OwnerId = userIdentity.UserId
			newRecord[i].WebinarId = recordDto.WebinarId
			// newRecord[i].EntityId = userIdentity.EntityId
			newRecord[i].SpeakerId = arrtempWebinarSpeakerId[i]
		}

		result = c.webinarSpeakerService.Insert(newRecord, recordDto.SpeakerType)

		if result.Status {
			response := helper.BuildResponse(true, "OK", result.Data)
			context.JSON(http.StatusOK, response)
		} else {
			response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *webinarSpeakerController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarSpeakerService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *webinarSpeakerController) GetSpeakersRatingByWebinarId(context *gin.Context) {
	commons.Logger()
	webinarId := context.Param("webinarId")
	if webinarId == "" {
		log.Error("Failed to get webinarId")
		response := helper.BuildErrorResponse("Failed to get webinarId", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarSpeakerService.GetSpeakersRatingByWebinarId(webinarId)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *webinarSpeakerController) GetSpeakerReviewById(context *gin.Context) {
	commons.Logger()
	id := context.Param("id")
	if id == "" {
		log.Error("Failed to get id")
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.webinarSpeakerService.GetSpeakerReviewById(id)
	if result.Status {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildResponse(false, result.Message, helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *webinarSpeakerController) GetAll(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = c.webinarSpeakerService.GetAll(filter)
	response := helper.BuildResponse(true, "Ok", result)
	context.JSON(http.StatusOK, response)
}
