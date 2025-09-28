package model

import "github.com/abrarshakhi/hostel-management-server/internal/database"

type Model interface {
	Update(db database.Service) error
}
