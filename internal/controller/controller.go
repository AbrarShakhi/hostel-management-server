package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/abrarshakhi/hostel-management-server/internal/middleware"
	"github.com/abrarshakhi/hostel-management-server/internal/service"
)

type controller struct {
	db    service.Database
	email service.Email
}

func InitUsersRoutes(rg *gin.RouterGroup, db service.Database, email service.Email) {
	middleware, controller := middleware.NewMiddleware(), &controller{db: db, email: email}

	rg.GET("/", controller.HelloWorld)
	rg.GET("/health", controller.Health)

	rg.POST("/login", controller.userLogin)
	rg.DELETE("/logout", controller.userLogOut)
	rg.GET("/auth-check", middleware.VerifyUser, controller.userAuthCheck)

	rg.PATCH("/change-password", controller.userChangePassword)
	rg.PATCH("/forget-password", controller.userForgetPassword)

	rg.POST("/send-otp", middleware.IdentifyOtpUser, controller.userSendOtp)
}

func InitAdminsRoutes(rg *gin.RouterGroup, db service.Database, email service.Email) {
	_, controller := middleware.NewMiddleware(), &controller{db: db, email: email}

	rg.GET("/", controller.HelloWorld)
	rg.GET("/health", controller.Health)
}
