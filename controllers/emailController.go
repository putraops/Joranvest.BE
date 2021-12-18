package controllers

import (
	"fmt"
	"joranvest/helper"
	"joranvest/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailController interface {
	SendEmailVerification(context *gin.Context)
}

type emailController struct {
	emailService service.EmailService
	jwtService   service.JWTService
}

func NewEmailController(emailService service.EmailService, jwtService service.JWTService) EmailController {
	return &emailController{
		emailService: emailService,
		jwtService:   jwtService,
	}
}

func (c *emailController) SendEmailVerification(context *gin.Context) {
	to := []string{"putraops@gmail.com"}

	result := c.emailService.SendEmailVerification(to, "id-test")
	if result.Status {
		response := helper.BuildResponse(true, "OK", result.Data)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse(result.Message, fmt.Sprintf("%v", result.Errors), helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
