package config

import "os"

type LogisticConfig struct {
	URL    string
	ApiKey string
}

func initLogisticConfig() LogisticConfig {
	return LogisticConfig{
		URL:    os.Getenv("LOGISTIC_URL"),
		ApiKey: os.Getenv("LOGISTIC_API_KEY"),
	}
}
