package config

import "os"

type URLConfig struct {
	DomainName string
	Frontend   string
	Backend    string
}

func initURLConfig() URLConfig {
	return URLConfig{
		DomainName: os.Getenv("DOMAIN_NAME"),
		Frontend:   os.Getenv("FE_URL"),
		Backend:    os.Getenv("BE_URL"),
	}
}
