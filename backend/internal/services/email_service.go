package services

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailService struct {
	host      string
	port      int
	user      string
	password  string
	fromEmail string
}

func NewEmailService(host string, port int, user, password, from string) *EmailService {
	return &EmailService{host: host, port: port, user: user, password: password, fromEmail: from}
}

// SendOTPAsync mengirim OTP dengan goroutine agar non-blocking.
func (s *EmailService) SendOTPAsync(toEmail, otpCode string) {
	go func() {
		if err := s.sendOTP(toEmail, otpCode); err != nil {
			log.Printf("failed sending OTP to %s: %v", toEmail, err)
		}
	}()
}

func (s *EmailService) sendOTP(toEmail, otpCode string) error {
	if s.host == "" || s.user == "" || s.password == "" {
		log.Printf("[DEV-MODE] OTP for %s is: %s", toEmail, otpCode)
		return nil
	}

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: Verifikasi Email E-Commerce\r\n" +
		"\r\n" +
		"Kode OTP Anda: "+otpCode+"\r\n" +
		"Berlaku selama 10 menit.\r\n")

	return smtp.SendMail(addr, auth, s.fromEmail, []string{toEmail}, msg)
}
