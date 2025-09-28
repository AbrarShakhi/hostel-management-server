package controller

import (
	"database/sql"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/abrarshakhi/hostel-management-server/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

func (h *Controller) UserLogin(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request payload"})
		return
	}

	var user *model.User_
	var err error
	if strings.Contains(req.Identifier, "@") {
		user, err = model.FindByEmail(h.db, req.Identifier)
	} else {
		user, err = model.FindByPhone(h.db, req.Identifier)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid user email, phone, password."})
		return
	}

	if !user.Password_.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Your account is not active yet. Active it first."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password_.String), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid user email, phone, password."})
		return
	}

	user.LastLogin = sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	if user.Update(h.db) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update login time"})
		return
	}

	claims := jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to generate token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("user_auth", signedToken, 3600*24*30, "", "", false, false)
	c.JSON(http.StatusOK, gin.H{
		"user_id": user.Id,
		"phone":   user.Phone,
		"email":   user.Email,
	})
}
