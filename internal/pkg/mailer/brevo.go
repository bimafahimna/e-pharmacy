package mailer

import (
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"gopkg.in/gomail.v2"
)

type brevoMailer struct {
	msg    *gomail.Message
	dialer *gomail.Dialer
	host   string
}

func NewBrevoMailer(config config.SmtpConfig) Mailer {
	msg := gomail.NewMessage()

	dialer := gomail.NewDialer(config.Server, config.Port, config.Email, config.Password)

	return &brevoMailer{
		msg:    msg,
		dialer: dialer,
		host:   config.ClientHost,
	}
}

func (m *brevoMailer) Host() string {
	return m.host
}

func (m *brevoMailer) Set(to, subject, body string) {
	m.msg.SetHeader("From", "admin@puxing.com")
	m.msg.SetHeader("To", to)
	m.msg.SetHeader("Subject", subject)
	m.msg.SetBody("text/html", body)
}

func (m *brevoMailer) Send() error {
	return m.dialer.DialAndSend(m.msg)
}
