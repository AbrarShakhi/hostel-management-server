package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type middleware struct {
}

func NewMiddleware() *middleware {
	return &middleware{}
}

func (m *middleware) unauthorizedCookieError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"msg": "cookie error",
	})
}
