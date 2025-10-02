package model

import (
	"database/sql"
	"time"

	"github.com/abrarshakhi/hostel-management-server/internal/service"
)

type UserOtp struct {
	userId    int
	otpCode   string
	createdAt time.Time
	expiresAt time.Time
	isUsed    bool
	attempts  int
}

func (m *UserOtp) UserId() int {
	return m.userId
}

func (m *UserOtp) Update(db service.Database) error {
	deleteQ := `DELETE FROM user_otp WHERE user_id = $1`
	_, err := db.Exec(deleteQ, m.userId)
	if err != nil {
		return err
	}

	insertQ := `INSERT INTO user_otp (user_id) VALUES ($1)`
	_, err = db.Exec(insertQ, m.userId)
	if err != nil {
		return err
	}

	row := db.QueryRow(`
		SELECT user_id, otp_code, created_at, expires_at, is_used, attempts
		FROM user_otp
		WHERE user_id = $1
		LIMIT 1`, m.userId)

	err = row.Scan(
		&m.userId,
		&m.otpCode,
		&m.createdAt,
		&m.expiresAt,
		&m.isUsed,
		&m.attempts,
	)
	return err
}

func (m *UserOtp) IsValidOtp(otpcode string) bool {
	if m == nil || m.isUsed || m.attempts >= 5 || time.Now().After(m.expiresAt) {
		return false
	}
	return otpcode == m.otpCode
}

func FindUserOtpById(db service.Database, userId int) (*UserOtp, error) {
	var userOtp UserOtp

	row := db.QueryRow(`
		SELECT user_id, otp_code, created_at, expires_at, is_used, attempts
		FROM user_otp
		WHERE user_id = $1
		LIMIT 1`, userId)

	err := row.Scan(
		&userOtp.userId,
		&userOtp.otpCode,
		&userOtp.createdAt,
		&userOtp.expiresAt,
		&userOtp.isUsed,
		&userOtp.attempts,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &userOtp, nil
}
