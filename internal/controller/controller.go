package controller

import "github.com/abrarshakhi/hostel-management-server/internal/service"

type Controller struct {
	db    service.Database
	email service.Email
}

func NewController(db service.Database, email service.Email) *Controller {
	return &Controller{db: db, email: email}
}
