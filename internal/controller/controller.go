package controller

import "github.com/abrarshakhi/hostel-management-server/internal/service"

type Controller struct {
	db service.Database
}

func NewController(db service.Database) *Controller {
	return &Controller{db: db}
}
