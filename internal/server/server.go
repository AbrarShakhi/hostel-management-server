package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/abrarshakhi/hostel-management-server/internal/service"
)

type Server struct {
	port  int
	db    service.Database
	email service.Email
}

func NewServer() *http.Server {
	db := service.DbInstance()
	if db == nil {
		log.Fatal("db instance is nil: Unable to create database service")
		return nil
	}
	email := service.EmailInstance()
	if email == nil {
		log.Fatal("email instance is nil: Unable to create email service")
		return nil
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:  port,
		db:    *db,
		email: *email,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
