package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

func (h *Handlers) UserLogin(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var query string
	if strings.Contains(req.Identifier, "@") {
		query = "SELECT id, email, phone, password_, first_name, last_name, last_login FROM user_ WHERE email = $1"
	} else {
		query = "SELECT id, email, phone, password_, first_name, last_name, last_login FROM user_ WHERE phone = $1"
	}

	row, err := h.db.Query(query, req.Identifier)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var (
		id         string
		email      string
		phone      string
		hashedPass string
		first_name string
		last_name  string
		last_login time.Time
	)
	row.Scan(&id, &email, &phone, &hashedPass, &first_name, &last_name, &last_login)

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	_, err = h.db.Exec(`UPDATE user_ SET last_login = $1 WHERE id = $2`, time.Now(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update login time"})
		return
	}

	claims := jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("user_auth", signedToken, 3600*24*30, "", "", false, false)
	c.JSON(http.StatusOK, gin.H{
		"user_id": id,
		"phone":   phone,
		"email":   email,
	})
}
