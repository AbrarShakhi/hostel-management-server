package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"sync"
)

type Email struct {
	smtpHost string
	smtpPort string
	username string
	password string
	from     string
}

var (
	smtpHost      = os.Getenv("SMTP_HOST")
	smtpPort      = os.Getenv("SMTP_PORT")
	emailUsername = os.Getenv("EMAIL_USERNAME")
	emailPassword = os.Getenv("EMAIL_PASSWORD")
	emailFrom     = os.Getenv("EMAIL_FROM")
	emailInstance *Email
	emailOnce     sync.Once
)

func EmailInstance() *Email {
	emailOnce.Do(func() {
		emailInstance = &Email{
			smtpHost: smtpHost,
			smtpPort: smtpPort,
			username: emailUsername,
			password: emailPassword,
			from:     emailFrom,
		}
	})
	return emailInstance
}

func (e *Email) SendTemplateEmail(to string, subject string, templatePath string, data any) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return e.sendEmail(to, subject, body.String())
}

func (e *Email) sendEmail(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", e.username, e.password, e.smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"From: " + e.from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body)

	addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)

	err := smtp.SendMail(addr, auth, e.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
