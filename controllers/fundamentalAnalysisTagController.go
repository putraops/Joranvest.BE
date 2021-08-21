package controllers

import (
	"net/http"

	"joranvest/helper"
	"joranvest/service"

	"github.com/gin-gonic/gin"
)

type FundamentalAnalysisTagController interface {
	GetById(context *gin.Context)
	GetAll(context *gin.Context)
}

type fundamentalAnalysisTagController struct {
	fundamentalAnalysisTagService service.FundamentalAnalysisTagService
	jwtService                    service.JWTService
}

func NewFundamentalAnalysisTagController(_service service.FundamentalAnalysisTagService, jwtService service.JWTService) FundamentalAnalysisTagController {
	return &fundamentalAnalysisTagController{
		fundamentalAnalysisTagService: _service,
		jwtService:                    jwtService,
	}
}

func (c *fundamentalAnalysisTagController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.fundamentalAnalysisTagService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *fundamentalAnalysisTagController) GetAll(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = c.fundamentalAnalysisTagService.GetAll(filter)
	response := helper.BuildResponse(true, "Ok", result)
	context.JSON(http.StatusOK, response)
}
