package email

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
)

type Service struct {
	config config.SMTPConfig
}

func NewService(cfg config.SMTPConfig) *Service {
	return &Service{
		config: cfg,
	}
}

func (s *Service) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	address := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	headers := []string{
		fmt.Sprintf("From: %s", s.config.SMTPUsername),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
	}
	msg := strings.Join(headers, "\r\n") + body

	err := smtp.SendMail(address, auth, s.config.SMTPUsername, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Println("Successfully sent to " + to)
	return nil
}

func (s *Service) SendForgotPasswordEmail(resetUrl, to string) error {
	msg := generatePasswordResetEmailHTML(resetUrl, time.Now().Year())
	return s.SendEmail(to, PasswordResetEmailSubject, msg)
}
