package server

import (
	"net/http"

	"github.com/abrarshakhi/hostel-management-server/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	h := handlers.NewHandlers(s.db)
	api := r.Group("/api")
	{
		api.GET("/", h.HelloWorldHandler)
		api.GET("/health", h.HealthHandler)
	}

	return r
}
