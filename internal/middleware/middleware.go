package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) unauthorizedCookieError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"msg": "cookie error",
	})
}
