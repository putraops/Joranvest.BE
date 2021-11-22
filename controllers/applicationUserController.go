package controllers

import (
	"fmt"
	"net/http"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ApplicationUserController interface {
	GetDatatables(context *gin.Context)
	Lookup(context *gin.Context)
	Update(context *gin.Context)
	ChangePassword(context *gin.Context)
	Profile(context *gin.Context)
	GetAll(context *gin.Context)
	GetById(context *gin.Context)
	GetViewById(context *gin.Context)
	DeleteById(context *gin.Context)
}

type applicationUserController struct {
	applicationUserService service.ApplicationUserService
	jwtService             service.JWTService
}

func NewApplicationUserController(applicationUserService service.ApplicationUserService, jwtService service.JWTService) ApplicationUserController {
	return &applicationUserController{
		applicationUserService: applicationUserService,
		jwtService:             jwtService,
	}
}

func (c *applicationUserController) GetDatatables(context *gin.Context) {
	var dt commons.DataTableRequest
	errDTO := context.Bind(&dt)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.applicationUserService.GetDatatables(dt)
	context.JSON(http.StatusOK, result)
}

func (c *applicationUserController) Lookup(context *gin.Context) {
	var request helper.ReactSelectRequest
	qry := context.Request.URL.Query()

	if _, found := qry["q"]; found {
		request.Q = fmt.Sprint(qry["q"][0])
	}
	request.Field = helper.StringifyToArray(fmt.Sprint(qry["field"]))

	var result = c.applicationUserService.Lookup(request)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *applicationUserController) GetAll(context *gin.Context) {
	var users = c.applicationUserService.GetAll()
	res := helper.BuildResponse(true, "Ok", users)
	context.JSON(http.StatusOK, res)
}

func (c *applicationUserController) Update(context *gin.Context) {
	var applicationUserUpdateDto dto.ApplicationUserUpdateDto
	errDTO := context.ShouldBind(&applicationUserUpdateDto)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	applicationUserUpdateDto.Id = id
	u := c.applicationUserService.Update(applicationUserUpdateDto)
	res := helper.BuildResponse(true, "Ok!", u)
	context.JSON(http.StatusOK, res)
}

func (c *applicationUserController) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user := c.applicationUserService.UserProfile(id)

	res := helper.BuildResponse(true, "Ok!", user)
	context.JSON(http.StatusOK, res)
}

func (c *applicationUserController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result helper.Response
	result = c.applicationUserService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *applicationUserController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.applicationUserService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *applicationUserController) GetViewById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.applicationUserService.GetViewById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *applicationUserController) ChangePassword(context *gin.Context) {
	var loginDto dto.LoginDto
	err := context.ShouldBind(&loginDto)
	fmt.Println(loginDto)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to request login", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	result := c.applicationUserService.ChangePassword(loginDto.Username, loginDto.Email, loginDto.Password)
	if result.Status {
		if v, ok := (result.Data).(models.ApplicationUser); ok {
			generatedToken := c.jwtService.GenerateToken(v.Id, v.EntityId)
			v.Token = generatedToken

			response := helper.BuildResponse(true, "Ok!", v)
			ctx.JSON(http.StatusOK, response)
			return
		}
	} else {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	}
}
