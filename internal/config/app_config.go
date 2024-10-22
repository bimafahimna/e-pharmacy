package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	ServerAddress     string
	ServerGracePeriod time.Duration
	BcryptCost        int
	LoggingLevel      int
}

func initAppConfig() AppConfig {
	gracePeriod, err := time.ParseDuration(os.Getenv("SERVER_GRACE_PERIOD"))
	if err != nil {
		log.Fatal("failed to parse SERVER_GRACE_PERIOD")
	}

	bcryptCost, err := strconv.ParseInt(os.Getenv("BCRYPT_COST"), 10, 32)
	if err != nil {
		log.Fatal("failed to parse BCRYPT_COST")
	}

	loggingLevel, err := strconv.ParseInt(os.Getenv("LOGGING_LEVEL"), 10, 32)
	if err != nil {
		log.Fatal("failed to parse LOGGING_LEVEL")
	}

	return AppConfig{
		ServerAddress:     os.Getenv("SERVER_ADDRESS"),
		ServerGracePeriod: gracePeriod,
		BcryptCost:        int(bcryptCost),
		LoggingLevel:      int(loggingLevel),
	}
}
