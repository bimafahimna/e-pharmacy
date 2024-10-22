package config

import (
	"log"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

func initDatabaseConfig() DatabaseConfig {
	port, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 32)
	if err != nil {
		log.Fatal("failed to parse DB_PORT")
	}

	return DatabaseConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     int(port),
		Name:     os.Getenv("DB_NAME"),
	}
}
