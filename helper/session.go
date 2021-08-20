package helper

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserSession struct {
	UserId   string
	EntityId string
	TenantId string
	Token    string
}

type AppSession interface {
	GetUserSession() UserSession
}

type appSession struct {
	c *gin.Context
}

func NewAppSession(ctx *gin.Context) AppSession {
	return &appSession{
		c: ctx,
	}
}

func (u *appSession) GetUserSession() UserSession {
	var userSession UserSession
	session := sessions.Default(u.c)
	userId := session.Get("UserId")
	entityid := session.Get("EntityId")
	userSession.UserId = fmt.Sprintf("%v", userId)
	userSession.EntityId = fmt.Sprintf("%v", entityid)
	return userSession
}
