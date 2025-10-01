package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/abrarshakhi/hostel-management-server/internal/middleware"
	"github.com/abrarshakhi/hostel-management-server/internal/service"
)

type Controller struct {
	db    service.Database
	email service.Email
}

func InitUsersRoutes(rg *gin.RouterGroup, db service.Database, email service.Email) {
	middleware, controller := middleware.NewMiddleware(), &Controller{db: db, email: email}

	rg.GET("/", controller.HelloWorld)
	rg.GET("/health", controller.Health)

	rg.POST("/login", controller.UserLogin)
	rg.GET("/logout", controller.UserLogOut)
	rg.GET("/auth-check", middleware.VerifyUser, controller.UserAuthCheck)
}

func InitAdminsRoutes(rg *gin.RouterGroup, db service.Database, email service.Email) {
	_, controller := middleware.NewMiddleware(), &Controller{db: db, email: email}

	rg.GET("/", controller.HelloWorld)
	rg.GET("/health", controller.Health)
}
