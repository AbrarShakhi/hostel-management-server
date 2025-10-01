package controller

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/abrarshakhi/hostel-management-server/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *controller) UserLogin(c *gin.Context) {
	req := struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request payload"})
		return
	}

	var user *model.Users
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
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid user email or phone."})
		return
	}

	if !user.HasPassword() {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Your account is not active yet. Active it first."})
		return
	}

	if !user.ComparePassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid password."})
		return
	}

	if user.HasLeft {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "You have left the hostel."})
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
	lastName := ""
	if user.LastName.Valid {
		lastName = user.LastName.String
	}

	claims := jwt.MapClaims{
		"sub": strconv.Itoa(user.UserId),
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
		"msg":        "User is authenticated",
		"user_id":    user.UserId,
		"phone":      user.Phone,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  lastName,
		"created_at": user.CreatedAt,
	})
}

func (h *controller) UserLogOut(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("user_auth", "", -1, "", "", false, false)
	c.JSON(http.StatusOK, gin.H{
		"msg": "Successfully logged out",
	})
}

func (h *controller) UserAuthCheck(c *gin.Context) {
	userIdAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	userId, ok := userIdAny.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to parse user ID"})
		return
	}

	user, err := model.FindById(h.db, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid user email or phone."})
		return
	}

	lastName := ""
	if user.LastName.Valid {
		lastName = user.LastName.String
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":        "User is authenticated",
		"user_id":    user.UserId,
		"phone":      user.Phone,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  lastName,
		"created_at": user.CreatedAt,
	})
}
