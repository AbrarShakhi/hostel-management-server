package model

import (
	"database/sql"
	"time"

	"github.com/abrarshakhi/hostel-management-server/internal/database"
)

type User_ struct {
	Id          string
	Email       string
	Phone       string
	Password_   sql.NullString
	FirstName   string
	LastName    sql.NullString
	DateOfBirth time.Time
	Gender      string
	Nationality string
	CreatedOn   time.Time
	LastLogin   sql.NullTime
	HasLeft     bool
}

func (m *User_) Update(db database.Service) error {
	query := `
		UPDATE user_
		SET email = $1,
		    phone = $2,
		    password_ = $3,
		    first_name = $4,
		    last_name = $5,
		    date_of_birth = $6,
		    gender = $7,
		    nationality = $8,
		    created_on = $9,
		    last_login = $10,
		    has_left = $11
		WHERE id = $12
	`

	_, err := db.Exec(query,
		m.Email,
		m.Phone,
		m.Password_,
		m.FirstName,
		m.LastName,
		m.DateOfBirth,
		m.Gender,
		m.Nationality,
		m.CreatedOn,
		m.LastLogin,
		m.HasLeft,
		m.Id,
	)

	return err
}

func FindByEmail(db database.Service, email string) (*User_, error) {
	var user User_

	row := db.QueryRow(`
		SELECT id, email, phone, password_, first_name, last_name, date_of_birth, 
		       gender, nationality, created_on, last_login, has_left 
		FROM user_ 
		WHERE email = $1 
		LIMIT 1`, email)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Phone,
		&user.Password_,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Gender,
		&user.Nationality,
		&user.CreatedOn,
		&user.LastLogin,
		&user.HasLeft,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func FindByPhone(db database.Service, phone string) (*User_, error) {
	var user User_

	row := db.QueryRow(`
		SELECT id, email, phone, password_, first_name, last_name, date_of_birth, 
		       gender, nationality, created_on, last_login, has_left 
		FROM user_ 
		WHERE phone = $1 
		LIMIT 1`, phone)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Phone,
		&user.Password_,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Gender,
		&user.Nationality,
		&user.CreatedOn,
		&user.LastLogin,
		&user.HasLeft,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
