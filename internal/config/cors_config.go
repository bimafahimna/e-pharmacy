package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
)

func initCorsConfig() cors.Config {
	allowCredentials, err := strconv.ParseBool(os.Getenv("CORS_ALLOW_CREDENTIALS"))
	if err != nil {
		log.Fatal("failed to parse CORS_ALLOW_CREDENTIALS")
	}

	return cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ","),
		AllowMethods:     strings.Split(os.Getenv("CORS_ALLOW_METHODS"), ","),
		AllowHeaders:     strings.Split(os.Getenv("CORS_ALLOW_HEADERS"), ","),
		AllowCredentials: allowCredentials,
	}
}
