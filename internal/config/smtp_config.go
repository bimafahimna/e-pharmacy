package config

import (
	"log"
	"os"
	"strconv"
)

type SmtpConfig struct {
	Server     string
	Port       int
	Email      string
	Password   string
	ClientHost string
}

func initSmtpConfig() SmtpConfig {
	port, err := strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 32)
	if err != nil {
		log.Fatal("failed to parse SMTP_PORT")
	}

	return SmtpConfig{
		Server:     os.Getenv("SMTP_SERVER"),
		Port:       int(port),
		Email:      os.Getenv("SMTP_EMAIL"),
		Password:   os.Getenv("SMTP_PASSWORD"),
		ClientHost: os.Getenv("SMTP_CLIENT_HOST"),
	}
}
