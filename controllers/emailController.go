package controllers

import (
	"joranvest/service"
)

type EmailController interface {
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
