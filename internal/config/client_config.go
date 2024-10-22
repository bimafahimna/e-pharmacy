package config

import "os"

type ClientConfig struct {
	AuthRedirectURL string
}

func initClientConfig() ClientConfig {
	return ClientConfig{
		AuthRedirectURL: os.Getenv("CLIENT_AUTH_REDIRECT_URL"),
	}
}
