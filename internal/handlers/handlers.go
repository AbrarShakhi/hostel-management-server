package handlers

import (
	"github.com/abrarshakhi/hostel-management-server/internal/database"
)

type Handlers struct {
	db database.Service
}

func New(db database.Service) *Handlers {
	return &Handlers{db: db}
}
