package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *middleware) IdentifyOtpUser(c *gin.Context) {
	identifier := c.Query("identifier")
	reason, err := strconv.Atoi(c.Query("reason"))
	if identifier == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing 'identifier' query parameter"})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "'reason' must be a valid integer"})
		return
	}
	if reason != 1 && reason != 2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid 'reason'"})
		return
	}

	if strings.Contains(identifier, "@") {
		c.Set("email", identifier)
	} else {
		c.Set("phone", identifier)
	}
	c.Set("reason", reason)
	c.Next()
}
