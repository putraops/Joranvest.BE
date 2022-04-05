package controllers

import (
	"joranvest/service"

	"gorm.io/gorm"
)

type EmailController interface {
}

type emailController struct {
	emailService service.EmailService
	jwtService   service.JWTService
	DB           *gorm.DB
}

func NewEmailController(db *gorm.DB, jwtService service.JWTService) EmailController {
	return &emailController{
		DB:           db,
		emailService: service.NewEmailService(db),
		jwtService:   jwtService,
	}
}
