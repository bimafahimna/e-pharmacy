package config

import (
	"log"
	"os"
	"time"
	_ "time/tzdata"
)

type CronConfig struct {
	TimeLocation *time.Location
}

func initCronConfig() CronConfig {
	timeLocation, err := time.LoadLocation(os.Getenv("TIME_LOCATION"))
	if err != nil {
		log.Fatal("failed to parse TIME_LOCATION")
	}
	return CronConfig{
		TimeLocation: timeLocation,
	}
}
