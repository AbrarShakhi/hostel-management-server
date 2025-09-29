package model

import (
	"database/sql"
	"time"

	"github.com/abrarshakhi/hostel-management-server/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type User_ struct {
	Id          string
	Email       string
	Phone       string
	password_   sql.NullString
	FirstName   string
	LastName    sql.NullString
	DateOfBirth time.Time
	Gender      string
	Nationality string
	CreatedOn   time.Time
	LastLogin   sql.NullTime
	HasLeft     bool
}

func (m *User_) Update(db service.Database) error {
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
		m.password_,
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

func (m *User_) ComparePassword(inputPassword string) bool {
	if !m.HasPassword() {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(m.password_.String), []byte(inputPassword)); err != nil {
		return false
	}
	return true
}

func (m *User_) HasPassword() bool {
	return m.password_.Valid
}

func FindByEmail(db service.Database, email string) (*User_, error) {
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
		&user.password_,
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

func FindByPhone(db service.Database, phone string) (*User_, error) {
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
		&user.password_,
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
