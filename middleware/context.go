package middleware

import (
	"github.com/gin-gonic/gin"
)

func ContextMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ctx_token", token)
		c.Next()
	}
}
