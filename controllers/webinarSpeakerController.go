package controllers

import (
	"net/http"

	"joranvest/helper"
	"joranvest/service"

	"github.com/gin-gonic/gin"
)

type WebinarSpeakerController interface {
	GetById(context *gin.Context)
	GetAll(context *gin.Context)
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
