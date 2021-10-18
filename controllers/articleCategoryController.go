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

type ArticleCategoryController interface {
	Lookup(context *gin.Context)
	GetDatatables(context *gin.Context)
	GetTree(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
	Save(context *gin.Context)
}

type articleCategoryController struct {
	articleCategoryService service.ArticleCategoryService
	jwtService             service.JWTService
}

func NewArticleCategoryController(articleCategoryService service.ArticleCategoryService, jwtService service.JWTService) ArticleCategoryController {
	return &articleCategoryController{
		articleCategoryService: articleCategoryService,
		jwtService:             jwtService,
	}
}

func (c *articleCategoryController) Lookup(context *gin.Context) {
	var request helper.ReactSelectRequest
	qry := context.Request.URL.Query()

	if _, found := qry["q"]; found {
		request.Q = fmt.Sprint(qry["q"][0])
	}
	request.Field = helper.StringifyToArray(fmt.Sprint(qry["field"]))

	var result = c.articleCategoryService.Lookup(request)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *articleCategoryController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.articleCategoryService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *articleCategoryController) GetTree(context *gin.Context) {
	var result = c.articleCategoryService.GetTree()
	context.JSON(http.StatusOK, result)
}

func (c *articleCategoryController) Save(context *gin.Context) {
	result := helper.Response{}
	var recordDto dto.ArticleCategoryDto

	errDTO := context.Bind(&recordDto)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userIdentity := c.jwtService.GetUserByToken(authHeader)

		var newRecord = models.ArticleCategory{}
		smapping.FillStruct(&newRecord, smapping.MapFields(&recordDto))
		newRecord.EntityId = userIdentity.EntityId

		if recordDto.Id == "" {
			newRecord.CreatedBy = userIdentity.UserId
			result = c.articleCategoryService.Insert(newRecord)
		} else {
			newRecord.UpdatedBy = userIdentity.UserId
			result = c.articleCategoryService.Update(newRecord)
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

func (c *articleCategoryController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.articleCategoryService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *articleCategoryController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.articleCategoryService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
