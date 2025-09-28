package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/abrarshakhi/hostel-management-server/internal/controller"
	"github.com/abrarshakhi/hostel-management-server/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	middleware := middleware.NewMiddleware()

	controller := controller.NewController(s.db)
	api := r.Group("/api")
	{
		api.GET("/", controller.HelloWorld)
		api.GET("/health", controller.Health)
		api.POST("/login", controller.UserLogin)
	}

	return r
}
