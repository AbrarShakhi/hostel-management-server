package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) VerifyUser(c *gin.Context) {
	tokenString, err := c.Cookie("user_auth")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "cookie error",
		})
		return
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header)
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "cookie error",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "cookie error",
		})
		return
	}

	if exp, ok := claims["exp"].(float64); ok && float64(time.Now().Unix()) > exp {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "cookie error",
		})
		return
	}
	if sub, ok := claims["sub"].(string); ok && sub == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "cookie error",
		})
		return
	}

	c.Set("user_id", claims["sub"])
	c.Next()
}
