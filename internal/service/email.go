package service

import (
	"bytes"
	"fmt"
	"net/smtp"
	"sync"
	"text/template"
)

type Email struct {
	smtpHost string
	smtpPort string
	username string
	password string
	from     string
}

var (
	emailInstance *Email
	once          sync.Once
)

func EmailInstance() *Email {
	once.Do(func() {
		emailInstance = &Email{
			smtpHost: "smtp.example.com",
			smtpPort: "587",
			username: "your@email.com",
			password: "your-password",
			from:     "no-reply@example.com",
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
