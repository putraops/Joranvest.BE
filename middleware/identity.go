package middleware

import (
	"fmt"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserSession struct {
	UserId   string
	EntityId string
	TenantId string
}

//AuthorizeJWT validates the token user given, return 401 if not valid
func SetIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		println("From Middleware")
		println("==========================")
		userLoginName := session.Get("userLoginName")
		userId := session.Get("UserId")
		entityId := session.Get("EntityId")
		tenantId := session.Get("TenantId")
		isAdmin := session.Get("IsAdmin")
		IsSuperAdmin := session.Get("IsSuperAdmin")

		os.Setenv("USER_ID", fmt.Sprintf("%v", userId))
		os.Setenv("USER_LOGINNAME", fmt.Sprintf("%v", userLoginName))
		os.Setenv("ENTITY_ID", fmt.Sprintf("%v", entityId))
		os.Setenv("TENANT_ID", fmt.Sprintf("%v", tenantId))
		os.Setenv("IS_ADMIN", fmt.Sprintf("%v", isAdmin))
		os.Setenv("IS_SUPERADMIN", fmt.Sprintf("%v", IsSuperAdmin))
	}
}
