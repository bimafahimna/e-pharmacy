package config

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Config struct {
	App        AppConfig
	Client     ClientConfig
	Cloudinary CloudinaryConfig
	Cors       cors.Config
	Database   DatabaseConfig
	Google     oauth2.Config
	Jwt        JwtConfig
	Redis      RedisConfig
	Smtp       SmtpConfig
	Cron       CronConfig
	Logistic   LogisticConfig
	URL        URLConfig
	Worker     WorkerConfig
	Cache      CacheConfig
}

func InitConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env")
	}

	return &Config{
		App:        initAppConfig(),
		Client:     initClientConfig(),
		Cloudinary: initCloudinaryConfig(),
		Cors:       initCorsConfig(),
		Database:   initDatabaseConfig(),
		Google:     initGoogleConfig(),
		Jwt:        initJwtConfig(),
		Redis:      initRedisConfig(),
		Smtp:       initSmtpConfig(),
		Cron:       initCronConfig(),
		Logistic:   initLogisticConfig(),
		URL:        initURLConfig(),
		Worker:     initWorkerConfig(),
		Cache:      initCacheConfig(),
	}
}
