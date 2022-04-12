package helper

import (
	"github.com/gin-gonic/gin"
)

type emailNotificationGroup struct {
	context                   *gin.Context
	PaymentNotificationAccess *PaymentNotificationAccess
}

type PaymentNotificationAccess struct {
	HasAccess bool
}

type EmailNotificationGroup interface {
	GetPaymentNotification() bool
}

func NewNotificationHelper(_context *gin.Context) EmailNotificationGroup {
	return &emailNotificationGroup{
		PaymentNotificationAccess: &PaymentNotificationAccess{
			HasAccess: true,
		},
		context: _context,
	}
}

func (c emailNotificationGroup) GetPaymentNotification() bool {
	return true
}
