package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/abrarshakhi/hostel-management-server/internal/handlers"
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

	handlers := handlers.New(s.db)
	api := r.Group("/api")
	{
		api.GET("/", handlers.HelloWorld)
		api.GET("/health", handlers.Health)
		api.GET("/auth-check", middleware.VerifyUser, handlers.Health)
	}

	return r
}
