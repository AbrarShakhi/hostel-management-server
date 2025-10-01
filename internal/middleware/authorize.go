package middleware

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *middleware) VerifyUser(c *gin.Context) {
	tokenString, err := c.Cookie("user_auth")
	if err != nil {
		m.unauthorizedCookieError(c)
		return
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header)
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		m.unauthorizedCookieError(c)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		m.unauthorizedCookieError(c)
		return
	}

	if exp, ok := claims["exp"].(float64); ok && float64(time.Now().Unix()) > exp {
		m.unauthorizedCookieError(c)
		return
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		m.unauthorizedCookieError(c)
		return
	}

	userId, err := strconv.Atoi(sub)
	if err != nil {
		m.unauthorizedCookieError(c)
		return
	}
	c.Set("user_id", userId)
	c.Next()
}
