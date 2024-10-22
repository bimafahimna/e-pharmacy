package config

import (
	"log"
	"os"
	"strings"
	"time"
)

type JwtConfig struct {
	Issuer         string
	SecretKey      string
	AllowedAlgs    []string
	ExpireDuration time.Duration
}

func initJwtConfig() JwtConfig {
	allowedAlgs := strings.Split(os.Getenv("JWT_ALLOWED_ALGS"), ",")

	expireDuration, err := time.ParseDuration(os.Getenv("JWT_EXPIRE_DURATION"))
	if err != nil {
		log.Fatal("failed to parse JWT_EXPIRE_DURATION")
	}

	return JwtConfig{
		Issuer:         os.Getenv("JWT_ISSUER"),
		SecretKey:      os.Getenv("JWT_SECRET_KEY"),
		AllowedAlgs:    allowedAlgs,
		ExpireDuration: expireDuration,
	}
}
