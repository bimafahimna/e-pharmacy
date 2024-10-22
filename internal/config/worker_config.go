package config

import "os"

type WorkerConfig struct {
	SecretKey string
}

func initWorkerConfig() WorkerConfig {
	return WorkerConfig{
		SecretKey: os.Getenv("WORKER_SECRET_KEY"),
	}
}
