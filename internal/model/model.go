package model

import "github.com/abrarshakhi/hostel-management-server/internal/service"

type Model interface {
	Update(db service.Database) error
}
