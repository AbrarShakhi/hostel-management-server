package controller

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/abrarshakhi/hostel-management-server/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *controller) userLogin(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, gin.H{"msg": "Invalid user email or phone."})
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
		"sub": strconv.Itoa(user.UserId()),
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
		"user_id":    user.UserId(),
		"phone":      user.Phone,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  lastName,
		"created_at": user.CreatedAt,
	})
}

func (h *controller) userLogOut(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("user_auth", "", -1, "", "", false, false)
	c.JSON(http.StatusOK, gin.H{
		"msg": "Successfully logged out",
	})
}

func (h *controller) userAuthCheck(c *gin.Context) {
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

	user, err := model.FindUserById(h.db, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Invalid user email or phone."})
		return
	}

	if user.HasLeft {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "You have left the hostel."})
		return
	}

	lastName := ""
	if user.LastName.Valid {
		lastName = user.LastName.String
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":        "User is authenticated",
		"user_id":    user.UserId(),
		"phone":      user.Phone,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  lastName,
		"created_at": user.CreatedAt,
	})
}

func (h *controller) userActivateAccount(c *gin.Context) {
	identifier := c.Query("identifier")
	otpCode := c.Query("otpcode")
	req := struct {
		NewPassword string `json:"new_password" binding:"required"`
	}{}

	if identifier == "" || otpCode == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing 'identifier' or 'otpcode' query parameter"})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request payload"})
		return
	}

	if len(req.NewPassword) < 8 || len(req.NewPassword) > 40 {
		c.JSON(http.StatusLengthRequired, gin.H{"msg": "Password length must be between 8 to 40 charecter long."})
		return
	}

	var user *model.Users
	var err error
	if strings.Contains(identifier, "@") {
		user, err = model.FindByEmail(h.db, identifier)
	} else {
		user, err = model.FindByPhone(h.db, identifier)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Invalid user email or phone."})
		return
	}

	if user.HasPassword() {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Your account is already active."})
		return
	}

	if user.HasLeft {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "You have left the hostel."})
		return
	}

	userOtp, err := model.FindUserOtpById(h.db, user.UserId())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
		return
	}
	if userOtp == nil || !userOtp.IsValidOtp(otpCode) {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Invalid OTP."})
		return
	}

	err = user.SetPassword(h.db, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend updating your password."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Account activated."})
}

func (h *controller) userForgetPassword(c *gin.Context) {
	identifier := c.Query("identifier")
	otpCode := c.Query("otpcode")
	req := struct {
		NewPassword string `json:"new_password" binding:"required"`
	}{}

	if identifier == "" || otpCode == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing 'identifier' or 'otpcode' query parameter"})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request payload"})
		return
	}

	if len(req.NewPassword) < 8 || len(req.NewPassword) > 40 {
		c.JSON(http.StatusLengthRequired, gin.H{"msg": "Password length must be between 8 to 40 charecter long."})
		return
	}

	var user *model.Users
	var err error
	if strings.Contains(identifier, "@") {
		user, err = model.FindByEmail(h.db, identifier)
	} else {
		user, err = model.FindByPhone(h.db, identifier)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Invalid user email or phone."})
		return
	}

	if !user.HasPassword() {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Your account is not active yet. Active it first."})
		return
	}

	if user.HasLeft {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "You have left the hostel."})
		return
	}

	userOtp, err := model.FindUserOtpById(h.db, user.UserId())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
		return
	}
	if userOtp == nil || !userOtp.IsValidOtp(otpCode) {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Invalid OTP."})
		return
	}

	err = user.SetPassword(h.db, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend updating your password."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Password change successfull."})
}

func (h *controller) userChangePassword(c *gin.Context) {
	req := struct {
		NewPassword string `json:"new_password" binding:"required"`
		OldPassword string `json:"old_password" binding:"required"`
	}{}
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request payload"})
		return
	}

	if len(req.NewPassword) < 8 || len(req.NewPassword) > 40 {
		c.JSON(http.StatusLengthRequired, gin.H{"msg": "Password length must be between 8 to 40 charecter long."})
		return
	}

	user, err := model.FindUserById(h.db, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Invalid user email or phone."})
		return
	}

	if !user.HasPassword() {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Your account is not active yet. Active it first."})
		return
	}
	if !user.ComparePassword(req.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Old password did not matched."})
		return
	}
	if user.HasLeft {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "You have left the hostel."})
		return
	}

	err = user.SetPassword(h.db, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend updating your password."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Password change successfull."})
}

func (h *controller) userSendOtp(c *gin.Context) {
	identifier := c.Query("identifier")
	if identifier == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing 'identifier' query parameter"})
		return
	}
	var reason int
	if reason, err := strconv.Atoi(c.Query("reason")); err != nil || (reason != 1 && reason != 2) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid 'reason', must be 1 or 2"})
		return
	}

	var user *model.Users
	var err error
	if strings.Contains(identifier, "@") {
		user, err = model.FindByEmail(h.db, identifier)
	} else {
		user, err = model.FindByPhone(h.db, identifier)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Invalid user email or phone."})
		return
	}

	if user.HasLeft {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "You have left the hostel."})
		return
	}

	userOtp, err := model.FindUserOtpById(h.db, user.UserId())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
		return
	}
	if userOtp.IsExpired() {
		err = userOtp.Update(h.db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something happend getting your data."})
			return
		}
	}

	data := struct {
		Name string
		OTP  string
	}{
		Name: user.FirstName,
		OTP:  userOtp.OtpCode(),
	}

	var (
		subject  string = "Here is your code for "
		template string = "send_otp.html"
	)
	if reason == 1 {
		subject += "Active your account"
	} else {
		subject = "Reset your password"
	}

	err = h.email.SendTemplateEmail(
		user.Email,
		subject,
		"internal/view/"+template,
		data,
	)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
}
