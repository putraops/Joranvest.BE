package controllers

import (
	"joranvest/repository"
	"joranvest/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailController interface {
	TestEmail(context *gin.Context)
}

type emailController struct {
	emailService      service.EmailService
	paymentRepository repository.PaymentRepository
	jwtService        service.JWTService
	DB                *gorm.DB
}

func NewEmailController(db *gorm.DB, jwtService service.JWTService) EmailController {
	return &emailController{
		DB:                db,
		emailService:      service.NewEmailService(db),
		paymentRepository: repository.NewPaymentRepository(db),
		jwtService:        jwtService,
	}
}

func (c emailController) TestEmail(context *gin.Context) {
	//res := c.paymentRepository.GetViewById("11cfda96-fed7-42e1-bddd-367b369bd6ac")

	// temp := res.Data.(entity_view_models.EntityPaymentView)

	// response := c.emailService.NewPayment(temp)
	context.JSON(http.StatusOK, nil)
	return
}
