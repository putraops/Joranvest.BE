package controllers

import (
	"net/http"

	"joranvest/helper"
	"joranvest/service"

	"github.com/gin-gonic/gin"
)

type ArticleTagController interface {
	GetById(context *gin.Context)
	GetAll(context *gin.Context)
}

type articleTagController struct {
	articleTagService service.ArticleTagService
	jwtService        service.JWTService
}

func NewArticleTagController(_service service.ArticleTagService, jwtService service.JWTService) ArticleTagController {
	return &articleTagController{
		articleTagService: _service,
		jwtService:        jwtService,
	}
}

func (c *articleTagController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.articleTagService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *articleTagController) GetAll(context *gin.Context) {
	qry := context.Request.URL.Query()
	filter := make(map[string]interface{})

	for k, v := range qry {
		filter[k] = v
	}

	var result = c.articleTagService.GetAll(filter)
	response := helper.BuildResponse(true, "Ok", result)
	context.JSON(http.StatusOK, response)
}
