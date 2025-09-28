package controller

import (
	"github.com/abrarshakhi/hostel-management-server/internal/database"
)

type Controller struct {
	db database.Service
}

func NewController(db database.Service) *Controller {
	return &Controller{db: db}
}
